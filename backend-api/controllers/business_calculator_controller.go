package controllers

import (
	"idash-backend-api/database"
	"idash-backend-api/models"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func GetBusinessCalculation(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	companyCode := c.QueryParam("company_code")

	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type BusinessCalculation struct {
		NetSales         float64 `json:"net_sales"`
		BankLoan         float64 `json:"bank_loan"`
		FinancialExpense float64 `json:"financial_expense"`
		LoanRatio        float64 `json:"loan_ratio"`
		ExpenseRatio     float64 `json:"expense_ratio"`
	}

	var calc BusinessCalculation

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
	database.DB.Raw(salesQuery, salesParams...).Scan(&calc)

	// Get bank loan
	loanQuery := `
		SELECT COALESCE(SUM(amount), 0) as bank_loan
		FROM bank_loan_status_raw_data
		WHERE month = ?
	`
	loanParams := []interface{}{endDate}
	if companyCode != "" {
		loanQuery += " AND company_id = ?"
		companyID, _ := strconv.Atoi(companyCode)
		loanParams = append(loanParams, companyID)
	}
	database.DB.Raw(loanQuery, loanParams...).Scan(&calc.BankLoan)

	// Calculate ratios
	if calc.NetSales > 0 {
		calc.LoanRatio = (calc.BankLoan / calc.NetSales) * 100
		calc.ExpenseRatio = (calc.FinancialExpense / calc.NetSales) * 100
	}

	return c.JSON(http.StatusOK, calc)
}

func GetPrimaryBusinessCalculation(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")

	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type PrimaryBusinessCalc struct {
		Lifting           float64 `json:"lifting"`
		PrimaryCollection float64 `json:"primary_collection"`
		CollectionRate    float64 `json:"collection_rate"`
		DueAmount         float64 `json:"due_amount"`
	}

	var calc PrimaryBusinessCalc

	// Get lifting and primary collection
	query := `
		SELECT 
			COALESCE(SUM(billed), 0) as lifting,
			COALESCE(SUM(primary_collection), 0) as primary_collection
		FROM channelwise_monthly_report
		WHERE data_month >= ? AND data_month <= ?
	`
	params := []interface{}{startDate, endDate}
	database.DB.Raw(query, params...).Scan(&calc)

	calc.DueAmount = calc.Lifting - calc.PrimaryCollection
	if calc.Lifting > 0 {
		calc.CollectionRate = (calc.PrimaryCollection / calc.Lifting) * 100
	}

	return c.JSON(http.StatusOK, calc)
}

func GetSecondaryBusinessCalculation(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")

	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type SecondaryBusinessCalc struct {
		IMS              float64 `json:"ims"`
		MarketCollection float64 `json:"market_collection"`
		CollectionRate   float64 `json:"collection_rate"`
		DueAmount        float64 `json:"due_amount"`
	}

	var calc SecondaryBusinessCalc

	// Get IMS and market collection
	query := `
		SELECT 
			COALESCE(SUM(ims), 0) as ims,
			COALESCE(SUM(market_collection), 0) as market_collection
		FROM channelwise_monthly_report
		WHERE data_month >= ? AND data_month <= ?
	`
	params := []interface{}{startDate, endDate}
	database.DB.Raw(query, params...).Scan(&calc)

	calc.DueAmount = calc.IMS - calc.MarketCollection
	if calc.IMS > 0 {
		calc.CollectionRate = (calc.MarketCollection / calc.IMS) * 100
	}

	return c.JSON(http.StatusOK, calc)
}

func GetPrimaryIntelligenceDetails(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	channelID := c.QueryParam("channel_id")

	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type IntelligenceDetails struct {
		ChannelID         int     `json:"channel_id"`
		ChannelName       string  `json:"channel_name"`
		Lifting           float64 `json:"lifting"`
		PrimaryCollection float64 `json:"primary_collection"`
		CollectionRate    float64 `json:"collection_rate"`
		DueAmount         float64 `json:"due_amount"`
	}

	var details []IntelligenceDetails
	query := `
		SELECT 
			channel_id,
			channel_name,
			COALESCE(SUM(billed), 0) as lifting,
			COALESCE(SUM(primary_collection), 0) as primary_collection,
			(COALESCE(SUM(billed), 0) - COALESCE(SUM(primary_collection), 0)) as due_amount,
			CASE 
				WHEN COALESCE(SUM(billed), 0) > 0 
				THEN (COALESCE(SUM(primary_collection), 0) / COALESCE(SUM(billed), 0)) * 100
				ELSE 0
			END as collection_rate
		FROM channelwise_monthly_report
		WHERE data_month >= ? AND data_month <= ?
	`
	params := []interface{}{startDate, endDate}

	if channelID != "" {
		query += " AND channel_id = ?"
		chID, _ := strconv.Atoi(channelID)
		params = append(params, chID)
	}

	query += " GROUP BY channel_id, channel_name ORDER BY lifting DESC"
	database.DB.Raw(query, params...).Scan(&details)

	return c.JSON(http.StatusOK, details)
}

func GetSecondaryIntelligenceDetails(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	channelID := c.QueryParam("channel_id")

	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type IntelligenceDetails struct {
		ChannelID        int     `json:"channel_id"`
		ChannelName      string  `json:"channel_name"`
		IMS              float64 `json:"ims"`
		MarketCollection float64 `json:"market_collection"`
		CollectionRate   float64 `json:"collection_rate"`
		DueAmount        float64 `json:"due_amount"`
	}

	var details []IntelligenceDetails
	query := `
		SELECT 
			channel_id,
			channel_name,
			COALESCE(SUM(ims), 0) as ims,
			COALESCE(SUM(market_collection), 0) as market_collection,
			(COALESCE(SUM(ims), 0) - COALESCE(SUM(market_collection), 0)) as due_amount,
			CASE 
				WHEN COALESCE(SUM(ims), 0) > 0 
				THEN (COALESCE(SUM(market_collection), 0) / COALESCE(SUM(ims), 0)) * 100
				ELSE 0
			END as collection_rate
		FROM channelwise_monthly_report
		WHERE data_month >= ? AND data_month <= ?
	`
	params := []interface{}{startDate, endDate}

	if channelID != "" {
		query += " AND channel_id = ?"
		chID, _ := strconv.Atoi(channelID)
		params = append(params, chID)
	}

	query += " GROUP BY channel_id, channel_name ORDER BY ims DESC"
	database.DB.Raw(query, params...).Scan(&details)

	return c.JSON(http.StatusOK, details)
}

func GetTopRetailerList(c echo.Context) error {
	date := c.QueryParam("date")
	limit := c.QueryParam("limit")
	if limit == "" {
		limit = "10"
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	limitInt, _ := strconv.Atoi(limit)
	targetDate, _ := time.Parse("2006-01-02", date)

	var retailers []models.TopRetailer
	database.DB.Where("date = ?", targetDate).
		Order("amount DESC").
		Limit(limitInt).
		Find(&retailers)

	return c.JSON(http.StatusOK, retailers)
}
