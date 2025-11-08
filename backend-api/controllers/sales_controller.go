package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"idash-backend-api/database"
	"idash-backend-api/models"

	"github.com/labstack/echo/v4"
)

func callAISummary(endpoint string, payload interface{}, target interface{}) error {
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

	if resp.StatusCode >= 400 {
		return fmt.Errorf("ai service returned %s", resp.Status)
	}

	if target != nil {
		if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
			return err
		}
	}

	return nil
}

func GetSalesSummary(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	// Parse year and month
	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])

	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	// Use Raw SQL for aggregation
	type SalesSummary struct {
		LiftingTarget     float64 `json:"lifting_target"`
		Lifting           float64 `json:"lifting"`
		PrimaryCollection float64 `json:"primary_collection"`
		IMSTarget         float64 `json:"ims_target"`
		IMS               float64 `json:"ims"`
		MarketCollection  float64 `json:"market_collection"`
		MemoTarget        float64 `json:"memo_target"`
		MemoQty           float64 `json:"memo_qty"`
		PGTarget          float64 `json:"pg_target"`
		PGCover           float64 `json:"pg_cover"`
	}

	var result SalesSummary
	database.DB.Raw(`
		SELECT 
			COALESCE(SUM(lifting_target), 0) as lifting_target,
			COALESCE(SUM(billed), 0) as lifting,
			COALESCE(SUM(primary_collection), 0) as primary_collection,
			COALESCE(SUM(ims_target), 0) as ims_target,
			COALESCE(SUM(ims), 0) as ims,
			COALESCE(SUM(market_collection), 0) as market_collection,
			COALESCE(SUM(memo_target), 0) as memo_target,
			COALESCE(SUM(memo_qty), 0) as memo_qty,
			COALESCE(SUM(pg_target), 0) as pg_target,
			COALESCE(SUM(pg_cover), 0) as pg_cover
		FROM channelwise_monthly_report
		WHERE data_month >= ? AND data_month <= ?
	`, startDate, endDate).Scan(&result)

	// Calculate dues
	primaryDue := result.Lifting - result.PrimaryCollection
	imsDue := result.IMS - result.MarketCollection

	summary := map[string]interface{}{
		"year_month":         yearMonth,
		"lifting_target":     result.LiftingTarget,
		"lifting":            result.Lifting,
		"primary_collection": result.PrimaryCollection,
		"primary_due":        primaryDue,
		"ims_target":         result.IMSTarget,
		"ims":                result.IMS,
		"market_collection":  result.MarketCollection,
		"ims_due":            imsDue,
		"memo_target":        result.MemoTarget,
		"memo_qty":           result.MemoQty,
		"pg_target":          result.PGTarget,
		"pg_cover":           result.PGCover,
	}

	return c.JSON(http.StatusOK, summary)
}

func GetSalesCumulative(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	startDate := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(year, 12, 31, 23, 59, 59, 0, time.UTC)

	type CumulativeData struct {
		DataMonth         time.Time `json:"data_month"`
		LiftingTarget     float64   `json:"lifting_target"`
		Lifting           float64   `json:"lifting"`
		PrimaryCollection float64   `json:"primary_collection"`
		IMSTarget         float64   `json:"ims_target"`
		IMS               float64   `json:"ims"`
		MarketCollection  float64   `json:"market_collection"`
	}

	var data []CumulativeData
	database.DB.Raw(`
		SELECT 
			data_month,
			COALESCE(SUM(lifting_target), 0) as lifting_target,
			COALESCE(SUM(billed), 0) as lifting,
			COALESCE(SUM(primary_collection), 0) as primary_collection,
			COALESCE(SUM(ims_target), 0) as ims_target,
			COALESCE(SUM(ims), 0) as ims,
			COALESCE(SUM(market_collection), 0) as market_collection
		FROM channelwise_monthly_report
		WHERE data_month >= ? AND data_month <= ?
		GROUP BY data_month
		ORDER BY data_month ASC
	`, startDate, endDate).Scan(&data)

	return c.JSON(http.StatusOK, data)
}

