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
	"time"

	"github.com/labstack/echo/v4"
)

func GetSupplyChainSummary(c echo.Context) error {
	month := c.QueryParam("month")
	company := c.QueryParam("company")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type SupplyChainSummary struct {
		TotalPO        float64 `json:"total_po"`
		TotalGRN       float64 `json:"total_grn"`
		TotalInvoice   float64 `json:"total_invoice"`
		PendingPO      float64 `json:"pending_po"`
		PendingGRN     float64 `json:"pending_grn"`
		PendingInvoice float64 `json:"pending_invoice"`
	}

	var summary SupplyChainSummary

	// Get PO data
	poQuery := `
		SELECT COALESCE(SUM(po_value), 0) as total_po
		FROM supply_chain_master_datas
		WHERE po_date >= ? AND po_date <= ?
	`
	poParams := []interface{}{startDate, endDate}
	if company != "" {
		poQuery += " AND company = ?"
		compID, _ := strconv.Atoi(company)
		poParams = append(poParams, compID)
	}
	database.DB.Raw(poQuery, poParams...).Scan(&summary.TotalPO)

	// Get GRN data
	grnQuery := `
		SELECT COALESCE(SUM(grn_amount), 0) as total_grn
		FROM supply_chain_grn_datas
		WHERE grn_date >= ? AND grn_date <= ?
	`
	grnParams := []interface{}{startDate, endDate}
	if company != "" {
		grnQuery += " AND company = ?"
		compID, _ := strconv.Atoi(company)
		grnParams = append(grnParams, compID)
	}
	database.DB.Raw(grnQuery, grnParams...).Scan(&summary.TotalGRN)

	// Get Invoice data
	invQuery := `
		SELECT COALESCE(SUM(total_invoice), 0) as total_invoice
		FROM supply_chain_invoice_datas
		WHERE iv_date >= ? AND iv_date <= ?
	`
	invParams := []interface{}{startDate, endDate}
	if company != "" {
		invQuery += " AND company = ?"
		compID, _ := strconv.Atoi(company)
		invParams = append(invParams, compID)
	}
	database.DB.Raw(invQuery, invParams...).Scan(&summary.TotalInvoice)

	// Calculate pending
	summary.PendingPO = summary.TotalPO - summary.TotalGRN
	summary.PendingGRN = summary.TotalGRN - summary.TotalInvoice
	summary.PendingInvoice = summary.TotalInvoice

	return c.JSON(http.StatusOK, summary)
}

func GetGRNData(c echo.Context) error {
	month := c.QueryParam("month")
	company := c.QueryParam("company")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	var grnData []models.SupplyChainGrnData
	query := database.DB.Where("grn_date >= ? AND grn_date <= ?", startDate, endDate)

	if company != "" {
		compID, _ := strconv.Atoi(company)
		query = query.Where("company = ?", compID)
	}

	query.Find(&grnData)

	return c.JSON(http.StatusOK, grnData)
}

func GetInvoiceData(c echo.Context) error {
	month := c.QueryParam("month")
	company := c.QueryParam("company")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	var invoiceData []models.SupplyChainInvoiceData
	query := database.DB.Where("iv_date >= ? AND iv_date <= ?", startDate, endDate)

	if company != "" {
		compID, _ := strconv.Atoi(company)
		query = query.Where("company = ?", compID)
	}

	query.Find(&invoiceData)

	return c.JSON(http.StatusOK, invoiceData)
}

