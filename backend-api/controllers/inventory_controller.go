package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"idash-backend-api/database"
	"idash-backend-api/models"

	"github.com/labstack/echo/v4"
)

func GetInventorySummary(c echo.Context) error {
	month := c.QueryParam("month")
	companyID := c.QueryParam("company_id")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	targetDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)

	type InventorySummary struct {
		TotalInventory float64 `json:"total_inventory"`
		TotalValue     float64 `json:"total_value"`
		AvgInventory   float64 `json:"avg_inventory"`
		COGS           float64 `json:"cogs"`
		GP             float64 `json:"gp"`
	}

	var summary InventorySummary

	// Get inventory data
	invQuery := `
		SELECT COALESCE(SUM(amount), 0) as total_inventory
		FROM inventory_raw_datas
		WHERE month = ? AND status = 1
	`
	invParams := []interface{}{targetDate}
	if companyID != "" {
		invQuery += " AND company_id = ?"
		compID, _ := strconv.Atoi(companyID)
		invParams = append(invParams, compID)
	}
	database.DB.Raw(invQuery, invParams...).Scan(&summary.TotalInventory)

	// Get COGS and GP
	cogsQuery := `
		SELECT 
			COALESCE(SUM(cogs), 0) as cogs,
			COALESCE(SUM(gp), 0) as gp
		FROM cogs_gps
		WHERE month = ?
	`
	cogsParams := []interface{}{targetDate}
	if companyID != "" {
		cogsQuery += " AND company_id = ?"
		compID, _ := strconv.Atoi(companyID)
		cogsParams = append(cogsParams, compID)
	}
	database.DB.Raw(cogsQuery, cogsParams...).Scan(&summary)

	summary.TotalValue = summary.TotalInventory
	summary.AvgInventory = summary.TotalInventory // Can be calculated from opening and closing

	return c.JSON(http.StatusOK, summary)
}

func GetInventoryRatio(c echo.Context) error {
	month := c.QueryParam("month")
	companyID := c.QueryParam("company_id")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	if companyID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "company_id is required"})
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	lastDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	firstDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)

	type InventoryRatio struct {
		InventoryTurnoverRatio float64 `json:"inventory_turnover_ratio"`
		GMROI                  float64 `json:"gmroi"`
		DaysOnHand             float64 `json:"days_on_hand"`
	}

	var ratio InventoryRatio

	// Get opening and closing inventory
	type InventoryData struct {
		Month  time.Time `json:"month"`
		Amount float64   `json:"amount"`
	}
	var invData []InventoryData
	database.DB.Raw(`
		SELECT month, SUM(amount) as amount
		FROM inventory_raw_datas
		WHERE company_id = ? AND month IN (?, ?) AND status = 1
		GROUP BY month
		ORDER BY month ASC
	`, companyID, firstDate, lastDate).Scan(&invData)

	var opening, closing float64
	if len(invData) > 0 {
		opening = invData[0].Amount
	}
	if len(invData) > 1 {
		closing = invData[1].Amount
	} else if len(invData) == 1 {
		closing = invData[0].Amount
	}

	avgInventory := (opening + closing) / 2

	// Get COGS
	var cogs float64
	database.DB.Raw(`
		SELECT COALESCE(SUM(cogs), 0) as cogs
		FROM cogs_gps
		WHERE company_id = ? AND month = ?
	`, companyID, lastDate).Scan(&cogs)

	// Calculate daily COGS
	daysInMonth := lastDate.Day()
	dailyCOGS := cogs / float64(daysInMonth)

	// Calculate inventory turnover ratio
	if dailyCOGS > 0 {
		ratio.InventoryTurnoverRatio = avgInventory / dailyCOGS
		ratio.DaysOnHand = ratio.InventoryTurnoverRatio
	}

	// Get GP for GMROI
	var gp float64
	database.DB.Raw(`
		SELECT COALESCE(SUM(gp), 0) as gp
		FROM cogs_gps
		WHERE company_id = ? AND month = ?
	`, companyID, lastDate).Scan(&gp)

	totalGP := gp * 12 // Annualized

	// Calculate GMROI
	if avgInventory > 0 {
		ratio.GMROI = (totalGP / avgInventory) * 100
	}

	return c.JSON(http.StatusOK, ratio)
}