func GetChannelwiseSales(c echo.Context) error {
	channelID := c.QueryParam("channel")
	yearMonth := c.QueryParam("yearMonth")
	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type ChannelwiseData struct {
		ChannelID         int       `json:"channel_id"`
		ChannelName       string    `json:"channel_name"`
		DataMonth         time.Time `json:"data_month"`
		LiftingTarget     float64   `json:"lifting_target"`
		Billed            float64   `json:"billed"`
		IMS               float64   `json:"ims"`
		PrimaryCollection float64   `json:"primary_collection"`
		IMSTarget         float64   `json:"ims_target"`
		MarketCollection  float64   `json:"market_collection"`
	}

	var data []ChannelwiseData
	query := `
		SELECT 
			channel_id,
			channel_name,
			data_month,
			COALESCE(SUM(lifting_target), 0) as lifting_target,
			COALESCE(SUM(billed), 0) as billed,
			COALESCE(SUM(ims), 0) as ims,
			COALESCE(SUM(primary_collection), 0) as primary_collection,
			COALESCE(SUM(ims_target), 0) as ims_target,
			COALESCE(SUM(market_collection), 0) as market_collection
		FROM channelwise_monthly_report
		WHERE data_month >= ? AND data_month <= ?
	`
	params := []interface{}{startDate, endDate}

	if channelID != "" {
		query += " AND channel_id = ?"
		chID, _ := strconv.Atoi(channelID)
		params = append(params, chID)
	}

	query += " GROUP BY channel_id, channel_name, data_month ORDER BY data_month ASC"
	database.DB.Raw(query, params...).Scan(&data)

	return c.JSON(http.StatusOK, data)
}

func GetDailySales(c echo.Context) error {
	fromDate := c.QueryParam("fromdate")
	toDate := c.QueryParam("todate")
	channelID := c.QueryParam("channel")

	if fromDate == "" {
		fromDate = time.Now().Format("2006-01-02")
	}
	if toDate == "" {
		toDate = time.Now().Format("2006-01-02")
	}

	from, _ := time.Parse("2006-01-02", fromDate)
	to, _ := time.Parse("2006-01-02", toDate)

	type DailyData struct {
		DataDate          time.Time `json:"data_date"`
		LiftingTarget     float64   `json:"lifting_target"`
		Lifting           float64   `json:"lifting"`
		PrimaryCollection float64   `json:"primary_collection"`
		IMSTarget         float64   `json:"ims_target"`
		IMS               float64   `json:"ims"`
		MarketCollection  float64   `json:"market_collection"`
	}

	var data []DailyData
	query := `
		SELECT 
			data_date,
			COALESCE(SUM(lifting_target), 0) as lifting_target,
			COALESCE(SUM(billed), 0) as lifting,
			COALESCE(SUM(lifting_collection), 0) as primary_collection,
			COALESCE(SUM(ims_target), 0) as ims_target,
			COALESCE(SUM(ims), 0) as ims,
			COALESCE(SUM(ims_collection), 0) as market_collection
		FROM channelwise_lic_data
		WHERE data_date >= ? AND data_date <= ?
	`
	params := []interface{}{from, to}

	if channelID != "" {
		query += " AND channel_id = ?"
		chID, _ := strconv.Atoi(channelID)
		params = append(params, chID)
	}

	query += " GROUP BY data_date ORDER BY data_date ASC"
	database.DB.Raw(query, params...).Scan(&data)

	return c.JSON(http.StatusOK, data)
}

func GetBestSellingProducts(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	channelID := c.QueryParam("channel")
	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "10"
	}

	limitInt, _ := strconv.Atoi(limit)

	var products []models.BestSellingProduct
	query := database.DB

	if yearMonth != "" {
		year, _ := strconv.Atoi(yearMonth[:4])
		month, _ := strconv.Atoi(yearMonth[5:7])
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)
		query = query.Where("year_month >= ? AND year_month <= ?", startDate, endDate)
	}

	if channelID != "" {
		chID, _ := strconv.Atoi(channelID)
		query = query.Where("channel_id = ?", chID)
	}

	query.Order("value DESC").Limit(limitInt).Find(&products)

	return c.JSON(http.StatusOK, products)
}

