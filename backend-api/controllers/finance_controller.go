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

func callFinanceAI(endpoint string, payload interface{}, target interface{}) error {
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

func GetFinanceSummary(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	companyCode := c.QueryParam("company_code")

	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type FinanceSummary struct {
		TotalBudget     float64 `json:"total_budget"`
		TotalExpense    float64 `json:"total_expense"`
		TotalBankLoan   float64 `json:"total_bank_loan"`
		NetSales        float64 `json:"net_sales"`
		FinancialExpense float64 `json:"financial_expense"`
		BudgetVsActual  float64 `json:"budget_vs_actual"`
	}

	var summary FinanceSummary

	// Get budget data
	budgetQuery := `
		SELECT COALESCE(SUM(amount), 0) as total_budget
		FROM bdgt_budgets
		WHERE month >= ? AND month <= ?
	`
	budgetParams := []interface{}{startDate, endDate}
	database.DB.Raw(budgetQuery, budgetParams...).Scan(&summary.TotalBudget)

	// Get expense data
	expenseQuery := `
		SELECT COALESCE(SUM(amount), 0) as total_expense
		FROM bdgt_expenses
		WHERE month >= ? AND month <= ?
	`
	expenseParams := []interface{}{startDate, endDate}
	database.DB.Raw(expenseQuery, expenseParams...).Scan(&summary.TotalExpense)

	// Get bank loan data
	bankLoanQuery := `
		SELECT COALESCE(SUM(amount), 0) as total_bank_loan
		FROM bank_loan_status_raw_data
		WHERE month = ?
	`
	bankLoanParams := []interface{}{endDate}
	if companyCode != "" {
		bankLoanQuery += " AND company_id = ?"
		companyID, _ := strconv.Atoi(companyCode)
		bankLoanParams = append(bankLoanParams, companyID)
	}
	database.DB.Raw(bankLoanQuery, bankLoanParams...).Scan(&summary.TotalBankLoan)

	// Get net sales and financial expense
	salesQuery := `
		SELECT 
			COALESCE(SUM(net_sales), 0) as net_sales,
			COALESCE(SUM(financial_expense), 0) as financial_expense
		FROM bank_loan
		WHERE year_month >= ? AND year_month <= ?
	`
	salesParams := []interface{}{startDate, endDate}
	if companyCode != "" {
		salesQuery += " AND company_code = ?"
		companyID, _ := strconv.Atoi(companyCode)
		salesParams = append(salesParams, companyID)
	}
	database.DB.Raw(salesQuery, salesParams...).Scan(&summary)

	summary.BudgetVsActual = summary.TotalBudget - summary.TotalExpense

	return c.JSON(http.StatusOK, summary)
}

func GetBudget(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	categoryID := c.QueryParam("category_id")
	departmentID := c.QueryParam("department_id")

	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	var budgets []models.BdgtBudget
	query := database.DB.Where("month >= ? AND month <= ?", startDate, endDate)

	if categoryID != "" {
		catID, _ := strconv.Atoi(categoryID)
		query = query.Where("category_id = ?", catID)
	}

	if departmentID != "" {
		deptID, _ := strconv.Atoi(departmentID)
		query = query.Where("department_id = ?", deptID)
	}

	query.Find(&budgets)

	return c.JSON(http.StatusOK, budgets)
}

func GetBudgetVsActual(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	categoryID := c.QueryParam("category_id")

	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type BudgetVsActual struct {
		CategoryID   int     `json:"category_id"`
		CategoryName string  `json:"category_name"`
		Budget       float64 `json:"budget"`
		Actual       float64 `json:"actual"`
		Variance     float64 `json:"variance"`
		VariancePercent float64 `json:"variance_percent"`
	}

	var data []BudgetVsActual
	query := `
		SELECT 
			b.category_id,
			c.category_name,
			COALESCE(SUM(b.amount), 0) as budget,
			COALESCE(SUM(e.amount), 0) as actual,
			COALESCE(SUM(b.amount), 0) - COALESCE(SUM(e.amount), 0) as variance,
			CASE 
				WHEN COALESCE(SUM(b.amount), 0) > 0 
				THEN ((COALESCE(SUM(b.amount), 0) - COALESCE(SUM(e.amount), 0)) / COALESCE(SUM(b.amount), 0)) * 100
				ELSE 0
			END as variance_percent
		FROM bdgt_budgets b
		LEFT JOIN bdgt_categories c ON b.category_id = c.id
		LEFT JOIN bdgt_expenses e ON b.id = e.budget_id AND e.month >= ? AND e.month <= ?
		WHERE b.month >= ? AND b.month <= ?
	`
	params := []interface{}{startDate, endDate, startDate, endDate}

	if categoryID != "" {
		query += " AND b.category_id = ?"
		catID, _ := strconv.Atoi(categoryID)
		params = append(params, catID)
	}

	query += " GROUP BY b.category_id, c.category_name"
	database.DB.Raw(query, params...).Scan(&data)

	return c.JSON(http.StatusOK, data)
}

func GetBankLoan(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	companyCode := c.QueryParam("company_code")

	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	var bankLoans []models.BankLoan
	query := database.DB.Where("year_month >= ? AND year_month <= ?", startDate, endDate)

	if companyCode != "" {
		companyID, _ := strconv.Atoi(companyCode)
		query = query.Where("company_code = ?", companyID)
	}

	query.Find(&bankLoans)

	return c.JSON(http.StatusOK, bankLoans)
}

func GetBankLoanStatus(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	companyCode := c.QueryParam("company_code")

	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	targetDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	type BankLoanStatus struct {
		Head   string  `json:"head"`
		Amount float64 `json:"amount"`
	}

	var status []BankLoanStatus
	query := `
		SELECT 
			head,
			COALESCE(SUM(amount), 0) as amount
		FROM bank_loan_status_raw_data
		WHERE month = ?
	`
	params := []interface{}{targetDate}

	if companyCode != "" {
		query += " AND company_id = ?"
		companyID, _ := strconv.Atoi(companyCode)
		params = append(params, companyID)
	}

	query += " GROUP BY head ORDER BY amount DESC"
	database.DB.Raw(query, params...).Scan(&status)

	return c.JSON(http.StatusOK, status)
}

func GetCompanyBankLoanSummary(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")

	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	targetDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)

	type CompanyBankLoan struct {
		CompanyID   int     `json:"company_id"`
		CompanyName string  `json:"company_name"`
		TotalLoan   float64 `json:"total_loan"`
	}

	var summary []CompanyBankLoan
	query := `
		SELECT 
			bl.company_id,
			c.name as company_name,
			COALESCE(SUM(bl.amount), 0) as total_loan
		FROM bank_loan_status_raw_data bl
		LEFT JOIN tbl_company c ON bl.company_id = c.id
		WHERE bl.month = ?
		GROUP BY bl.company_id, c.name
		ORDER BY total_loan DESC
	`
	database.DB.Raw(query, targetDate).Scan(&summary)

	return c.JSON(http.StatusOK, summary)
}

