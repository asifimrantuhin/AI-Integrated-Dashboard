package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"idash-backend-api/database"
	"idash-backend-api/models"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type BIKPI struct {
	Label  string  `json:"label"`
	Value  float64 `json:"value"`
	Change float64 `json:"change"`
}

type BIDepartmentScore struct {
	Department string  `json:"department"`
	Revenue    float64 `json:"revenue"`
	Cost       float64 `json:"cost"`
	Margin     float64 `json:"margin"`
	Attendance float64 `json:"attendance"`
}

type BIForecastSummary struct {
	Type       string  `json:"type"`
	Label      string  `json:"label"`
	Value      float64 `json:"value"`
	Confidence float64 `json:"confidence"`
	ModelUsed  string  `json:"model_used"`
	ForecastAt string  `json:"forecast_at"`
}

type BITrendPoint struct {
	Month     string  `json:"month"`
	Revenue   float64 `json:"revenue"`
	Cost      float64 `json:"cost"`
	Margin    float64 `json:"margin"`
	Inventory float64 `json:"inventory"`
}

type BIOverviewResponse struct {
	KPIs             []BIKPI                 `json:"kpis"`
	Departments      []BIDepartmentScore     `json:"departments"`
	Forecasts        []BIForecastSummary     `json:"forecasts"`
	Trend            []BITrendPoint          `json:"trend"`
	Alerts           []string                `json:"alerts"`
	AIInsights       []string                `json:"ai_insights"`
	Pipeline         *SalesPipelineSnapshot  `json:"pipeline,omitempty"`
	SalesTargets     []SalesTargetVariance   `json:"sales_targets,omitempty"`
	Promotions       []SalesPromotionInsight `json:"promotions,omitempty"`
	ExecutiveSummary []string                `json:"executive_summary,omitempty"`
	LastUpdated      time.Time               `json:"last_updated"`
}

type BIAssistantRequest struct {
	Question string                 `json:"question"`
	Context  map[string]interface{} `json:"context"`
}

type BIAssistantResponse struct {
	Answer string `json:"answer"`
	Source string `json:"source"`
}

