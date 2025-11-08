package controllers

import (
	"fmt"
	"net/http"
	"time"

	"idash-backend-api/database"

	"github.com/labstack/echo/v4"
)

type KPIValue struct {
	Label string      `json:"label"`
	Value interface{} `json:"value"`
	Change float64    `json:"change,omitempty"`
}

type DashboardResponse struct {
	KPIs       []KPIValue `json:"kpis"`
	Charts     map[string]interface{} `json:"charts"`
	AIInsights []string    `json:"ai_insights"`
	LastUpdated time.Time  `json:"last_updated"`
}

func GetExecutiveDashboard(c echo.Context) error {
	resp := DashboardResponse{
		KPIs:       make([]KPIValue, 0),
		Charts:     make(map[string]interface{}),
		AIInsights: make([]string, 0),
		LastUpdated: time.Now(),
	}

	db := database.DB

	// Revenue MTD
	type revenueResult struct {
		Value float64
	}
	var revenue revenueResult
	db.Raw(`
		SELECT COALESCE(SUM(billed),0) AS value
		FROM channelwise_monthly_report
		WHERE DATE_FORMAT(data_month, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')
	`).Scan(&revenue)
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Revenue (MTD)", Value: revenue.Value})

	// Profit margin (using budget vs actual as proxy)
	type profitResult struct {
		Value float64
	}
	var profit profitResult
	db.Raw(`
		SELECT COALESCE((SUM(budget_amount) - SUM(actual_amount)) / NULLIF(SUM(budget_amount),0) * 100, 0) AS value
		FROM budget_summaries
		WHERE DATE_FORMAT(month, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')
	`).Scan(&profit)
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Profit Margin %", Value: profit.Value})

	// Inventory value
	type inventoryResult struct {
		Value float64
	}
	var inventory inventoryResult
	db.Raw(`
		SELECT COALESCE(SUM(amount),0) AS value
		FROM inventory_raw_datas
		WHERE DATE_FORMAT(month, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')
	`).Scan(&inventory)
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Inventory Value", Value: inventory.Value})

	// Attrition (HR)
	type attritionResult struct {
		Value float64
	}
	var attrition attritionResult
	db.Raw(`
		SELECT COALESCE((SUM(resigned_employee) / NULLIF(SUM(total_active_staff + total_active_worker),0)) * 100, 0) AS value
		FROM employee_tran_overs
		JOIN employee_basic_infos ON DATE_FORMAT(employee_tran_overs.month, '%Y-%m') = DATE_FORMAT(employee_basic_infos.report_date, '%Y-%m')
		WHERE employee_tran_overs.year = YEAR(CURDATE())
	`).Scan(&attrition)
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Attrition %", Value: attrition.Value})

	resp.AIInsights = append(resp.AIInsights, generateAIInsight("sales"))
	resp.AIInsights = append(resp.AIInsights, generateAIInsight("finance"))
	resp.AIInsights = append(resp.AIInsights, generateAIInsight("inventory"))

	return c.JSON(http.StatusOK, resp)
}

func GetFinanceDashboard(c echo.Context) error {
	resp := DashboardResponse{
		KPIs:       make([]KPIValue, 0),
		Charts:     make(map[string]interface{}),
		AIInsights: make([]string, 0),
		LastUpdated: time.Now(),
	}

	db := database.DB

	// Expense MTD
	var expenses struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE(SUM(actual_amount),0) AS value
		FROM budget_summaries
		WHERE DATE_FORMAT(month, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')
	`).Scan(&expenses)
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Expenses (MTD)", Value: expenses.Value})

	// Budget variance
	var variance struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE((SUM(budget_amount) - SUM(actual_amount)),0) AS value
		FROM budget_summaries
		WHERE DATE_FORMAT(month, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')
	`).Scan(&variance)
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Budget Variance", Value: variance.Value})

	// Loan outstanding
	var loans struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE(SUM(amount),0) AS value
		FROM bank_loan_status_raw_data
		WHERE DATE_FORMAT(month, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')
	`).Scan(&loans)
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Loan Outstanding", Value: loans.Value})

	resp.AIInsights = append(resp.AIInsights, generateAIInsight("finance"))
	resp.AIInsights = append(resp.AIInsights, generateAIInsight("sales"))

	return c.JSON(http.StatusOK, resp)
}

func GetSalesDashboard(c echo.Context) error {
	resp := DashboardResponse{
		KPIs:       make([]KPIValue, 0),
		Charts:     make(map[string]interface{}),
		AIInsights: make([]string, 0),
		LastUpdated: time.Now(),
	}

	db := database.DB

	var target struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE(SUM(lifting_target),0) AS value
		FROM channelwise_monthly_report
		WHERE DATE_FORMAT(data_month, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')
	`).Scan(&target)

	var achievement struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE(SUM(billed),0) AS value
		FROM channelwise_monthly_report
		WHERE DATE_FORMAT(data_month, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')
	`).Scan(&achievement)

	completion := 0.0
	if target.Value > 0 {
		completion = (achievement.Value / target.Value) * 100
	}

	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Target", Value: target.Value})
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Achievement", Value: achievement.Value, Change: completion})

	// Best product
	var bestProduct struct {
		ProductName string
		Value      float64
	}
	db.Raw(`
		SELECT product_name, COALESCE(SUM(value),0) AS value
		FROM best_selling_products
		WHERE DATE_FORMAT(year_month, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')
		GROUP BY product_name
		ORDER BY value DESC LIMIT 1
	`).Scan(&bestProduct)
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Top Product", Value: bestProduct.ProductName})

	resp.AIInsights = append(resp.AIInsights, generateAIInsight("sales"))

	return c.JSON(http.StatusOK, resp)
}

func GetProductionDashboard(c echo.Context) error {
	resp := DashboardResponse{
		KPIs:       make([]KPIValue, 0),
		Charts:     make(map[string]interface{}),
		AIInsights: make([]string, 0),
		LastUpdated: time.Now(),
	}

	db := database.DB

	var efficiency struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE(AVG(efficiency_percentage),0) AS value
		FROM production_efficiency
		WHERE production_date BETWEEN DATE_SUB(CURDATE(), INTERVAL 30 DAY) AND CURDATE()
	`).Scan(&efficiency)

	var downtime struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE(SUM(downtime_minutes),0) AS value
		FROM machine_maintenances
		WHERE maintenance_date BETWEEN DATE_SUB(CURDATE(), INTERVAL 30 DAY) AND CURDATE()
	`).Scan(&downtime)

	var wastage struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE(AVG(wastage),0) AS value
		FROM wastage_datas
		WHERE month BETWEEN DATE_SUB(CURDATE(), INTERVAL 6 MONTH) AND CURDATE()
	`).Scan(&wastage)

	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Efficiency %", Value: efficiency.Value})
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Downtime (mins)", Value: downtime.Value})
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Avg Wastage", Value: wastage.Value})

	resp.AIInsights = append(resp.AIInsights, generateAIInsight("production"))

	return c.JSON(http.StatusOK, resp)
}