func GetInventoryValuation(c echo.Context) error {
	month := c.QueryParam("month")
	companyID := c.QueryParam("company_id")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	targetDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)

	type InventoryValuation struct {
		GLID     int     `json:"gl_id"`
		GLCode   string  `json:"gl_code"`
		GLName   string  `json:"gl_name"`
		Amount   float64 `json:"amount"`
	}

	var valuation []InventoryValuation
	query := `
		SELECT 
			ird.gl_id,
			iga.gl_code,
			iga.gl_name,
			COALESCE(SUM(ird.amount), 0) as amount
		FROM inventory_raw_datas ird
		LEFT JOIN inventory_gl_accounts iga ON ird.gl_id = iga.id
		WHERE ird.month = ? AND ird.status = 1
	`
	params := []interface{}{targetDate}

	if companyID != "" {
		query += " AND ird.company_id = ?"
		compID, _ := strconv.Atoi(companyID)
		params = append(params, compID)
	}

	query += " GROUP BY ird.gl_id, iga.gl_code, iga.gl_name ORDER BY amount DESC"
	database.DB.Raw(query, params...).Scan(&valuation)

	return c.JSON(http.StatusOK, valuation)
}

func GetCOGSGP(c echo.Context) error {
	month := c.QueryParam("month")
	companyID := c.QueryParam("company_id")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	var cogsGP []models.CogsGp
	query := database.DB.Where("month >= ? AND month <= ?", startDate, endDate)

	if companyID != "" {
		compID, _ := strconv.Atoi(companyID)
		query = query.Where("company_id = ?", compID)
	}

	query.Find(&cogsGP)

	return c.JSON(http.StatusOK, cogsGP)
}

func GetCompanyInventorySummary(c echo.Context) error {
	month := c.QueryParam("month")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	targetDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)

	type CompanyInventory struct {
		CompanyID   int     `json:"company_id"`
		CompanyName string  `json:"company_name"`
		TotalAmount float64 `json:"total_amount"`
	}

	var summary []CompanyInventory
	query := `
		SELECT 
			ird.company_id,
			c.name as company_name,
			COALESCE(SUM(ird.amount), 0) as total_amount
		FROM inventory_raw_datas ird
		LEFT JOIN tbl_company c ON ird.company_id = c.id
		WHERE ird.month = ? AND ird.status = 1
		GROUP BY ird.company_id, c.name
		ORDER BY total_amount DESC
	`
	database.DB.Raw(query, targetDate).Scan(&summary)

	return c.JSON(http.StatusOK, summary)
}

func GetInventoryTrends(c echo.Context) error {
	companyID := c.QueryParam("company_id")
	months := c.QueryParam("months")
	if months == "" {
		months = "12"
	}

	if companyID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "company_id is required"})
	}

	monthsInt, _ := strconv.Atoi(months)
	startDate := time.Now().AddDate(0, -monthsInt, 0)

	type TrendData struct {
		Month   time.Time `json:"month"`
		Amount  float64   `json:"amount"`
		COGS    float64   `json:"cogs"`
		GP      float64   `json:"gp"`
	}

	var trends []TrendData
	query := `
		SELECT 
			ird.month,
			COALESCE(SUM(ird.amount), 0) as amount,
			COALESCE(SUM(cg.cogs), 0) as cogs,
			COALESCE(SUM(cg.gp), 0) as gp
		FROM inventory_raw_datas ird
		LEFT JOIN cogs_gps cg ON ird.company_id = cg.company_id AND ird.month = cg.month
		WHERE ird.company_id = ? AND ird.month >= ? AND ird.status = 1
		GROUP BY ird.month
		ORDER BY ird.month ASC
	`
	database.DB.Raw(query, companyID, startDate).Scan(&trends)

	return c.JSON(http.StatusOK, trends)
}