func GetExpenses(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	budgetID := c.QueryParam("budget_id")

	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	var expenses []models.BdgtExpense
	query := database.DB.Where("month >= ? AND month <= ?", startDate, endDate)

	if budgetID != "" {
		budgetIDInt, _ := strconv.Atoi(budgetID)
		query = query.Where("budget_id = ?", budgetIDInt)
	}

	query.Find(&expenses)

	return c.JSON(http.StatusOK, expenses)
}

func GetBudgetCategories(c echo.Context) error {
	var categories []models.BdgtCategory
	database.DB.Where("status = ?", 1).Find(&categories)
	return c.JSON(http.StatusOK, categories)
}

func GetBudgetDepartments(c echo.Context) error {
	var departments []models.BdgtDepartment
	database.DB.Where("status = ?", 1).Find(&departments)
	return c.JSON(http.StatusOK, departments)
}

type FinanceKPI struct {
	Label  string  `json:"label"`
	Value  float64 `json:"value"`
	Change float64 `json:"change"`
}

type DepartmentPerformance struct {
	DepartmentID   int     `json:"department_id"`
	DepartmentName string  `json:"department_name"`
	Budget         float64 `json:"budget"`
	Actual         float64 `json:"actual"`
	Variance       float64 `json:"variance"`
	VariancePct    float64 `json:"variance_percent"`
}

type CategoryExpense struct {
	CategoryID   int     `json:"category_id"`
	CategoryName string  `json:"category_name"`
	Actual       float64 `json:"actual"`
}

type LoanExposure struct {
	Head   string  `json:"head"`
	Amount float64 `json:"amount"`
}

type FinanceTrendPoint struct {
	Month  string  `json:"month"`
	Budget float64 `json:"budget"`
	Actual float64 `json:"actual"`
}

type FinanceForecastPoint struct {
	Date       string  `json:"date"`
	Forecast   float64 `json:"forecast"`
	UpperBound float64 `json:"upper_bound"`
	LowerBound float64 `json:"lower_bound"`
}

type FinanceForecastSummary struct {
	TotalForecast        float64               `json:"total_forecast"`
	AverageDailyForecast float64               `json:"average_daily_forecast"`
	ConfidenceLevel      float64               `json:"confidence_level"`
	ModelUsed            string                `json:"model_used"`
	ForecastData         []FinanceForecastPoint `json:"forecast_data"`
}

