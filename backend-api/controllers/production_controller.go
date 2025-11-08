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

func GetProductionSummary(c echo.Context) error {
	month := c.QueryParam("month")
	factory := c.QueryParam("factory")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type ProductionSummary struct {
		TotalProduction float64 `json:"total_production"`
		TotalWastage    float64 `json:"total_wastage"`
		TotalCost       float64 `json:"total_cost"`
		AvgGrossProfit  float64 `json:"avg_gross_profit"`
	}

	var summary ProductionSummary

	// Get production data
	query := `
		SELECT 
			COALESCE(SUM(amonthly_amount), 0) as total_production,
			COALESCE(SUM(cmonthly_amount), 0) as total_cost,
			COALESCE(AVG(amonthly_per), 0) as avg_gross_profit
		FROM production_analyses
		WHERE month >= ? AND month <= ?
	`
	params := []interface{}{startDate, endDate}

	if factory != "" {
		query += " AND factory = ?"
		params = append(params, factory)
	}

	database.DB.Raw(query, params...).Scan(&summary)

	// Get wastage data
	var wastage float64
	wastageQuery := `
		SELECT COALESCE(SUM(wastage), 0) as total_wastage
		FROM wastage_datas
		WHERE month >= ? AND month <= ?
	`
	wastageParams := []interface{}{startDate, endDate}
	if factory != "" {
		wastageQuery += " AND factory = ?"
		wastageParams = append(wastageParams, factory)
	}
	database.DB.Raw(wastageQuery, wastageParams...).Scan(&wastage)
	summary.TotalWastage = wastage

	return c.JSON(http.StatusOK, summary)
}

func GetProductionAnalysis(c echo.Context) error {
	month := c.QueryParam("month")
	factory := c.QueryParam("factory")
	summaryGroup := c.QueryParam("summary_group")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	var analyses []models.ProductionAnalysis
	query := database.DB.Where("month >= ? AND month <= ?", startDate, endDate)

	if factory != "" {
		query = query.Where("factory = ?", factory)
	}

	if summaryGroup != "" {
		query = query.Where("summary_group = ?", summaryGroup)
	}

	query.Find(&analyses)

	return c.JSON(http.StatusOK, analyses)
}

func GetProductionLastMonths(c echo.Context) error {
	month := c.QueryParam("month")
	factory := c.QueryParam("factory")
	summaryGroup := c.QueryParam("summary_group")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	firstMonth := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	lastMonth := firstMonth.AddDate(0, -2, 0)

	type LastMonthsData struct {
		Month          time.Time `json:"month"`
		Factory        string    `json:"factory"`
		SummaryGroup   string    `json:"summary_group"`
		AMonthlyAmount float64   `json:"amonthly_amount"`
		AMonthlyPer    float64   `json:"amonthly_per"`
		CAvgAmount     float64   `json:"cavg_amount"`
		CAvgPer        float64   `json:"cavg_per"`
		CMonthlyAmount float64   `json:"cmonthly_amount"`
		CMonthlyPer    float64   `json:"cmonthly_per"`
		PAvgAmount     float64   `json:"pavg_amount"`
		PAvgPer        float64   `json:"pavg_per"`
		PMonthlyAmount float64   `json:"pmonthly_amount"`
		PMonthlyPer    float64   `json:"pmonthly_per"`
		TMonthlyAmount float64   `json:"tmonthly_amount"`
		TMonthlyPer    float64   `json:"tmonthly_per"`
	}

	var data []LastMonthsData
	query := `
		SELECT 
			month,
			factory,
			summary_group,
			SUM(amonthly_amount) as amonthly_amount,
			SUM(amonthly_per) as amonthly_per,
			SUM(cavg_amount) as cavg_amount,
			SUM(cavg_per) as cavg_per,
			SUM(cmonthly_amount) as cmonthly_amount,
			SUM(cmonthly_per) as cmonthly_per,
			SUM(pavg_amount) as pavg_amount,
			SUM(pavg_per) as pavg_per,
			SUM(pmonthly_amount) as pmonthly_amount,
			SUM(pmonthly_per) as pmonthly_per,
			SUM(tmonthly_amount) as tmonthly_amount,
			SUM(tmonthly_per) as tmonthly_per
		FROM production_analyses
		WHERE month >= ? AND month <= ? AND factory = ? AND summary_group = ?
		GROUP BY month, factory, summary_group
		ORDER BY month ASC
	`
	database.DB.Raw(query, lastMonth, firstMonth, factory, summaryGroup).Scan(&data)

	return c.JSON(http.StatusOK, data)
}