func GetBIOverview(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	prevStart := startDate.AddDate(0, -1, 0)
	prevEnd := endDate.AddDate(0, -1, 0)

	resp := BIOverviewResponse{
		KPIs:             make([]BIKPI, 0),
		Departments:      make([]BIDepartmentScore, 0),
		Forecasts:        make([]BIForecastSummary, 0),
		Trend:            make([]BITrendPoint, 0),
		Alerts:           make([]string, 0),
		AIInsights:       make([]string, 0),
		SalesTargets:     make([]SalesTargetVariance, 0),
		Promotions:       make([]SalesPromotionInsight, 0),
		ExecutiveSummary: make([]string, 0),
		LastUpdated:      time.Now(),
	}

	type summaryRow struct {
		Revenue   float64
		Cost      float64
		Margin    float64
		Inventory float64
	}

	var current summaryRow
	database.DB.Raw(`
		SELECT COALESCE(SUM(billed),0) AS revenue,
			COALESCE(SUM(cogs),0) AS cost,
			COALESCE(SUM(gp),0) AS margin
		FROM channelwise_monthly_report
		LEFT JOIN cogs_gps ON channelwise_monthly_report.company_id = cogs_gps.company_id AND channelwise_monthly_report.data_month = cogs_gps.month
		WHERE data_month BETWEEN ? AND ?
	`, startDate, endDate).Scan(&current)

	var previous summaryRow
	database.DB.Raw(`
		SELECT COALESCE(SUM(billed),0) AS revenue,
			COALESCE(SUM(cogs),0) AS cost,
			COALESCE(SUM(gp),0) AS margin
		FROM channelwise_monthly_report
		LEFT JOIN cogs_gps ON channelwise_monthly_report.company_id = cogs_gps.company_id AND channelwise_monthly_report.data_month = cogs_gps.month
		WHERE data_month BETWEEN ? AND ?
	`, prevStart, prevEnd).Scan(&previous)

	database.DB.Raw(`
		SELECT COALESCE(SUM(amount),0) AS inventory
		FROM inventory_raw_datas
		WHERE month BETWEEN ? AND ? AND status = 1
	`, startDate, endDate).Scan(&current.Inventory)

	database.DB.Raw(`
		SELECT COALESCE(SUM(amount),0) AS inventory
		FROM inventory_raw_datas
		WHERE month BETWEEN ? AND ? AND status = 1
	`, prevStart, prevEnd).Scan(&previous.Inventory)

	delta := func(cur, prev float64) float64 {
		if prev == 0 {
			return 0
		}
		return ((cur - prev) / prev) * 100
	}

	resp.KPIs = append(resp.KPIs,
		BIKPI{Label: "Revenue", Value: current.Revenue, Change: delta(current.Revenue, previous.Revenue)},
		BIKPI{Label: "Cost", Value: current.Cost, Change: delta(current.Cost, previous.Cost)},
		BIKPI{Label: "Gross Margin", Value: current.Margin, Change: delta(current.Margin, previous.Margin)},
		BIKPI{Label: "Inventory Value", Value: current.Inventory, Change: delta(current.Inventory, previous.Inventory)},
	)

	type pipelineAgg struct {
		TotalValue     float64
		DeliveredValue float64
	}

	var currPipeline pipelineAgg
	database.DB.Raw(`
		SELECT
			COALESCE(SUM(order_amount), 0) AS total_value,
			COALESCE(SUM(CASE WHEN status = 'delivered' THEN order_amount ELSE 0 END), 0) AS delivered_value
		FROM sales_order_book
		WHERE order_date BETWEEN ? AND ?
	`, startDate, endDate).Scan(&currPipeline)

	var prevPipeline pipelineAgg
	database.DB.Raw(`
		SELECT
			COALESCE(SUM(order_amount), 0) AS total_value,
			COALESCE(SUM(CASE WHEN status = 'delivered' THEN order_amount ELSE 0 END), 0) AS delivered_value
		FROM sales_order_book
		WHERE order_date BETWEEN ? AND ?
	`, prevStart, prevEnd).Scan(&prevPipeline)

	pipelineDelta := delta(currPipeline.TotalValue, prevPipeline.TotalValue)
	resp.KPIs = append(resp.KPIs, BIKPI{Label: "Pipeline Value", Value: currPipeline.TotalValue, Change: pipelineDelta})

	currConversion := 0.0
	if currPipeline.TotalValue > 0 {
		currConversion = (currPipeline.DeliveredValue / currPipeline.TotalValue) * 100
	}
	prevConversion := 0.0
	if prevPipeline.TotalValue > 0 {
		prevConversion = (prevPipeline.DeliveredValue / prevPipeline.TotalValue) * 100
	}
	resp.KPIs = append(resp.KPIs, BIKPI{Label: "Pipeline Conversion %", Value: currConversion, Change: currConversion - prevConversion})

	// Finance variance
	var variance struct{ Value float64 }
	database.DB.Raw(`
		SELECT COALESCE(SUM(budget_amount - actual_amount),0) AS value
		FROM budget_summaries
		WHERE month BETWEEN ? AND ?
	`, startDate, endDate).Scan(&variance)
	resp.KPIs = append(resp.KPIs, BIKPI{Label: "Budget Variance", Value: variance.Value, Change: 0})

	// Attrition rate
	var headcount int64
	database.DB.Model(&models.EmployeeBasicInfo{}).Where("status = ?", 1).Count(&headcount)
	var attrition int64
	database.DB.Model(&models.EmployeeTranOver{}).Where("date BETWEEN ? AND ? AND type = ?", startDate, endDate, "attrition").Count(&attrition)
	attritionRate := 0.0
	if headcount > 0 {
		attritionRate = (float64(attrition) / float64(headcount)) * 100
	}
	resp.KPIs = append(resp.KPIs, BIKPI{Label: "Attrition %", Value: attritionRate, Change: 0})

	// Department scores (finance + hr data)
	type deptRow struct {
		Department string
		Revenue    float64
		Cost       float64
		Margin     float64
		Attendance float64
	}

	deptQuery := `
		SELECT ebi.department,
			COALESCE(SUM(cmm.billed),0) AS revenue,
			COALESCE(SUM(cg.cogs),0) AS cost,
			COALESCE(SUM(cg.gp),0) AS margin,
			COALESCE(AVG(CASE WHEN ea.status = 'present' THEN 1 ELSE 0 END) * 100,0) AS attendance
		FROM employee_basic_infos ebi
		LEFT JOIN employee_attendances ea ON ebi.employee_id = ea.employee_id AND ea.date BETWEEN ? AND ?
		LEFT JOIN channelwise_monthly_report cmm ON ebi.company_id = cmm.company_id AND cmm.data_month BETWEEN ? AND ?
		LEFT JOIN cogs_gps cg ON cmm.company_id = cg.company_id AND cmm.data_month = cg.month
		WHERE ebi.status = 1
		GROUP BY ebi.department
	`
	var deptRows []deptRow
	database.DB.Raw(deptQuery, startDate, endDate, startDate, endDate).Scan(&deptRows)
	for _, row := range deptRows {
		resp.Departments = append(resp.Departments, BIDepartmentScore{
			Department: row.Department,
			Revenue:    row.Revenue,
			Cost:       row.Cost,
			Margin:     row.Margin,
			Attendance: row.Attendance,
		})
	}
	resp.Departments = resp.Departments[:min(len(resp.Departments), 10)]

	sort.Slice(resp.Departments, func(i, j int) bool {
		return resp.Departments[i].Margin > resp.Departments[j].Margin
	})

	// Forecasts from ai_forecasts
	forecastTypes := []struct {
		Type  string
		Label string
	}{
		{"sales", "Sales Forecast"},
		{"production", "Production Forecast"},
		{"finance", "Expense Forecast"},
		{"inventory", "Inventory Forecast"},
		{"hr", "Attrition Forecast"},
		{"supplychain", "Procurement Forecast"},
	}

	for _, ft := range forecastTypes {
		var model models.AIForecast
		if err := database.DB.Where("forecast_type = ?", ft.Type).Order("created_at DESC").First(&model).Error; err == nil {
			resp.Forecasts = append(resp.Forecasts, BIForecastSummary{
				Type:       ft.Type,
				Label:      ft.Label,
				Value:      model.ForecastedValue,
				Confidence: model.ConfidenceLevel,
				ModelUsed:  model.ModelUsed,
				ForecastAt: model.CreatedAt.Format(time.RFC3339),
			})
		}
	}

	// Trend (last 6 months revenue/cost/margin/inventory)
	trendStart := startDate.AddDate(0, -5, 0)
	var trendRows []struct {
		Month     string
		Revenue   float64
		Cost      float64
		Margin    float64
		Inventory float64
	}
	database.DB.Raw(`
		SELECT DATE_FORMAT(cmm.data_month,'%Y-%m') AS month,
			COALESCE(SUM(cmm.billed),0) AS revenue,
			COALESCE(SUM(cg.cogs),0) AS cost,
			COALESCE(SUM(cg.gp),0) AS margin,
			(SELECT COALESCE(SUM(amount),0) FROM inventory_raw_datas ird WHERE ird.month = cmm.data_month AND ird.status = 1) AS inventory
		FROM channelwise_monthly_report cmm
		LEFT JOIN cogs_gps cg ON cmm.company_id = cg.company_id AND cmm.data_month = cg.month
		WHERE cmm.data_month BETWEEN ? AND ?
		GROUP BY month
		ORDER BY month
	`, trendStart, endDate).Scan(&trendRows)

	for _, row := range trendRows {
		resp.Trend = append(resp.Trend, BITrendPoint{
			Month:     row.Month,
			Revenue:   row.Revenue,
			Cost:      row.Cost,
			Margin:    row.Margin,
			Inventory: row.Inventory,
		})
	}

	// Alerts
	marginRate := 0.0
	if current.Revenue > 0 {
		marginRate = (current.Margin / current.Revenue) * 100
	}
	if marginRate < 25 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Gross margin %.1f%% below threshold", marginRate))
	}
	if attritionRate > 5 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Attrition rate %.1f%% elevated", attritionRate))
	}
	if current.Inventory > current.Cost*0.5 && current.Cost > 0 {
		resp.Alerts = append(resp.Alerts, "Inventory carrying cost high relative to COGS")
	}

	resp.AIInsights = append(resp.AIInsights,
		fmt.Sprintf("Sales vs cost delta this month: ৳ %.2f", current.Revenue-current.Cost),
		fmt.Sprintf("Top department by margin: %s", topMarginDepartment(resp.Departments)),
	)

	// Sales-focused AI insight enrichment for executive dashboard
	type execTargetRow struct {
		ChannelID         int
		ChannelName       string
		RevenueTarget     float64
		ActualRevenue     float64
		PromotionBudget   float64
		GrossMarginTarget float64
		VolumeTarget      float64
		NewCustomerTarget float64
	}

	var execTargets []execTargetRow
	database.DB.Raw(`
		SELECT
			COALESCE(sct.channel_id, 0) AS channel_id,
			COALESCE(sct.channel_name, '') AS channel_name,
			SUM(sct.revenue_target) AS revenue_target,
			COALESCE(SUM(cmm.billed), 0) AS actual_revenue,
			SUM(sct.promotion_budget) AS promotion_budget,
			SUM(sct.gross_margin_target) AS gross_margin_target,
			SUM(sct.volume_target) AS volume_target,
			SUM(sct.new_customer_target) AS new_customer_target
		FROM sales_channel_targets sct
		LEFT JOIN channelwise_monthly_report cmm
			ON cmm.channel_id = sct.channel_id AND DATE_FORMAT(cmm.data_month, '%Y-%m') = DATE_FORMAT(sct.data_month, '%Y-%m')
		WHERE sct.data_month BETWEEN ? AND ?
		GROUP BY sct.channel_id, sct.channel_name
		ORDER BY revenue_target DESC
		LIMIT 6
	`, startDate, endDate).Scan(&execTargets)

	targetsPayload := make([]map[string]interface{}, 0, len(execTargets))
	for _, row := range execTargets {
		gap := row.RevenueTarget - row.ActualRevenue
		achievement := 0.0
		if row.RevenueTarget > 0 {
			achievement = (row.ActualRevenue / row.RevenueTarget) * 100
		}
		targetsPayload = append(targetsPayload, map[string]interface{}{
			"channel_id":          row.ChannelID,
			"channel_name":        row.ChannelName,
			"revenue_target":      row.RevenueTarget,
			"actual_revenue":      row.ActualRevenue,
			"revenue_gap":         gap,
			"achievement":         achievement,
			"promotion_budget":    row.PromotionBudget,
			"gross_margin_target": row.GrossMarginTarget,
			"volume_target":       row.VolumeTarget,
			"new_customer_target": row.NewCustomerTarget,
		})
		resp.SalesTargets = append(resp.SalesTargets, SalesTargetVariance{
			ChannelID:         row.ChannelID,
			ChannelName:       row.ChannelName,
			RevenueTarget:     row.RevenueTarget,
			ActualRevenue:     row.ActualRevenue,
			Achievement:       achievement,
			RevenueGap:        gap,
			PromotionBudget:   row.PromotionBudget,
			GrossMarginTarget: row.GrossMarginTarget,
			VolumeTarget:      row.VolumeTarget,
			NewCustomerTarget: row.NewCustomerTarget,
		})
	}

	type execPipelineStage struct {
		Status        string
		Orders        int
		Value         float64
		DiscountValue float64
		MarginValue   float64
		AvgAgeDays    float64
		Delivered     float64
	}

	var execPipelineStages []execPipelineStage
	database.DB.Raw(`
		SELECT status,
			COUNT(*) AS orders,
			COALESCE(SUM(order_amount), 0) AS value,
			COALESCE(SUM(discount_amount), 0) AS discount_value,
			COALESCE(SUM(gross_margin), 0) AS margin_value,
			COALESCE(AVG(DATEDIFF(COALESCE(fulfilled_at, NOW()), order_date)), 0) AS avg_age_days,
			COALESCE(SUM(CASE WHEN status = 'delivered' THEN order_amount ELSE 0 END), 0) AS delivered
		FROM sales_order_book
		WHERE order_date BETWEEN ? AND ?
		GROUP BY status
	`, startDate, endDate).Scan(&execPipelineStages)

	pipelineStagesPayload := make([]map[string]interface{}, 0, len(execPipelineStages))
	respPipelineStages := make([]SalesPipelineStage, 0, len(execPipelineStages))
	totalOrders := 0
	for _, stage := range execPipelineStages {
		pending := stage.Value
		if strings.EqualFold(stage.Status, "delivered") {
			pending = stage.Value - stage.Delivered
			if pending < 0 {
				pending = 0
			}
		}
		avgDiscount := 0.0
		if stage.Value > 0 {
			avgDiscount = (stage.DiscountValue / stage.Value) * 100
		}
		avgMargin := 0.0
		if stage.Value > 0 {
			avgMargin = (stage.MarginValue / stage.Value) * 100
		}
		totalOrders += stage.Orders
		pipelineStagesPayload = append(pipelineStagesPayload, map[string]interface{}{
			"status":          stage.Status,
			"orders":          stage.Orders,
			"value":           stage.Value,
			"delivered_value": stage.Delivered,
			"pending_value":   pending,
			"avg_discount":    avgDiscount,
			"avg_margin":      avgMargin,
			"avg_age_days":    stage.AvgAgeDays,
		})
		respPipelineStages = append(respPipelineStages, SalesPipelineStage{
			Status:         stage.Status,
			Orders:         stage.Orders,
			Value:          stage.Value,
			DeliveredValue: stage.Delivered,
			PendingValue:   pending,
			AvgDiscount:    avgDiscount,
			AvgMargin:      avgMargin,
			AvgAgeDays:     stage.AvgAgeDays,
		})
	}
	pendingTotal := currPipeline.TotalValue - currPipeline.DeliveredValue
	if pendingTotal < 0 {
		pendingTotal = 0
	}
	if totalOrders > 0 || currPipeline.TotalValue > 0 {
		resp.Pipeline = &SalesPipelineSnapshot{
			TotalOrders:    totalOrders,
			TotalValue:     currPipeline.TotalValue,
			DeliveredValue: currPipeline.DeliveredValue,
			PendingValue:   pendingTotal,
			ConversionRate: currConversion,
			Stages:         respPipelineStages,
		}
	}

	pipelinePayload := map[string]interface{}{
		"total_orders":    totalOrders,
		"total_value":     currPipeline.TotalValue,
		"delivered_value": currPipeline.DeliveredValue,
		"pending_value":   pendingTotal,
		"conversion_rate": currConversion,
		"stages":          pipelineStagesPayload,
	}

	type execPromotionRow struct {
		CampaignCode     string
		CampaignName     string
		ChannelName      string
		SpendAmount      float64
		RevenueUplift    float64
		UpliftPercentage float64
		ROI              float64
		AudienceTags     string
	}

	var execPromotions []execPromotionRow
	database.DB.Raw(`
		SELECT campaign_code, campaign_name, COALESCE(channel_name, '') AS channel_name,
			spend_amount, revenue_uplift, uplift_percentage, roi, audience_tags
		FROM sales_promotion_performance
		WHERE (start_date BETWEEN ? AND ?)
			OR (end_date BETWEEN ? AND ?)
			OR (start_date <= ? AND (end_date IS NULL OR end_date >= ?))
		ORDER BY COALESCE(end_date, start_date, NOW()) DESC
		LIMIT 5
	`, startDate, endDate, startDate, endDate, endDate, startDate).Scan(&execPromotions)

	promotionsPayload := make([]map[string]interface{}, 0, len(execPromotions))
	for _, promo := range execPromotions {
		tags := make([]string, 0)
		if promo.AudienceTags != "" {
			var jsonTags []string
			if err := json.Unmarshal([]byte(promo.AudienceTags), &jsonTags); err == nil {
				tags = jsonTags
			} else {
				for _, part := range strings.Split(promo.AudienceTags, ",") {
					trimmed := strings.TrimSpace(part)
					if trimmed != "" {
						tags = append(tags, trimmed)
					}
				}
			}
		}
		promotionsPayload = append(promotionsPayload, map[string]interface{}{
			"campaign_code":     promo.CampaignCode,
			"campaign_name":     promo.CampaignName,
			"channel_name":      promo.ChannelName,
			"spend_amount":      promo.SpendAmount,
			"revenue_uplift":    promo.RevenueUplift,
			"uplift_percentage": promo.UpliftPercentage,
			"roi":               promo.ROI,
			"audience_tags":     tags,
		})
		resp.Promotions = append(resp.Promotions, SalesPromotionInsight{
			CampaignCode:     promo.CampaignCode,
			CampaignName:     promo.CampaignName,
			ChannelName:      promo.ChannelName,
			SpendAmount:      promo.SpendAmount,
			RevenueUplift:    promo.RevenueUplift,
			UpliftPercentage: promo.UpliftPercentage,
			ROI:              promo.ROI,
			AudienceTags:     tags,
		})
	}

	insightPayload := map[string]interface{}{
		"targets":    targetsPayload,
		"pipeline":   pipelinePayload,
		"promotions": promotionsPayload,
	}
	var salesInsights struct {
		ExecutiveSummary []string `json:"executive_summary"`
		Pipeline         struct {
			Narrative string   `json:"narrative"`
			Alerts    []string `json:"alerts"`
		} `json:"pipeline"`
		Targets struct {
			Alerts []string `json:"alerts"`
		} `json:"targets"`
		Promotions struct {
			Highlights []string `json:"highlights"`
		} `json:"promotions"`
		RecommendedActions []map[string]interface{} `json:"recommended_actions"`
	}
	if err := callAISummary("/api/enrich/sales/insights", insightPayload, &salesInsights); err == nil {
		if len(salesInsights.ExecutiveSummary) > 0 {
			resp.ExecutiveSummary = append(resp.ExecutiveSummary, salesInsights.ExecutiveSummary...)
		}
		if salesInsights.Pipeline.Narrative != "" {
			resp.AIInsights = append(resp.AIInsights, salesInsights.Pipeline.Narrative)
		}
		if len(salesInsights.Pipeline.Alerts) > 0 {
			resp.Alerts = append(resp.Alerts, salesInsights.Pipeline.Alerts...)
		}
		if len(salesInsights.Targets.Alerts) > 0 {
			resp.Alerts = append(resp.Alerts, salesInsights.Targets.Alerts...)
		}
		if len(salesInsights.Promotions.Highlights) > 0 {
			resp.AIInsights = append(resp.AIInsights, salesInsights.Promotions.Highlights...)
		}
		for _, action := range salesInsights.RecommendedActions {
			if label, ok := action["action"].(string); ok && label != "" {
				resp.AIInsights = append(resp.AIInsights, label)
			}
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func PostBIAssistant(c echo.Context) error {
	var req BIAssistantRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request"})
	}
	if req.Question == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "question is required"})
	}

	reqPayload := map[string]interface{}{
		"question": req.Question,
		"context":  req.Context,
	}
	body, err := json.Marshal(reqPayload)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to marshal request"})
	}

	url := os.Getenv("AI_SERVICE_URL")
	if url == "" {
		url = "http://localhost:8000"
	}

	resp, err := http.Post(fmt.Sprintf("%s/assistant/ask", url), "application/json", bytes.NewBuffer(body))
	if err != nil {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "failed to reach AI service"})
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return c.JSON(http.StatusBadGateway, map[string]string{"error": "AI service error"})
	}

	var aiResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&aiResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "invalid AI response"})
	}

	answer := ""
	if val, ok := aiResp["answer"].(string); ok {
		answer = val
	}
	source := ""
	if val, ok := aiResp["source"].(string); ok {
		source = val
	}

	return c.JSON(http.StatusOK, BIAssistantResponse{Answer: answer, Source: source})
}

func topMarginDepartment(depts []BIDepartmentScore) string {
	if len(depts) == 0 {
		return "N/A"
	}
	top := depts[0]
	for _, d := range depts {
		if d.Margin > top.Margin {
			top = d
		}
	}
	if top.Department == "" {
		return "Unassigned"
	}
	return top.Department
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