func GetInventoryGLAccounts(c echo.Context) error {
	var accounts []models.InventoryGlAccounts
	database.DB.Where("status = ?", 1).Find(&accounts)
	return c.JSON(http.StatusOK, accounts)
}

type InventoryKPI struct {
	Label  string  `json:"label"`
	Value  float64 `json:"value"`
	Change float64 `json:"change"`
}

type InventoryCategory struct {
	GLID   int     `json:"gl_id"`
	GLCode string  `json:"gl_code"`
	GLName string  `json:"gl_name"`
	Amount float64 `json:"amount"`
}

type InventoryCompany struct {
	CompanyID   int     `json:"company_id"`
	CompanyName string  `json:"company_name"`
	Amount      float64 `json:"amount"`
}

type InventoryTrendPoint struct {
	Month  string  `json:"month"`
	Amount float64 `json:"amount"`
	COGS   float64 `json:"cogs"`
	GP     float64 `json:"gp"`
}

type InventoryTurnoverSummary struct {
	AverageInventory float64 `json:"average_inventory"`
	COGS             float64 `json:"cogs"`
	TurnoverDays     float64 `json:"turnover_days"`
	GMROI            float64 `json:"gmroi"`
}

type InventoryForecastPoint struct {
	Date       string  `json:"date"`
	Forecast   float64 `json:"forecast"`
	UpperBound float64 `json:"upper_bound"`
	LowerBound float64 `json:"lower_bound"`
}

type InventoryForecastSummary struct {
	TotalForecast        float64                  `json:"total_forecast"`
	AverageDailyForecast float64                  `json:"average_daily_forecast"`
	ConfidenceLevel      float64                  `json:"confidence_level"`
	ModelUsed            string                   `json:"model_used"`
	ForecastData         []InventoryForecastPoint `json:"forecast_data"`
}

type InventoryOverviewResponse struct {
	KPIs        []InventoryKPI           `json:"kpis"`
	Categories  []InventoryCategory      `json:"categories"`
	Companies   []InventoryCompany       `json:"companies"`
	Turnover    InventoryTurnoverSummary `json:"turnover"`
	Trend       []InventoryTrendPoint    `json:"trend"`
	Forecast    *InventoryForecastSummary `json:"forecast"`
	Alerts      []string                 `json:"alerts"`
	Prescriptions map[string]interface{} `json:"prescriptions,omitempty"`
	SlowMovers []interface{}            `json:"slow_movers,omitempty"`
	LastUpdated time.Time                `json:"last_updated"`
}

type inventorySummaryRow struct {
	Amount float64
	COGS   float64
	GP     float64
}

type inventoryForecastDetail struct {
	Date       string  `json:"date"`
	Forecast   float64 `json:"forecast"`
	UpperBound float64 `json:"upper_bound"`
	LowerBound float64 `json:"lower_bound"`
}