type FinanceOverviewResponse struct {
	KPIs        []FinanceKPI             `json:"kpis"`
	Departments []DepartmentPerformance  `json:"departments"`
	Categories  []CategoryExpense        `json:"categories"`
	Loans       []LoanExposure           `json:"loans"`
	Trend       []FinanceTrendPoint      `json:"trend"`
	Forecast    *FinanceForecastSummary  `json:"forecast"`
	Alerts      []string                 `json:"alerts"`
	Prescriptions map[string]interface{} `json:"prescriptions,omitempty"`
	Scenario    map[string]interface{}   `json:"scenario,omitempty"`
	LastUpdated time.Time                `json:"last_updated"`
}

type financeSummaryRow struct {
	Budget float64
	Actual float64
}

type financeForecastDetail struct {
	Date       string  `json:"date"`
	Forecast   float64 `json:"forecast"`
	UpperBound float64 `json:"upper_bound"`
	LowerBound float64 `json:"lower_bound"`
}

func GetFinanceOverview(c echo.Context) error {
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

	resp := FinanceOverviewResponse{
		KPIs:        make([]FinanceKPI, 0),
		Departments: make([]DepartmentPerformance, 0),
		Categories:  make([]CategoryExpense, 0),
		Loans:       make([]LoanExposure, 0),
		Trend:       make([]FinanceTrendPoint, 0),
		Alerts:      make([]string, 0),
		LastUpdated: time.Now(),
	}

	// Budget vs actual summary
	var current financeSummaryRow
	database.DB.Raw(`
		SELECT COALESCE(SUM(budget_amount),0) AS budget, COALESCE(SUM(actual_amount),0) AS actual
		FROM budget_summaries
		WHERE month BETWEEN ? AND ?
	`, startDate, endDate).Scan(&current)

	var previous financeSummaryRow
	database.DB.Raw(`
		SELECT COALESCE(SUM(budget_amount),0) AS budget, COALESCE(SUM(actual_amount),0) AS actual
		FROM budget_summaries
		WHERE month BETWEEN ? AND ?
	`, prevStart, prevEnd).Scan(&previous)

	var loans struct{ Value float64 }
	database.DB.Raw(`
		SELECT COALESCE(SUM(amount),0) AS value
		FROM bank_loan_status_raw_data
		WHERE month = ?
	`, endDate).Scan(&loans)

	var revenue struct{ Value float64 }
	database.DB.Raw(`
		SELECT COALESCE(SUM(net_sales),0) AS value
		FROM bank_loan
		WHERE year_month BETWEEN ? AND ?
	`, startDate, endDate).Scan(&revenue)

	delta := func(cur, prev float64) float64 {
		if prev == 0 {
			return 0
		}
		return ((cur - prev) / prev) * 100
	}

	variance := current.Budget - current.Actual

	resp.KPIs = append(resp.KPIs,
		FinanceKPI{Label: "Budget", Value: current.Budget, Change: delta(current.Budget, previous.Budget)},
		FinanceKPI{Label: "Actual Spend", Value: current.Actual, Change: delta(current.Actual, previous.Actual)},
		FinanceKPI{Label: "Variance", Value: variance, Change: delta(variance, previous.Budget-previous.Actual)},
		FinanceKPI{Label: "Loan Exposure", Value: loans.Value, Change: 0},
		FinanceKPI{Label: "Net Sales", Value: revenue.Value, Change: 0},
	)

	// Department performance
	var departmentRows []struct {
		DepartmentID   int
		DepartmentName string
		Budget         float64
		Actual         float64
	}
	database.DB.Raw(`
		SELECT bs.department_id, COALESCE(d.name, 'Unassigned') AS department_name,
			COALESCE(SUM(bs.budget_amount),0) AS budget,
			COALESCE(SUM(bs.actual_amount),0) AS actual
		FROM budget_summaries bs
		LEFT JOIN bdgt_departments d ON bs.department_id = d.id
		WHERE bs.month BETWEEN ? AND ?
		GROUP BY bs.department_id, d.name
		ORDER BY actual DESC
		LIMIT 10
	`, startDate, endDate).Scan(&departmentRows)

	for _, row := range departmentRows {
		variance := row.Budget - row.Actual
		variancePct := 0.0
		if row.Budget != 0 {
			variancePct = (variance / row.Budget) * 100
		}
		resp.Departments = append(resp.Departments, DepartmentPerformance{
			DepartmentID:   row.DepartmentID,
			DepartmentName: row.DepartmentName,
			Budget:         row.Budget,
			Actual:         row.Actual,
			Variance:       variance,
			VariancePct:    variancePct,
		})
	}

	// Category expenses
	var categoryRows []CategoryExpense
	database.DB.Raw(`
		SELECT bs.category_id AS category_id, COALESCE(c.name, 'Unassigned') AS category_name,
			COALESCE(SUM(bs.actual_amount),0) AS actual
		FROM budget_summaries bs
		LEFT JOIN bdgt_categories c ON bs.category_id = c.id
		WHERE bs.month BETWEEN ? AND ?
		GROUP BY bs.category_id, c.name
		ORDER BY actual DESC
		LIMIT 10
	`, startDate, endDate).Scan(&categoryRows)
	resp.Categories = categoryRows

	// Loan exposure
	database.DB.Raw(`
		SELECT head, COALESCE(SUM(amount),0) AS amount
		FROM bank_loan_status_raw_data
		WHERE month = ?
		GROUP BY head
		ORDER BY amount DESC
		LIMIT 10
	`, endDate).Scan(&resp.Loans)

	// Trend over last 12 months
	trendStart := startDate.AddDate(0, -11, 0)
	var trendRows []struct {
		Month  string
		Budget float64
		Actual float64
	}
	database.DB.Raw(`
		SELECT DATE_FORMAT(month, '%Y-%m') AS month,
			COALESCE(SUM(budget_amount),0) AS budget,
			COALESCE(SUM(actual_amount),0) AS actual
		FROM budget_summaries
		WHERE month BETWEEN ? AND ?
		GROUP BY month
		ORDER BY month
	`, trendStart, endDate).Scan(&trendRows)

	for _, row := range trendRows {
		resp.Trend = append(resp.Trend, FinanceTrendPoint{
			Month:  row.Month,
			Budget: row.Budget,
			Actual: row.Actual,
		})
	}

	// Forecast data
	var forecastModel models.AIForecast
	if err := database.DB.Where("forecast_type = ?", "finance").Order("created_at DESC").First(&forecastModel).Error; err == nil {
		var details []financeForecastDetail
		if forecastModel.ForecastDetails != "" {
			_ = json.Unmarshal([]byte(forecastModel.ForecastDetails), &details)
		}

		forecastSummary := &FinanceForecastSummary{
			TotalForecast:   forecastModel.ForecastedValue,
			ConfidenceLevel: forecastModel.ConfidenceLevel,
			ModelUsed:       forecastModel.ModelUsed,
			ForecastData:    make([]FinanceForecastPoint, 0, len(details)),
		}

		if len(details) > 0 {
			total := 0.0
			for _, d := range details {
				forecastSummary.ForecastData = append(forecastSummary.ForecastData, FinanceForecastPoint{
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
	if variance < 0 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Budget overrun of ৳ %.0f", -variance))
	}
	if current.Actual > current.Budget*1.1 && current.Budget > 0 {
		resp.Alerts = append(resp.Alerts, "Actual spend exceeds budget by more than 10%")
	}
	if len(resp.Departments) > 0 && resp.Departments[0].VariancePct < -15 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Department %s variance %.1f%%", resp.Departments[0].DepartmentName, resp.Departments[0].VariancePct))
	}
	if resp.Forecast != nil && resp.Forecast.TotalForecast > current.Actual*1.1 {
		resp.Alerts = append(resp.Alerts, "AI forecast indicates higher upcoming expenses")
	}
	if loans.Value > current.Actual*0.6 && current.Actual > 0 {
		resp.Alerts = append(resp.Alerts, "Loan exposure exceeds 60% of monthly spend")
	}

	financePayload := map[string]interface{}{
		"module":     "finance",
		"start_date": trendStart.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"horizon":    90,
	}
	var financeAI FinancialPrescriptionResponse
	if err := callFinanceAI("/api/prescribe/finance", financePayload, &financeAI); err == nil {
		resp.Prescriptions = financeAI.Prescriptions
	}

	scenarioPayload := map[string]interface{}{
		"horizon": 120,
		"base_metrics": map[string]float64{
			"sales":        revenue.Value,
			"gross_margin": current.Budget - current.Actual,
		},
		"adjustments": map[string]float64{
			"price_change_pct":   1.5,
			"volume_change_pct":  0,
			"cost_change_pct":   -1,
		},
	}
	var scenarioResp ScenarioSimulationResponse
	if err := callFinanceAI("/api/scenario/whatif", scenarioPayload, &scenarioResp); err == nil {
		resp.Scenario = map[string]interface{}{
			"projected_sales":    scenarioResp.ProjectedSales,
			"projected_margin":   scenarioResp.ProjectedMargin,
			"incremental_profit": scenarioResp.IncrementalProfit,
			"narrative":          scenarioResp.Narrative,
		}
	}

	return c.JSON(http.StatusOK, resp)
}