func GetBestSellingPGs(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	channelID := c.QueryParam("channel")
	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "10"
	}

	limitInt, _ := strconv.Atoi(limit)

	var pgs []models.BestSellingPg
	query := database.DB

	if yearMonth != "" {
		year, _ := strconv.Atoi(yearMonth[:4])
		month, _ := strconv.Atoi(yearMonth[5:7])
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)
		query = query.Where("year_month >= ? AND year_month <= ?", startDate, endDate)
	}

	if channelID != "" {
		chID, _ := strconv.Atoi(channelID)
		query = query.Where("channel_id = ?", chID)
	}

	query.Order("value DESC").Limit(limitInt).Find(&pgs)

	return c.JSON(http.StatusOK, pgs)
}

func GetSlowMovingProducts(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "10"
	}

	limitInt, _ := strconv.Atoi(limit)

	var products []models.BestSellingProduct
	query := database.DB

	if yearMonth != "" {
		year, _ := strconv.Atoi(yearMonth[:4])
		month, _ := strconv.Atoi(yearMonth[5:7])
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)
		query = query.Where("year_month >= ? AND year_month <= ?", startDate, endDate)
	}

	query.Order("qty ASC").Limit(limitInt).Find(&products)

	return c.JSON(http.StatusOK, products)
}

func GetTopDistributors(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	targetDate, _ := time.Parse("2006-01-02", date)

	var distributors []models.TopChannelDB
	query := database.DB.Where("date = ? AND type = ?", targetDate, 0)
	query.Order("amount DESC").Limit(10).Find(&distributors)

	return c.JSON(http.StatusOK, distributors)
}

func GetTopRetailers(c echo.Context) error {
	date := c.QueryParam("date")
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	targetDate, _ := time.Parse("2006-01-02", date)

	var retailers []models.TopRetailer
	query := database.DB.Where("date = ?", targetDate)
	query.Order("amount DESC").Limit(10).Find(&retailers)

	return c.JSON(http.StatusOK, retailers)
}

func GetOrderVsDelivery(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	var orders []models.SalesOrder
	var deliveries []models.SalesDelivery

	database.DB.Where("order_date >= ? AND order_date <= ?", startDate, endDate).Find(&orders)
	database.DB.Where("delivery_date >= ? AND delivery_date <= ?", startDate, endDate).Find(&deliveries)

	totalOrders := 0.0
	totalDeliveries := 0.0

	for _, order := range orders {
		totalOrders += order.Amount
	}

	for _, delivery := range deliveries {
		totalDeliveries += delivery.Amount
	}

	result := map[string]interface{}{
		"year_month":       yearMonth,
		"total_orders":     totalOrders,
		"total_deliveries": totalDeliveries,
		"pending":          totalOrders - totalDeliveries,
	}

	return c.JSON(http.StatusOK, result)
}

type SalesKPI struct {
	Label  string  `json:"label"`
	Value  float64 `json:"value"`
	Change float64 `json:"change"`
}

type ChannelPerformance struct {
	ChannelID    int     `json:"channel_id"`
	ChannelName  string  `json:"channel_name"`
	Billed       float64 `json:"billed"`
	Target       float64 `json:"target"`
	Primary      float64 `json:"primary_collection"`
	IMS          float64 `json:"ims"`
	Achievement  float64 `json:"achievement"`
	Contribution float64 `json:"contribution"`
}

type ProductMetric struct {
	Name      string  `json:"name"`
	Value     float64 `json:"value"`
	Quantity  float64 `json:"quantity"`
	ChannelID int     `json:"channel_id"`
}

type ProductPerformance struct {
	TopProducts       []ProductMetric `json:"top_products"`
	SlowMovers        []ProductMetric `json:"slow_movers"`
	BestProductGroups []ProductMetric `json:"best_product_groups"`
}

