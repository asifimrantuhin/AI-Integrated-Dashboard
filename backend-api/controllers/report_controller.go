package controllers

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"idash-backend-api/database"

	"github.com/labstack/echo/v4"
)

type ReportDefinition struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Endpoint    string   `json:"endpoint"`
	Filters     []string `json:"filters"`
	Formats     []string `json:"formats"`
}

var availableReports = []ReportDefinition{
	{ID: "sales-monthly", Name: "Sales Monthly Summary", Description: "Channel and product sales by month", Endpoint: "/api/sales/summary", Filters: []string{"yearMonth", "channel"}, Formats: []string{"json", "csv"}},
	{ID: "production-efficiency", Name: "Production Efficiency", Description: "Factory efficiency and downtime", Endpoint: "/api/production/overview", Filters: []string{"factory"}, Formats: []string{"json", "csv"}},
	{ID: "finance-budget", Name: "Budget vs Actual", Description: "Variance by category and department", Endpoint: "/api/finance/overview", Filters: []string{"department"}, Formats: []string{"json", "csv"}},
}

func ListReports(c echo.Context) error {
	return c.JSON(http.StatusOK, availableReports)
}

func GenerateReport(c echo.Context) error {
	reportID := c.Param("id")
	format := c.QueryParam("format")
	if format == "" {
		format = "json"
	}

	var selected *ReportDefinition
	for _, r := range availableReports {
		if r.ID == reportID {
			selected = &r
			break
		}
	}

	if selected == nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "report not found"})
	}

	data, err := executeReportQuery(reportID, c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	switch format {
	case "csv":
		return streamCSV(c, reportID, data)
	default:
		return c.JSON(http.StatusOK, data)
	}
}

func executeReportQuery(reportID string, c echo.Context) ([]map[string]interface{}, error) {
	var rows *sql.Rows
	var err error
	switch reportID {
	case "sales-monthly":
		yearMonth := c.QueryParam("yearMonth")
		if yearMonth == "" {
			yearMonth = time.Now().Format("2006-01")
		}
		year, _ := strconv.Atoi(yearMonth[:4])
		month, _ := strconv.Atoi(yearMonth[5:7])
		startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
		endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)
		rows, err = database.DB.Raw(`SELECT channel_name, data_month, SUM(billed) billed, SUM(ims) ims, SUM(primary_collection) primary_collection FROM channelwise_monthly_report WHERE data_month BETWEEN ? AND ? GROUP BY channel_name, data_month ORDER BY data_month`, startDate, endDate).Rows()
	case "production-efficiency":
		rows, err = database.DB.Raw(`SELECT factory_id, production_line_id, AVG(efficiency_percentage) efficiency, SUM(downtime_minutes) downtime FROM production_efficiency GROUP BY factory_id, production_line_id`).Rows()
	case "finance-budget":
		rows, err = database.DB.Raw(`SELECT department_id, SUM(budget_amount) budget, SUM(actual_amount) actual FROM budget_summaries GROUP BY department_id`).Rows()
	default:
		return nil, fmt.Errorf("report not implemented")
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	results := make([]map[string]interface{}, 0)

	for rows.Next() {
		record := make([]interface{}, len(columns))
		recordPtrs := make([]interface{}, len(columns))
		for i := range record {
			recordPtrs[i] = &record[i]
		}

		if err := rows.Scan(recordPtrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{}, len(columns))
		for i, col := range columns {
			value := record[i]
			if b, ok := value.([]byte); ok {
				rowMap[col] = string(b)
			} else {
				rowMap[col] = value
			}
		}
		results = append(results, rowMap)
	}

	return results, nil
}

func streamCSV(c echo.Context, reportID string, data []map[string]interface{}) error {
	if len(data) == 0 {
		return c.Blob(http.StatusOK, "text/csv", []byte{})
	}

	headers := make([]string, 0, len(data[0]))
	for key := range data[0] {
		headers = append(headers, key)
	}

	c.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s.csv", reportID))
	c.Response().Header().Set(echo.HeaderContentType, "text/csv")

	writer := csv.NewWriter(c.Response())
	if err := writer.Write(headers); err != nil {
		return err
	}

	for _, row := range data {
		record := make([]string, len(headers))
		for i, header := range headers {
			value := row[header]
			switch v := value.(type) {
			case nil:
				record[i] = ""
			case string:
				record[i] = v
			default:
				record[i] = fmt.Sprintf("%v", v)
			}
		}
		writer.Write(record)
	}

	writer.Flush()
	return writer.Error()
}