func GetHRDashboard(c echo.Context) error {
	resp := DashboardResponse{
		KPIs:       make([]KPIValue, 0),
		Charts:     make(map[string]interface{}),
		AIInsights: make([]string, 0),
		LastUpdated: time.Now(),
	}

	db := database.DB

	var headcount struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE(SUM(total_active_staff + total_active_worker),0) AS value
		FROM employee_basic_infos
		WHERE report_date = (SELECT MAX(report_date) FROM employee_basic_infos)
	`).Scan(&headcount)

	var attendance struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE(AVG(total_present/(NULLIF(total_present + total_absent + total_leave,0))),0) * 100 AS value
		FROM employee_attendances
		WHERE date BETWEEN DATE_SUB(CURDATE(), INTERVAL 30 DAY) AND CURDATE()
	`).Scan(&attendance)

	var promotions struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE(SUM(promoted_count),0) AS value
		FROM yearly_employee_promotions
		WHERE year = YEAR(CURDATE())
	`).Scan(&promotions)

	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Total Headcount", Value: headcount.Value})
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Attendance %", Value: attendance.Value})
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Promotions YTD", Value: promotions.Value})

	resp.AIInsights = append(resp.AIInsights, generateAIInsight("hr"))

	return c.JSON(http.StatusOK, resp)
}

func GetSupplyChainDashboard(c echo.Context) error {
	resp := DashboardResponse{
		KPIs:       make([]KPIValue, 0),
		Charts:     make(map[string]interface{}),
		AIInsights: make([]string, 0),
		LastUpdated: time.Now(),
	}

	db := database.DB

	var po struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE(SUM(po_amount),0) AS value
		FROM supply_chain_pos
		WHERE DATE_FORMAT(po_date, '%Y-%m') = DATE_FORMAT(CURDATE(), '%Y-%m')
	`).Scan(&po)

	var grn struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE(SUM(grn1_amount),0) AS value
		FROM supply_chain_raw_datas
		WHERE grn1_date BETWEEN DATE_SUB(CURDATE(), INTERVAL 30 DAY) AND CURDATE()
	`).Scan(&grn)

	var suppliers struct{ Value float64 }
	db.Raw(`
		SELECT COALESCE(AVG(overall_score),0) AS value
		FROM supplier_performance
		WHERE evaluation_date BETWEEN DATE_SUB(CURDATE(), INTERVAL 90 DAY) AND CURDATE()
	`).Scan(&suppliers)

	resp.KPIs = append(resp.KPIs, KPIValue{Label: "PO Value", Value: po.Value})
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "GRN Amount", Value: grn.Value})
	resp.KPIs = append(resp.KPIs, KPIValue{Label: "Supplier Score", Value: suppliers.Value})

	resp.AIInsights = append(resp.AIInsights, generateAIInsight("supplychain"))

	return c.JSON(http.StatusOK, resp)
}

func generateAIInsight(forecastType string) string {
	var forecast struct {
		Value float64
		Upper float64
		Lower float64
		ForecastDate time.Time
	}

	database.DB.Raw(`
		SELECT forecasted_value as value, upper_bound, lower_bound, forecast_date
		FROM ai_forecasts
		WHERE forecast_type = ?
		ORDER BY created_at DESC
		LIMIT 1
	`, forecastType).Scan(&forecast)

	switch forecastType {
	case "sales":
		return formatInsight("Sales forecast projects %.2f with confidence range %.2f - %.2f", forecast.Value, forecast.Lower, forecast.Upper)
	case "finance":
		return formatInsight("Financial outlook suggests variance buffer of %.2f", forecast.Value)
	case "inventory":
		return formatInsight("Inventory projection indicates level at %.2f units", forecast.Value)
	case "production":
		return formatInsight("Production capacity forecast targets %.2f with safety range %.2f - %.2f", forecast.Value, forecast.Lower, forecast.Upper)
	case "supplychain":
		return formatInsight("Supply chain risk score forecasted at %.2f", forecast.Value)
	case "hr":
		return formatInsight("Employee attrition risk forecasted at %.2f%%", forecast.Value)
	default:
		return "AI forecast data not available yet"
	}
}

func formatInsight(msg string, values ...interface{}) string {
	return fmt.Sprintf(msg, values...)
}