func GetPOManagement(c echo.Context) error {
	month := c.QueryParam("month")
	company := c.QueryParam("company")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type POManagement struct {
		PONumber      string    `json:"po_number"`
		PODate        time.Time `json:"po_date"`
		POValue       float64   `json:"po_value"`
		GRNAmount     float64   `json:"grn_amount"`
		InvoiceAmount float64   `json:"invoice_amount"`
		PendingAmount float64   `json:"pending_amount"`
		Status        string    `json:"status"`
	}

	var poManagement []POManagement
	query := `
		SELECT 
			scm.po_number,
			scm.po_date,
			scm.po_value,
			COALESCE(SUM(scgd.grn_amount), 0) as grn_amount,
			COALESCE(SUM(scid.total_invoice), 0) as invoice_amount,
			(scm.po_value - COALESCE(SUM(scgd.grn_amount), 0)) as pending_amount,
			CASE 
				WHEN scm.po_value = COALESCE(SUM(scgd.grn_amount), 0) THEN 'completed'
				WHEN COALESCE(SUM(scgd.grn_amount), 0) > 0 THEN 'partial'
				ELSE 'pending'
			END as status
		FROM supply_chain_master_datas scm
		LEFT JOIN supply_chain_grn_datas scgd ON scm.po_number = scgd.po_number
		LEFT JOIN supply_chain_invoice_datas scid ON scm.po_number = scid.invoice_number
		WHERE scm.po_date >= ? AND scm.po_date <= ?
	`
	params := []interface{}{startDate, endDate}

	if company != "" {
		query += " AND scm.company = ?"
		compID, _ := strconv.Atoi(company)
		params = append(params, compID)
	}

	query += " GROUP BY scm.po_number, scm.po_date, scm.po_value ORDER BY scm.po_date DESC"
	database.DB.Raw(query, params...).Scan(&poManagement)

	return c.JSON(http.StatusOK, poManagement)
}

func GetPendingItems(c echo.Context) error {
	month := c.QueryParam("month")
	company := c.QueryParam("company")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type PendingItem struct {
		PONumber      string    `json:"po_number"`
		PODate        time.Time `json:"po_date"`
		POValue       float64   `json:"po_value"`
		GRNAmount     float64   `json:"grn_amount"`
		PendingAmount float64   `json:"pending_amount"`
		DaysPending   int       `json:"days_pending"`
	}

	var pendingItems []PendingItem
	query := `
		SELECT 
			scm.po_number,
			scm.po_date,
			scm.po_value,
			COALESCE(SUM(scgd.grn_amount), 0) as grn_amount,
			(scm.po_value - COALESCE(SUM(scgd.grn_amount), 0)) as pending_amount,
			DATEDIFF(?, scm.po_date) as days_pending
		FROM supply_chain_master_datas scm
		LEFT JOIN supply_chain_grn_datas scgd ON scm.po_number = scgd.po_number
		WHERE scm.po_date >= ? AND scm.po_date <= ?
			AND (scm.po_value - COALESCE(SUM(scgd.grn_amount), 0)) > 0
	`
	params := []interface{}{time.Now(), startDate, endDate}

	if company != "" {
		query += " AND scm.company = ?"
		compID, _ := strconv.Atoi(company)
		params = append(params, compID)
	}

	query += " GROUP BY scm.po_number, scm.po_date, scm.po_value HAVING pending_amount > 0 ORDER BY days_pending DESC"
	database.DB.Raw(query, params...).Scan(&pendingItems)

	return c.JSON(http.StatusOK, pendingItems)
}

func GetCompanyPOMonthly(c echo.Context) error {
	month := c.QueryParam("month")
	company := c.QueryParam("company")

	if month == "" {
		month = time.Now().Format("2006-01")
	}
	if company == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "company is required"})
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type POMonthly struct {
		PurchaseOrg string  `json:"purchase_org"`
		POAmount    float64 `json:"po_amount"`
	}

	var poMonthly []POMonthly
	query := `
		SELECT 
			purchase_org,
			COALESCE(SUM(po_value), 0) as po_amount
		FROM supply_chain_master_datas
		WHERE company = ? AND po_date >= ? AND po_date <= ?
		GROUP BY purchase_org
		ORDER BY po_amount DESC
	`
	compID, _ := strconv.Atoi(company)
	database.DB.Raw(query, compID, startDate, endDate).Scan(&poMonthly)

	return c.JSON(http.StatusOK, poMonthly)
}

func GetCompanyGRNMonthly(c echo.Context) error {
	month := c.QueryParam("month")
	company := c.QueryParam("company")

	if month == "" {
		month = time.Now().Format("2006-01")
	}
	if company == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "company is required"})
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type GRNMonthly struct {
		PurchaseOrg string  `json:"purchase_org"`
		GRNAmount   float64 `json:"grn_amount"`
	}

	var grnMonthly []GRNMonthly
	query := `
		SELECT 
			scm.purchase_org,
			COALESCE(SUM(scgd.grn_amount), 0) as grn_amount
		FROM supply_chain_grn_datas scgd
		JOIN supply_chain_master_datas scm ON scgd.po_number = scm.po_number
		WHERE scgd.company = ? AND scgd.grn_date >= ? AND scgd.grn_date <= ?
		GROUP BY scm.purchase_org
		ORDER BY grn_amount DESC
	`
	compID, _ := strconv.Atoi(company)
	database.DB.Raw(query, compID, startDate, endDate).Scan(&grnMonthly)

	return c.JSON(http.StatusOK, grnMonthly)
}