type SalesTrendPoint struct {
	Month  string  `json:"month"`
	Billed float64 `json:"billed"`
	Target float64 `json:"target"`
}

type ForecastPoint struct {
	Date       string  `json:"date"`
	Forecast   float64 `json:"forecast"`
	UpperBound float64 `json:"upper_bound"`
	LowerBound float64 `json:"lower_bound"`
}

type SalesForecastSummary struct {
	TotalForecast        float64         `json:"total_forecast"`
	AverageDailyForecast float64         `json:"average_daily_forecast"`
	ConfidenceLevel      float64         `json:"confidence_level"`
	ModelUsed            string          `json:"model_used"`
	ForecastData         []ForecastPoint `json:"forecast_data"`
}

type SalesOverviewResponse struct {
	KPIs            []SalesKPI            `json:"kpis"`
	Channels        []ChannelPerformance  `json:"channels"`
	Products        ProductPerformance    `json:"products"`
	Trend           []SalesTrendPoint     `json:"trend"`
	Forecast        *SalesForecastSummary `json:"forecast"`
	Alerts          []string              `json:"alerts"`
	Predictions     []map[string]interface{} `json:"predictions,omitempty"`
	Recommendations []map[string]interface{} `json:"recommendations,omitempty"`
	Anomalies       map[string]interface{}   `json:"anomalies,omitempty"`
	Scenario        map[string]interface{}   `json:"scenario,omitempty"`
	LastUpdated     time.Time             `json:"last_updated"`
}

type forecastDetail struct {
	Date       string  `json:"date"`
	Forecast   float64 `json:"forecast"`
	UpperBound float64 `json:"upper_bound"`
	LowerBound float64 `json:"lower_bound"`
}