func GetInventoryOverview(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	companyID := c.QueryParam("company_id")
	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	prevStart := startDate.AddDate(0, -1, 0)
	prevEnd := endDate.AddDate(0, -1, 0)

	resp := InventoryOverviewResponse{
		KPIs:        make([]InventoryKPI, 0),
		Categories:  make([]InventoryCategory, 0),
		Companies:   make([]InventoryCompany, 0),
		Trend:       make([]InventoryTrendPoint, 0),
		Alerts:      make([]string, 0),
		LastUpdated: time.Now(),
	}

	// Current summary
	var current inventorySummaryRow
	summaryQuery := `
		SELECT COALESCE(SUM(amount),0) AS amount
		FROM inventory_raw_datas
		WHERE month BETWEEN ? AND ? AND status = 1`
	summaryParams := []interface{}{startDate, endDate}
	if companyID != "" {
		summaryQuery += " AND company_id = ?"
		summaryParams = append(summaryParams, companyID)
	}
	database.DB.Raw(summaryQuery, summaryParams...).Scan(&current.Amount)

	database.DB.Raw(`
		SELECT COALESCE(SUM(cogs),0) AS cogs, COALESCE(SUM(gp),0) AS gp
		FROM cogs_gps
		WHERE month BETWEEN ? AND ?`+withCompanyClause(companyID)+`
	`, append([]interface{}{startDate, endDate}, companyIDParam(companyID)...)...).Scan(&current)

	var previous inventorySummaryRow
	database.DB.Raw(`
		SELECT COALESCE(SUM(amount),0) AS amount
		FROM inventory_raw_datas
		WHERE month BETWEEN ? AND ? AND status = 1`+withCompanyClause(companyID)+`
	`, append([]interface{}{prevStart, prevEnd}, companyIDParam(companyID)...)...).Scan(&previous.Amount)
	database.DB.Raw(`
		SELECT COALESCE(SUM(cogs),0) AS cogs, COALESCE(SUM(gp),0) AS gp
		FROM cogs_gps
		WHERE month BETWEEN ? AND ?`+withCompanyClause(companyID)+`
	`, append([]interface{}{prevStart, prevEnd}, companyIDParam(companyID)...)...).Scan(&previous)

	delta := func(cur, prev float64) float64 {
		if prev == 0 {
			return 0
		}
		return ((cur - prev) / prev) * 100
	}

	resp.KPIs = append(resp.KPIs,
		InventoryKPI{Label: "Inventory Value", Value: current.Amount, Change: delta(current.Amount, previous.Amount)},
		InventoryKPI{Label: "COGS", Value: current.COGS, Change: delta(current.COGS, previous.COGS)},
		InventoryKPI{Label: "Gross Profit", Value: current.GP, Change: delta(current.GP, previous.GP)},
	)

	// Turnover metrics
	avgInventory := (current.Amount + previous.Amount) / 2
	turnover := InventoryTurnoverSummary{AverageInventory: avgInventory, COGS: current.COGS}
	if current.COGS > 0 {
		daysInMonth := float64(endDate.Day())
		dailyCOGS := current.COGS / daysInMonth
		if dailyCOGS > 0 {
			turnover.TurnoverDays = avgInventory / dailyCOGS
		}
	}
	if avgInventory > 0 {
		turnover.GMROI = (current.GP / avgInventory) * 100
	}
	resp.Turnover = turnover

	// Category breakdown
	var categoryRows []InventoryCategory
	categoryQuery := `
		SELECT ird.gl_id AS gl_id, COALESCE(iga.gl_code,'' ) AS gl_code, COALESCE(iga.gl_name,'') AS gl_name,
			COALESCE(SUM(ird.amount),0) AS amount
		FROM inventory_raw_datas ird
		LEFT JOIN inventory_gl_accounts iga ON ird.gl_id = iga.id
		WHERE ird.month BETWEEN ? AND ? AND ird.status = 1` + withCompanyClause(companyID) + `
		GROUP BY ird.gl_id, iga.gl_code, iga.gl_name
		ORDER BY amount DESC
		LIMIT 10
	`
	database.DB.Raw(categoryQuery, append([]interface{}{startDate, endDate}, companyIDParam(companyID)...)...).Scan(&categoryRows)
	resp.Categories = categoryRows

	// Company summary
	var companyRows []InventoryCompany
	database.DB.Raw(`
		SELECT ird.company_id AS company_id, COALESCE(c.name,'') AS company_name,
			COALESCE(SUM(ird.amount),0) AS amount
		FROM inventory_raw_datas ird
		LEFT JOIN tbl_company c ON ird.company_id = c.id
		WHERE ird.month BETWEEN ? AND ? AND ird.status = 1` + withCompanyClause(companyID) + `
		GROUP BY ird.company_id, c.name
		ORDER BY amount DESC
		LIMIT 10
	`, append([]interface{}{startDate, endDate}, companyIDParam(companyID)...)...).Scan(&companyRows)
	resp.Companies = companyRows

	// Trend over last 12 months
	trendStart := startDate.AddDate(0, -11, 0)
	var trendRows []struct {
		Month  string
		Amount float64
		COGS   float64
		GP     float64
	}
	database.DB.Raw(`
		SELECT DATE_FORMAT(ird.month,'%Y-%m') AS month,
			COALESCE(SUM(ird.amount),0) AS amount,
			COALESCE(SUM(cg.cogs),0) AS cogs,
			COALESCE(SUM(cg.gp),0) AS gp
		FROM inventory_raw_datas ird
		LEFT JOIN cogs_gps cg ON ird.company_id = cg.company_id AND ird.month = cg.month
		WHERE ird.month BETWEEN ? AND ? AND ird.status = 1` + withCompanyClause(companyID) + `
		GROUP BY month
		ORDER BY month
	`, append([]interface{}{trendStart, endDate}, companyIDParam(companyID)...)...).Scan(&trendRows)

	for _, row := range trendRows {
		resp.Trend = append(resp.Trend, InventoryTrendPoint{
			Month:  row.Month,
			Amount: row.Amount,
			COGS:   row.COGS,
			GP:     row.GP,
		})
	}

	// Forecast
	var forecastModel models.AIForecast
	if err := database.DB.Where("forecast_type = ?", "inventory").Order("created_at DESC").First(&forecastModel).Error; err == nil {
		var details []inventoryForecastDetail
		if forecastModel.ForecastDetails != "" {
			_ = json.Unmarshal([]byte(forecastModel.ForecastDetails), &details)
		}

		forecastSummary := &InventoryForecastSummary{
			TotalForecast:   forecastModel.ForecastedValue,
			ConfidenceLevel: forecastModel.ConfidenceLevel,
			ModelUsed:       forecastModel.ModelUsed,
			ForecastData:    make([]InventoryForecastPoint, 0, len(details)),
		}

		if len(details) > 0 {
			total := 0.0
			for _, d := range details {
				forecastSummary.ForecastData = append(forecastSummary.ForecastData, InventoryForecastPoint{
					Date:       d.Date,
					Forecast:   d.Forecast,
					UpperBound: d.UpperBound,
					LowerBound: d.LowerBound,
				})
				total += d.Forecast
			}
			forecastSummary.TotalForecast = total
			forecastSummary.AverageDailyForecast = total / float64(len(details))
		}

		resp.Forecast = forecastSummary
	}

	// Alerts
	if turnover.TurnoverDays > 45 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Inventory turnover days is %.0f, above target", turnover.TurnoverDays))
	}
	if len(resp.Categories) > 0 && resp.Categories[0].Amount > resp.Categories[len(resp.Categories)-1].Amount*4 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Category %s dominates inventory value", resp.Categories[0].GLName))
	}
	if resp.Forecast != nil && resp.Forecast.TotalForecast > current.Amount*1.1 {
		resp.Alerts = append(resp.Alerts, "AI forecast indicates inventory increase exceeding 10%")
	}

	// AI-driven inventory optimisation
	inventoryPayload := map[string]interface{}{
		"module":     "inventory",
		"start_date": trendStart.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"horizon":    60,
	}
	var inventoryAI InventoryPrescriptionResponse
	if err := callAISummary("/api/prescribe/inventory", inventoryPayload, &inventoryAI); err == nil {
		resp.Prescriptions = inventoryAI.Prescriptions
		if slow, ok := inventoryAI.Prescriptions["slow_movers"].([]interface{}); ok {
			resp.SlowMovers = slow
		}
	}

	return c.JSON(http.StatusOK, resp)
}

func withCompanyClause(companyID string) string {
	if companyID != "" {
		return " AND company_id = ?"
	}
	return ""
}

func companyIDParam(companyID string) []interface{} {
	if companyID != "" {
		return []interface{}{companyID}
	}
	return []interface{}{}
}