func GetCompanyInvoiceMonthly(c echo.Context) error {
	month := c.QueryParam("month")
	company := c.QueryParam("company")

	if month == "" {
		month = time.Now().Format("2006-01")
	}
	if company == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "company is required"})
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type InvoiceMonthly struct {
		TotalInvoice float64 `json:"total_invoice"`
	}

	var invoiceMonthly InvoiceMonthly
	query := `
		SELECT COALESCE(SUM(total_invoice), 0) as total_invoice
		FROM supply_chain_invoice_datas
		WHERE company = ? AND iv_date >= ? AND iv_date <= ?
	`
	compID, _ := strconv.Atoi(company)
	database.DB.Raw(query, compID, startDate, endDate).Scan(&invoiceMonthly)

	return c.JSON(http.StatusOK, invoiceMonthly)
}

type SupplyChainKPI struct {
	Label  string  `json:"label"`
	Value  float64 `json:"value"`
	Change float64 `json:"change"`
}

type SupplierPerformanceMetric struct {
	SupplierName     string  `json:"supplier_name"`
	OverallScore     float64 `json:"overall_score"`
	OnTimePercentage float64 `json:"on_time_percentage"`
	QualityScore     float64 `json:"quality_score"`
	CostScore        float64 `json:"cost_score"`
	Rating           string  `json:"rating"`
}

type PendingOrderMetric struct {
	PONumber      string    `json:"po_number"`
	PODate        time.Time `json:"po_date"`
	POValue       float64   `json:"po_value"`
	GRNAmount     float64   `json:"grn_amount"`
	PendingAmount float64   `json:"pending_amount"`
	DaysPending   int       `json:"days_pending"`
}

type SupplyTrendPoint struct {
	Month        string  `json:"month"`
	POValue      float64 `json:"po_value"`
	GRNValue     float64 `json:"grn_value"`
	InvoiceValue float64 `json:"invoice_value"`
}

type SupplyForecastPoint struct {
	Date       string  `json:"date"`
	Forecast   float64 `json:"forecast"`
	UpperBound float64 `json:"upper_bound"`
	LowerBound float64 `json:"lower_bound"`
}

type SupplyForecastSummary struct {
	TotalForecast        float64               `json:"total_forecast"`
	AverageDailyForecast float64               `json:"average_daily_forecast"`
	ConfidenceLevel      float64               `json:"confidence_level"`
	ModelUsed            string                `json:"model_used"`
	ForecastData         []SupplyForecastPoint `json:"forecast_data"`
}

type SupplyChainOverviewResponse struct {
	KPIs          []SupplyChainKPI            `json:"kpis"`
	Suppliers     []SupplierPerformanceMetric `json:"suppliers"`
	PendingOrders []PendingOrderMetric        `json:"pending_orders"`
	Trend         []SupplyTrendPoint          `json:"trend"`
	Forecast      *SupplyForecastSummary      `json:"forecast"`
	Alerts        []string                    `json:"alerts"`
	Insights      []string                    `json:"insights"`
	Recommendations []map[string]interface{}  `json:"recommendations,omitempty"`
	ExecutiveSummary []string                 `json:"executive_summary,omitempty"`
	LastUpdated   time.Time                   `json:"last_updated"`
}

type supplyTotals struct {
	POValue      float64
	GRNValue     float64
	InvoiceValue float64
}

type supplyForecastDetail struct {
	Date       string  `json:"date"`
	Forecast   float64 `json:"forecast"`
	UpperBound float64 `json:"upper_bound"`
	LowerBound float64 `json:"lower_bound"`
}

type supplyInsightResponse struct {
	ExecutiveSummary   []string                 `json:"executive_summary"`
	Insights           []string                 `json:"insights"`
	RecommendedActions []map[string]interface{} `json:"recommended_actions"`
}