func GetSalesOverview(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	prevStart := startDate.AddDate(-1, 0, 0)
	prevEnd := endDate.AddDate(-1, 0, 0)

	resp := SalesOverviewResponse{
		KPIs:     make([]SalesKPI, 0),
		Channels: make([]ChannelPerformance, 0),
		Products: ProductPerformance{
			TopProducts:       make([]ProductMetric, 0),
			SlowMovers:        make([]ProductMetric, 0),
			BestProductGroups: make([]ProductMetric, 0),
		},
		Trend:       make([]SalesTrendPoint, 0),
		Alerts:      make([]string, 0),
		LastUpdated: time.Now(),
	}

	// Current summary
	type summaryRow struct {
		Target            float64
		Billed            float64
		PrimaryCollection float64
		IMS               float64
		IMSTarget         float64
		MarketCollection  float64
	}

	var current summaryRow
	database.DB.Raw(`
		SELECT 
			COALESCE(SUM(lifting_target),0) AS target,
			COALESCE(SUM(billed),0) AS billed,
			COALESCE(SUM(primary_collection),0) AS primary_collection,
			COALESCE(SUM(ims),0) AS ims,
			COALESCE(SUM(ims_target),0) AS ims_target,
			COALESCE(SUM(market_collection),0) AS market_collection
		FROM channelwise_monthly_report
		WHERE data_month BETWEEN ? AND ?
	`, startDate, endDate).Scan(&current)

	var previous summaryRow
	database.DB.Raw(`
		SELECT 
			COALESCE(SUM(lifting_target),0) AS target,
			COALESCE(SUM(billed),0) AS billed,
			COALESCE(SUM(primary_collection),0) AS primary_collection,
			COALESCE(SUM(ims),0) AS ims,
			COALESCE(SUM(ims_target),0) AS ims_target,
			COALESCE(SUM(market_collection),0) AS market_collection
		FROM channelwise_monthly_report
		WHERE data_month BETWEEN ? AND ?
	`, prevStart, prevEnd).Scan(&previous)

	achievement := 0.0
	if current.Target > 0 {
		achievement = (current.Billed / current.Target) * 100
	}
	yoy := 0.0
	if previous.Billed > 0 {
		yoy = ((current.Billed - previous.Billed) / previous.Billed) * 100
	}
	primaryDue := current.Billed - current.PrimaryCollection
	imsDue := current.IMS - current.MarketCollection

	resp.KPIs = append(resp.KPIs,
		SalesKPI{Label: "Monthly Sales", Value: current.Billed, Change: yoy},
		SalesKPI{Label: "Target Achievement %", Value: achievement, Change: 0},
		SalesKPI{Label: "Primary Collection", Value: current.PrimaryCollection, Change: 0},
		SalesKPI{Label: "Primary Due", Value: primaryDue, Change: 0},
		SalesKPI{Label: "IMS", Value: current.IMS, Change: 0},
		SalesKPI{Label: "IMS Due", Value: imsDue, Change: 0},
		SalesKPI{Label: "Market Collection", Value: current.MarketCollection, Change: 0},
	)

	// Channel performance
	type channelRow struct {
		ChannelID   int
		ChannelName string
		Billed      float64
		Target      float64
		Primary     float64
		IMS         float64
	}

	var channels []channelRow
	database.DB.Raw(`
		SELECT channel_id, channel_name,
			COALESCE(SUM(billed),0) AS billed,
			COALESCE(SUM(lifting_target),0) AS target,
			COALESCE(SUM(primary_collection),0) AS primary,
			COALESCE(SUM(ims),0) AS ims
		FROM channelwise_monthly_report
		WHERE data_month BETWEEN ? AND ?
		GROUP BY channel_id, channel_name
		ORDER BY billed DESC
	`, startDate, endDate).Scan(&channels)

	totalBilled := current.Billed
	for _, ch := range channels {
		achievement := 0.0
		if ch.Target > 0 {
			achievement = (ch.Billed / ch.Target) * 100
		}
		contribution := 0.0
		if totalBilled > 0 {
			contribution = (ch.Billed / totalBilled) * 100
		}
		resp.Channels = append(resp.Channels, ChannelPerformance{
			ChannelID:    ch.ChannelID,
			ChannelName:  ch.ChannelName,
			Billed:       ch.Billed,
			Target:       ch.Target,
			Primary:      ch.Primary,
			IMS:          ch.IMS,
			Achievement:  achievement,
			Contribution: contribution,
		})
	}

	// Product performance
	var topProducts []models.BestSellingProduct
	database.DB.Where("year_month BETWEEN ? AND ?", startDate, endDate).
		Order("value DESC").Limit(10).Find(&topProducts)
	for _, p := range topProducts {
		resp.Products.TopProducts = append(resp.Products.TopProducts, ProductMetric{
			Name:      p.ProductName,
			Value:     p.Value,
			Quantity:  p.Qty,
			ChannelID: p.ChannelID,
		})
	}

	var slowProducts []models.BestSellingProduct
	database.DB.Where("year_month BETWEEN ? AND ?", startDate, endDate).
		Order("value ASC").Limit(10).Find(&slowProducts)
	for _, p := range slowProducts {
		resp.Products.SlowMovers = append(resp.Products.SlowMovers, ProductMetric{
			Name:      p.ProductName,
			Value:     p.Value,
			Quantity:  p.Qty,
			ChannelID: p.ChannelID,
		})
	}

	var bestPGs []models.BestSellingPg
	database.DB.Where("year_month BETWEEN ? AND ?", startDate, endDate).
		Order("value DESC").Limit(10).Find(&bestPGs)
	for _, pg := range bestPGs {
		resp.Products.BestProductGroups = append(resp.Products.BestProductGroups, ProductMetric{
			Name:     pg.CategoryName,
			Value:    pg.Value,
			Quantity: pg.Qty,
		})
	}

	// Trend data (last 12 months)
	type trendRow struct {
		Month  string
		Billed float64
		Target float64
	}
	var trendRows []trendRow
	trendStart := startDate.AddDate(0, -11, 0)
	database.DB.Raw(`
		SELECT DATE_FORMAT(data_month, '%Y-%m') AS month,
			COALESCE(SUM(billed),0) AS billed,
			COALESCE(SUM(lifting_target),0) AS target
		FROM channelwise_monthly_report
		WHERE data_month BETWEEN ? AND ?
		GROUP BY month
		ORDER BY month
	`, trendStart, endDate).Scan(&trendRows)

	for _, t := range trendRows {
		resp.Trend = append(resp.Trend, SalesTrendPoint{
			Month:  t.Month,
			Billed: t.Billed,
			Target: t.Target,
		})
	}

	// Forecast data from AI forecasts table
	var forecastModel models.AIForecast
	if err := database.DB.Where("forecast_type = ?", "sales").Order("created_at DESC").First(&forecastModel).Error; err == nil {
		var details []forecastDetail
		if forecastModel.ForecastDetails != "" {
			_ = json.Unmarshal([]byte(forecastModel.ForecastDetails), &details)
		}

		forecastSummary := &SalesForecastSummary{
			ForecastType:       forecastModel.ForecastType,
			EntityID:           forecastModel.EntityID,
			ForecastPeriodDays: 30,
			ConfidenceLevel:    forecastModel.ConfidenceLevel,
			TotalForecast:      forecastModel.ForecastedValue,
			ModelUsed:          forecastModel.ModelUsed,
			ForecastData:       make([]ForecastPoint, 0),
		}

		for _, d := range details {
			forecastSummary.ForecastData = append(forecastSummary.ForecastData, ForecastPoint{
				Date:       d.Date,
				Forecast:   d.Forecast,
				UpperBound: d.UpperBound,
				LowerBound: d.LowerBound,
			})
		}

		resp.Forecast = forecastSummary
	}

	// Predictive summary & prescriptive recommendations
	predictionPayload := map[string]interface{}{
		"start_date": trendStart.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"horizon":    30,
		"granularity": "channel",
	}
	var predictionResp SalesPredictionResponse
	if err := callAISummary("/api/predict/sales/summary", predictionPayload, &predictionResp); err == nil {
		resp.Predictions = predictionResp.Predictions
		resp.Recommendations = predictionResp.Recommendations
	}

	// Enriched anomaly detection
	anomalyPayload := map[string]interface{}{
		"metric":       "sales",
		"analysis_type": "anomaly",
		"start_date":   trendStart.Format("2006-01-02"),
		"end_date":     endDate.Format("2006-01-02"),
	}
	var anomalyResp map[string]interface{}
	if err := callAISummary("/api/analyze/anomaly/enriched", anomalyPayload, &anomalyResp); err == nil {
		resp.Anomalies = anomalyResp
	}

	// Default what-if scenario (price +2%, production +3%)
	scenarioPayload := map[string]interface{}{
		"horizon": 90,
		"base_metrics": map[string]float64{
			"sales":        current.Billed,
			"gross_margin": current.PrimaryCollection,
		},
		"adjustments": map[string]float64{
			"price_change_pct":   2,
			"volume_change_pct":  1,
			"cost_change_pct":   -0.5,
		},
	}
	var scenarioResp ScenarioSimulationResponse
	if err := callAISummary("/api/scenario/whatif", scenarioPayload, &scenarioResp); err == nil {
		resp.Scenario = map[string]interface{}{
			"projected_sales":    scenarioResp.ProjectedSales,
			"projected_margin":   scenarioResp.ProjectedMargin,
			"incremental_profit": scenarioResp.IncrementalProfit,
			"narrative":          scenarioResp.Narrative,
		}
	}

	// Alerts
	if achievement < 95 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Target achievement is %.1f%%, below expected threshold", achievement))
	}
	if yoy < 0 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Year-over-year sales decreased by %.1f%%", yoy))
	}
	if primaryDue > current.Billed*0.15 {
		resp.Alerts = append(resp.Alerts, "Primary collection gap exceeds 15% of billed amount")
	}
	if resp.Forecast != nil && resp.Forecast.TotalForecast < current.Billed {
		resp.Alerts = append(resp.Alerts, "AI forecast indicates potential slowdown versus current month performance")
	}

	return c.JSON(http.StatusOK, resp)
}