func GetWastageData(c echo.Context) error {
	month := c.QueryParam("month")
	factory := c.QueryParam("factory")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	var wastageData []models.WastageData
	query := database.DB.Where("month >= ? AND month <= ?", startDate, endDate)

	if factory != "" {
		query = query.Where("factory = ?", factory)
	}

	query.Find(&wastageData)

	return c.JSON(http.StatusOK, wastageData)
}

func GetCostAnalysis(c echo.Context) error {
	month := c.QueryParam("month")
	factory := c.QueryParam("factory")
	costType := c.QueryParam("cost_type")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	var costData []models.CostAnalysis
	query := database.DB.Where("month >= ? AND month <= ?", startDate, endDate)

	if factory != "" {
		query = query.Where("factory = ?", factory)
	}

	if costType != "" {
		query = query.Where("cost_type = ?", costType)
	}

	query.Find(&costData)

	return c.JSON(http.StatusOK, costData)
}

func GetProductionTrends(c echo.Context) error {
	factory := c.QueryParam("factory")
	months := c.QueryParam("months")
	if months == "" {
		months = "12"
	}

	monthsInt, _ := strconv.Atoi(months)
	startDate := time.Now().AddDate(0, -monthsInt, 0)

	type TrendData struct {
		Month      time.Time `json:"month"`
		Production float64   `json:"production"`
		Wastage    float64   `json:"wastage"`
		Cost       float64   `json:"cost"`
	}

	var trends []TrendData
	query := `
		SELECT 
			month,
			SUM(amonthly_amount) as production,
			(SELECT COALESCE(SUM(wastage), 0) FROM wastage_datas WHERE wastage_datas.month = production_analyses.month AND wastage_datas.factory = production_analyses.factory) as wastage,
			SUM(cmonthly_amount) as cost
		FROM production_analyses
		WHERE month >= ? AND factory = ?
		GROUP BY month
		ORDER BY month ASC
	`
	database.DB.Raw(query, startDate, factory).Scan(&trends)

	return c.JSON(http.StatusOK, trends)
}

type ProductionKPI struct {
	Label  string  `json:"label"`
	Value  float64 `json:"value"`
	Change float64 `json:"change"`
}

type LinePerformance struct {
	FactoryID        uint    `json:"factory_id"`
	ProductionLineID *uint   `json:"production_line_id"`
	ActualOutput     float64 `json:"actual_output"`
	PlannedOutput    float64 `json:"planned_output"`
	Efficiency       float64 `json:"efficiency_percentage"`
	DowntimeMinutes  float64 `json:"downtime_minutes"`
	OEE              float64 `json:"oee"`
}

type WastageMetric struct {
	Factory string  `json:"factory"`
	Wastage float64 `json:"wastage"`
	Amount  float64 `json:"amount"`
	Rate    float64 `json:"rate"`
}

type MaintenanceMetric struct {
	MachineCode string    `json:"machine_code"`
	MachineName string    `json:"machine_name"`
	Downtime    float64   `json:"downtime_minutes"`
	Events      int       `json:"events"`
	Cost        float64   `json:"cost"`
	LastDate    time.Time `json:"last_date"`
}