func callSupplyChainAI(endpoint string, payload interface{}, target interface{}) error {
    aiServiceURL := os.Getenv("AI_SERVICE_URL")
    if aiServiceURL == "" {
        aiServiceURL = "http://localhost:8000"
    }

    body, err := json.Marshal(payload)
    if err != nil {
        return err
    }

    resp, err := http.Post(aiServiceURL+endpoint, "application/json", bytes.NewBuffer(body))
    if err != nil {
        return err
    }
    defer resp.Body.Close()

    if resp.StatusCode >= http.StatusBadRequest {
        return fmt.Errorf("ai service returned %s", resp.Status)
    }

    if target != nil {
        if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
            return err
        }
    }

    return nil
}

func GetSupplyChainOverview(c echo.Context) error {
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

	var companyInt int
	var hasCompany bool
	if companyID != "" {
		companyInt, _ = strconv.Atoi(companyID)
		hasCompany = true
	}

	resp := SupplyChainOverviewResponse{
		KPIs:          make([]SupplyChainKPI, 0),
		Suppliers:     make([]SupplierPerformanceMetric, 0),
		PendingOrders: make([]PendingOrderMetric, 0),
		Trend:         make([]SupplyTrendPoint, 0),
		Alerts:        make([]string, 0),
		Insights:      make([]string, 0),
		Recommendations: make([]map[string]interface{}, 0),
		ExecutiveSummary: make([]string, 0),
		LastUpdated:   time.Now(),
	}

	// Helper for summations
	sumValue := func(query string, start, end time.Time) float64 {
		params := []interface{}{start, end}
		if hasCompany {
			params = append(params, companyInt)
		}
		var value float64
		database.DB.Raw(query, params...).Scan(&value)
		return value
	}

	poQuery := `SELECT COALESCE(SUM(po_value),0) FROM supply_chain_master_datas WHERE po_date BETWEEN ? AND ?`
	if hasCompany {
		poQuery += " AND company = ?"
	}
	grnQuery := `SELECT COALESCE(SUM(grn_amount),0) FROM supply_chain_grn_datas WHERE grn_date BETWEEN ? AND ?`
	if hasCompany {
		grnQuery += " AND company = ?"
	}
	invoiceQuery := `SELECT COALESCE(SUM(total_invoice),0) FROM supply_chain_invoice_datas WHERE iv_date BETWEEN ? AND ?`
	if hasCompany {
		invoiceQuery += " AND company = ?"
	}

	current := supplyTotals{
		POValue:      sumValue(poQuery, startDate, endDate),
		GRNValue:     sumValue(grnQuery, startDate, endDate),
		InvoiceValue: sumValue(invoiceQuery, startDate, endDate),
	}
	previous := supplyTotals{
		POValue:      sumValue(poQuery, prevStart, prevEnd),
		GRNValue:     sumValue(grnQuery, prevStart, prevEnd),
		InvoiceValue: sumValue(invoiceQuery, prevStart, prevEnd),
	}

	pendingPO := current.POValue - current.GRNValue
	pendingGRN := current.GRNValue - current.InvoiceValue
	fulfillmentRate := 0.0
	if current.POValue > 0 {
		fulfillmentRate = (current.GRNValue / current.POValue) * 100
	}

	// Average lead time
	leadQuery := `
		SELECT COALESCE(AVG(DATEDIFF(scgd.grn_date, scm.po_date)),0)
		FROM supply_chain_master_datas scm
		JOIN supply_chain_grn_datas scgd ON scm.po_number = scgd.po_number
		WHERE scm.po_date BETWEEN ? AND ?`
	leadParams := []interface{}{startDate, endDate}
	if hasCompany {
		leadQuery += " AND scm.company = ?"
		leadParams = append(leadParams, companyInt)
	}
	var avgLead float64
	database.DB.Raw(leadQuery, leadParams...).Scan(&avgLead)

	// Supplier on-time rate
	supplierAvgQuery := `
		SELECT COALESCE(AVG(on_time_percentage),0)
		FROM supplier_performance
		WHERE evaluation_date BETWEEN ? AND ?`
	supplierParams := []interface{}{startDate, endDate}
	if hasCompany {
		supplierAvgQuery += " AND company_id = ?"
		supplierParams = append(supplierParams, companyInt)
	}
	var avgOnTime float64
	database.DB.Raw(supplierAvgQuery, supplierParams...).Scan(&avgOnTime)

	delta := func(cur, prev float64) float64 {
		if prev == 0 {
			return 0
		}
		return ((cur - prev) / prev) * 100
	}

	resp.KPIs = append(resp.KPIs,
		SupplyChainKPI{Label: "PO Value", Value: current.POValue, Change: delta(current.POValue, previous.POValue)},
		SupplyChainKPI{Label: "GRN Value", Value: current.GRNValue, Change: delta(current.GRNValue, previous.GRNValue)},
		SupplyChainKPI{Label: "Invoice Value", Value: current.InvoiceValue, Change: delta(current.InvoiceValue, previous.InvoiceValue)},
		SupplyChainKPI{Label: "Pending PO", Value: pendingPO, Change: 0},
		SupplyChainKPI{Label: "Pending GRN", Value: pendingGRN, Change: 0},
		SupplyChainKPI{Label: "Fulfillment Rate", Value: fulfillmentRate, Change: delta(fulfillmentRate, 0)},
		SupplyChainKPI{Label: "Avg Lead Time", Value: avgLead, Change: 0},
		SupplyChainKPI{Label: "On-Time Delivery %", Value: avgOnTime, Change: 0},
	)

	// Supplier performance list
	supplierQuery := `
		SELECT supplier_name, overall_score, on_time_percentage, quality_score, cost_score, rating
		FROM supplier_performance
		WHERE evaluation_date BETWEEN ? AND ?`
	supplierParams = []interface{}{startDate, endDate}
	if hasCompany {
		supplierQuery += " AND company_id = ?"
		supplierParams = append(supplierParams, companyInt)
	}
	supplierQuery += " ORDER BY overall_score DESC LIMIT 10"
	var suppliers []SupplierPerformanceMetric
	database.DB.Raw(supplierQuery, supplierParams...).Scan(&suppliers)
	resp.Suppliers = suppliers

	// Pending orders
	pendingQuery := `
		SELECT 
			scm.po_number,
			scm.po_date,
			scm.po_value,
			COALESCE(SUM(scgd.grn_amount),0) AS grn_amount,
			(scm.po_value - COALESCE(SUM(scgd.grn_amount),0)) AS pending_amount,
			DATEDIFF(?, scm.po_date) AS days_pending
		FROM supply_chain_master_datas scm
		LEFT JOIN supply_chain_grn_datas scgd ON scm.po_number = scgd.po_number
		WHERE scm.po_date BETWEEN ? AND ?`
	pendingParams := []interface{}{time.Now(), startDate, endDate}
	if hasCompany {
		pendingQuery += " AND scm.company = ?"
		pendingParams = append(pendingParams, companyInt)
	}
	pendingQuery += " GROUP BY scm.po_number, scm.po_date, scm.po_value HAVING pending_amount > 0 ORDER BY pending_amount DESC LIMIT 10"
	var pending []PendingOrderMetric
	database.DB.Raw(pendingQuery, pendingParams...).Scan(&pending)
	resp.PendingOrders = pending

	// Trend for last 12 months
	trendStart := startDate.AddDate(0, -11, 0)
	trendMap := map[string]*SupplyTrendPoint{}

	poTrendQuery := `SELECT DATE_FORMAT(po_date,'%Y-%m') AS month, COALESCE(SUM(po_value),0) AS value FROM supply_chain_master_datas WHERE po_date BETWEEN ? AND ?`
	poTrendParams := []interface{}{trendStart, endDate}
	if hasCompany {
		poTrendQuery += " AND company = ?"
		poTrendParams = append(poTrendParams, companyInt)
	}
	poTrendQuery += " GROUP BY month"
	var poTrend []struct {
		Month string
		Value float64
	}
	database.DB.Raw(poTrendQuery, poTrendParams...).Scan(&poTrend)
	for _, row := range poTrend {
		trendMap[row.Month] = &SupplyTrendPoint{Month: row.Month, POValue: row.Value}
	}

	grnTrendQuery := `SELECT DATE_FORMAT(grn_date,'%Y-%m') AS month, COALESCE(SUM(grn_amount),0) AS value FROM supply_chain_grn_datas WHERE grn_date BETWEEN ? AND ?`
	grnTrendParams := []interface{}{trendStart, endDate}
	if hasCompany {
		grnTrendQuery += " AND company = ?"
		grnTrendParams = append(grnTrendParams, companyInt)
	}
	grnTrendQuery += " GROUP BY month"
	var grnTrend []struct {
		Month string
		Value float64
	}
	database.DB.Raw(grnTrendQuery, grnTrendParams...).Scan(&grnTrend)
	for _, row := range grnTrend {
		point, exists := trendMap[row.Month]
		if !exists {
			point = &SupplyTrendPoint{Month: row.Month}
			trendMap[row.Month] = point
		}
		point.GRNValue = row.Value
	}

	invoiceTrendQuery := `SELECT DATE_FORMAT(iv_date,'%Y-%m') AS month, COALESCE(SUM(total_invoice),0) AS value FROM supply_chain_invoice_datas WHERE iv_date BETWEEN ? AND ?`
	invoiceTrendParams := []interface{}{trendStart, endDate}
	if hasCompany {
		invoiceTrendQuery += " AND company = ?"
		invoiceTrendParams = append(invoiceTrendParams, companyInt)
	}
	invoiceTrendQuery += " GROUP BY month"
	var invoiceTrend []struct {
		Month string
		Value float64
	}
	database.DB.Raw(invoiceTrendQuery, invoiceTrendParams...).Scan(&invoiceTrend)
	for _, row := range invoiceTrend {
		point, exists := trendMap[row.Month]
		if !exists {
			point = &SupplyTrendPoint{Month: row.Month}
			trendMap[row.Month] = point
		}
		point.InvoiceValue = row.Value
	}

	months := make([]string, 0, len(trendMap))
	for month := range trendMap {
		months = append(months, month)
	}
	sort.Strings(months)
	for _, month := range months {
		resp.Trend = append(resp.Trend, *trendMap[month])
	}

	// Forecast data
	var forecastModel models.AIForecast
	if err := database.DB.Where("forecast_type = ?", "supplychain").Order("created_at DESC").First(&forecastModel).Error; err == nil {
		var details []supplyForecastDetail
		if forecastModel.ForecastDetails != "" {
			_ = json.Unmarshal([]byte(forecastModel.ForecastDetails), &details)
		}

		forecastSummary := &SupplyForecastSummary{
			TotalForecast:   forecastModel.ForecastedValue,
			ConfidenceLevel: forecastModel.ConfidenceLevel,
			ModelUsed:       forecastModel.ModelUsed,
			ForecastData:    make([]SupplyForecastPoint, 0, len(details)),
		}

		if len(details) > 0 {
			total := 0.0
			for _, d := range details {
				forecastSummary.ForecastData = append(forecastSummary.ForecastData, SupplyForecastPoint{
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
	if fulfillmentRate < 90 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Fulfillment rate is %.1f%%, below target", fulfillmentRate))
	}
	if pendingPO > current.POValue*0.2 && current.POValue > 0 {
		resp.Alerts = append(resp.Alerts, "Pending PO value exceeds 20% of total PO")
	}
	if avgLead > 15 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Average lead time %.1f days is high", avgLead))
	}
	if avgOnTime < 85 && len(resp.Suppliers) > 0 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Average on-time delivery %.1f%% is below threshold", avgOnTime))
	}
	if resp.Forecast != nil && current.POValue > 0 && resp.Forecast.TotalForecast > current.POValue*1.1 {
		resp.Alerts = append(resp.Alerts, "AI forecast indicates 10% increase in upcoming procurement")
	}

	insightPayload := map[string]interface{}{
		"kpis":          resp.KPIs,
		"suppliers":     resp.Suppliers,
		"pending_orders": resp.PendingOrders,
		"trend":         resp.Trend,
		"alerts":        resp.Alerts,
	}
	if resp.Forecast != nil {
		insightPayload["forecast"] = resp.Forecast
	}

	var insightResp supplyInsightResponse
	if err := callSupplyChainAI("/api/enrich/supplychain/insights", insightPayload, &insightResp); err == nil {
		if len(insightResp.ExecutiveSummary) > 0 {
			resp.ExecutiveSummary = append(resp.ExecutiveSummary, insightResp.ExecutiveSummary...)
		}
		if len(insightResp.Insights) > 0 {
			resp.Insights = append(resp.Insights, insightResp.Insights...)
		}
		if len(insightResp.RecommendedActions) > 0 {
			resp.Recommendations = append(resp.Recommendations, insightResp.RecommendedActions...)
		}
	}

	return c.JSON(http.StatusOK, resp)
}