type ProductionTrendPoint struct {
	Date     string  `json:"date"`
	Actual   float64 `json:"actual_output"`
	Planned  float64 `json:"planned_output"`
	Downtime float64 `json:"downtime_minutes"`
}

type ProductionForecastPoint struct {
	Date       string  `json:"date"`
	Forecast   float64 `json:"forecast"`
	UpperBound float64 `json:"upper_bound"`
	LowerBound float64 `json:"lower_bound"`
}

type ProductionForecastSummary struct {
	TotalForecast        float64                   `json:"total_forecast"`
	AverageDailyForecast float64                   `json:"average_daily_forecast"`
	ConfidenceLevel      float64                   `json:"confidence_level"`
	ModelUsed            string                    `json:"model_used"`
	ForecastData         []ProductionForecastPoint `json:"forecast_data"`
}

type ProductionOverviewResponse struct {
	KPIs        []ProductionKPI            `json:"kpis"`
	Lines       []LinePerformance          `json:"lines"`
	Wastage     []WastageMetric            `json:"wastage"`
	Maintenance []MaintenanceMetric        `json:"maintenance"`
	Trend       []ProductionTrendPoint     `json:"trend"`
	Forecast    *ProductionForecastSummary `json:"forecast"`
	Alerts      []string                   `json:"alerts"`
	LastUpdated time.Time                  `json:"last_updated"`
}

type productionForecastDetail struct {
	Date       string  `json:"date"`
	Forecast   float64 `json:"forecast"`
	UpperBound float64 `json:"upper_bound"`
	LowerBound float64 `json:"lower_bound"`
}

type productionSummaryRow struct {
	ActualOutput    float64
	PlannedOutput   float64
	Efficiency      float64
	DowntimeMinutes float64
	OEE             float64
}

type wastageRow struct {
	Factory string
	Wastage float64
	Amount  float64
}

type lineRow struct {
	FactoryID        uint
	ProductionLineID *uint
	ActualOutput     float64
	PlannedOutput    float64
	Efficiency       float64
	DowntimeMinutes  float64
	OEE              float64
}

type maintenanceRow struct {
	MachineCode string
	MachineName string
	Downtime    float64
	Events      int
	Cost        float64
	LastDate    time.Time
}

type trendRow struct {
	Date     string
	Actual   float64
	Planned  float64
	Downtime float64
}

func GetProductionOverview(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	factory := c.QueryParam("factory")
	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	prevStart := startDate.AddDate(0, -1, 0)
	prevEnd := endDate.AddDate(0, -1, 0)

	resp := ProductionOverviewResponse{
		KPIs:        make([]ProductionKPI, 0),
		Lines:       make([]LinePerformance, 0),
		Wastage:     make([]WastageMetric, 0),
		Maintenance: make([]MaintenanceMetric, 0),
		Trend:       make([]ProductionTrendPoint, 0),
		Alerts:      make([]string, 0),
		LastUpdated: time.Now(),
	}

	// Production summary
	var current productionSummaryRow
	query := `
		SELECT
			COALESCE(SUM(actual_output),0) AS actual_output,
			COALESCE(SUM(planned_output),0) AS planned_output,
			COALESCE(AVG(efficiency_percentage),0) AS efficiency,
			COALESCE(SUM(downtime_minutes),0) AS downtime_minutes,
			COALESCE(AVG(oee),0) AS oee
		FROM production_efficiency
		WHERE production_date BETWEEN ? AND ?`
	params := []interface{}{startDate, endDate}
	if factory != "" {
		query += " AND factory_id = ?"
		params = append(params, factory)
	}
	database.DB.Raw(query, params...).Scan(&current)

	var previous productionSummaryRow
	prevQuery := `
		SELECT
			COALESCE(SUM(actual_output),0) AS actual_output,
			COALESCE(SUM(planned_output),0) AS planned_output,
			COALESCE(AVG(efficiency_percentage),0) AS efficiency,
			COALESCE(SUM(downtime_minutes),0) AS downtime_minutes,
			COALESCE(AVG(oee),0) AS oee
		FROM production_efficiency
		WHERE production_date BETWEEN ? AND ?`
	prevParams := []interface{}{prevStart, prevEnd}
	if factory != "" {
		prevQuery += " AND factory_id = ?"
		prevParams = append(prevParams, factory)
	}
	database.DB.Raw(prevQuery, prevParams...).Scan(&previous)

	delta := func(current, prev float64) float64 {
		if prev == 0 {
			return 0
		}
		return ((current - prev) / prev) * 100
	}

	resp.KPIs = append(resp.KPIs,
		ProductionKPI{Label: "Actual Output", Value: current.ActualOutput, Change: delta(current.ActualOutput, previous.ActualOutput)},
		ProductionKPI{Label: "Planned Output", Value: current.PlannedOutput, Change: delta(current.PlannedOutput, previous.PlannedOutput)},
		ProductionKPI{Label: "Efficiency %", Value: current.Efficiency, Change: delta(current.Efficiency, previous.Efficiency)},
		ProductionKPI{Label: "OEE %", Value: current.OEE, Change: delta(current.OEE, previous.OEE)},
		ProductionKPI{Label: "Downtime (mins)", Value: current.DowntimeMinutes, Change: delta(current.DowntimeMinutes, previous.DowntimeMinutes)},
	)

	// Wastage metrics
	var wastageRows []wastageRow
	wastageQuery := `
		SELECT factory, COALESCE(SUM(wastage),0) AS wastage, COALESCE(SUM(amount),0) AS amount
		FROM wastage_datas
		WHERE month BETWEEN ? AND ?`
	wastageParams := []interface{}{startDate, endDate}
	if factory != "" {
		wastageQuery += " AND factory = ?"
		wastageParams = append(wastageParams, factory)
	}
	wastageQuery += " GROUP BY factory ORDER BY wastage DESC"
	database.DB.Raw(wastageQuery, wastageParams...).Scan(&wastageRows)

	productionByFactory := make(map[string]float64)
	factoryProductionQuery := `
		SELECT factory, COALESCE(SUM(amonthly_amount),0) AS production
		FROM production_analyses
		WHERE month BETWEEN ? AND ?
		GROUP BY factory`
	database.DB.Raw(factoryProductionQuery, startDate, endDate).Scan(&productionByFactory)

	for _, row := range wastageRows {
		production := productionByFactory[row.Factory]
		rate := 0.0
		if production > 0 {
			rate = (row.Wastage / production) * 100
		}
		resp.Wastage = append(resp.Wastage, WastageMetric{
			Factory: row.Factory,
			Wastage: row.Wastage,
			Amount:  row.Amount,
			Rate:    rate,
		})
	}

	// Line performance
	var lineRows []lineRow
	lineQuery := `
		SELECT factory_id, production_line_id,
			COALESCE(SUM(actual_output),0) AS actual_output,
			COALESCE(SUM(planned_output),0) AS planned_output,
			COALESCE(AVG(efficiency_percentage),0) AS efficiency,
			COALESCE(SUM(downtime_minutes),0) AS downtime_minutes,
			COALESCE(AVG(oee),0) AS oee
		FROM production_efficiency
		WHERE production_date BETWEEN ? AND ?`
	lineParams := []interface{}{startDate, endDate}
	if factory != "" {
		lineQuery += " AND factory_id = ?"
		lineParams = append(lineParams, factory)
	}
	lineQuery += " GROUP BY factory_id, production_line_id ORDER BY actual_output DESC LIMIT 10"
	database.DB.Raw(lineQuery, lineParams...).Scan(&lineRows)

	for _, row := range lineRows {
		resp.Lines = append(resp.Lines, LinePerformance{
			FactoryID:        row.FactoryID,
			ProductionLineID: row.ProductionLineID,
			ActualOutput:     row.ActualOutput,
			PlannedOutput:    row.PlannedOutput,
			Efficiency:       row.Efficiency,
			DowntimeMinutes:  row.DowntimeMinutes,
			OEE:              row.OEE,
		})
	}

	// Maintenance metrics
	var maintenanceRows []maintenanceRow
	maintQuery := `
		SELECT machine_code, machine_name, 
			COALESCE(SUM(downtime_minutes),0) AS downtime,
			COUNT(*) AS events,
			COALESCE(SUM(cost),0) AS cost,
			MAX(maintenance_date) AS last_date
		FROM machine_maintenances
		WHERE maintenance_date BETWEEN ? AND ?`
	maintParams := []interface{}{startDate, endDate}
	if factory != "" {
		maintQuery += " AND factory_id = ?"
		maintParams = append(maintParams, factory)
	}
	maintQuery += " GROUP BY machine_code, machine_name ORDER BY downtime DESC LIMIT 10"
	database.DB.Raw(maintQuery, maintParams...).Scan(&maintenanceRows)

	for _, row := range maintenanceRows {
		resp.Maintenance = append(resp.Maintenance, MaintenanceMetric{
			MachineCode: row.MachineCode,
			MachineName: row.MachineName,
			Downtime:    row.Downtime,
			Events:      row.Events,
			Cost:        row.Cost,
			LastDate:    row.LastDate,
		})
	}

	// Trend data (last 6 months)
	trendStart := startDate.AddDate(0, -5, 0)
	var trendRows []trendRow
	database.DB.Raw(`
		SELECT DATE_FORMAT(production_date, '%Y-%m') AS date,
			COALESCE(SUM(actual_output),0) AS actual,
			COALESCE(SUM(planned_output),0) AS planned,
			COALESCE(SUM(downtime_minutes),0) AS downtime
		FROM production_efficiency
		WHERE production_date BETWEEN ? AND ?
		GROUP BY date
		ORDER BY date
	`, trendStart, endDate).Scan(&trendRows)

	for _, row := range trendRows {
		resp.Trend = append(resp.Trend, ProductionTrendPoint{
			Date:     row.Date,
			Actual:   row.Actual,
			Planned:  row.Planned,
			Downtime: row.Downtime,
		})
	}

	// Forecast
	var forecastModel models.AIForecast
	if err := database.DB.Where("forecast_type = ?", "production").Order("created_at DESC").First(&forecastModel).Error; err == nil {
		var details []productionForecastDetail
		if forecastModel.ForecastDetails != "" {
			_ = json.Unmarshal([]byte(forecastModel.ForecastDetails), &details)
		}

		forecastSummary := &ProductionForecastSummary{
			TotalForecast:   forecastModel.ForecastedValue,
			ConfidenceLevel: forecastModel.ConfidenceLevel,
			ModelUsed:       forecastModel.ModelUsed,
			ForecastData:    make([]ProductionForecastPoint, 0, len(details)),
		}

		if len(details) > 0 {
			total := 0.0
			for _, d := range details {
				forecastSummary.ForecastData = append(forecastSummary.ForecastData, ProductionForecastPoint{
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
	if current.Efficiency < 85 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Production efficiency dropped to %.1f%%", current.Efficiency))
	}
	if current.OEE < 70 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("OEE is at %.1f%%, below optimal threshold", current.OEE))
	}
	if current.DowntimeMinutes > previous.DowntimeMinutes*1.1 && previous.DowntimeMinutes > 0 {
		resp.Alerts = append(resp.Alerts, "Downtime increased over previous period")
	}
	if len(resp.Wastage) > 0 && resp.Wastage[0].Rate > 5 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Factory %s wastage rate %.1f%% exceeds target", resp.Wastage[0].Factory, resp.Wastage[0].Rate))
	}

	return c.JSON(http.StatusOK, resp)
}
