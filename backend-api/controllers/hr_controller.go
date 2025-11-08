package controllers

import (
	"encoding/json"
	"fmt"
	"idash-backend-api/database"
	"idash-backend-api/models"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

func GetHRSummary(c echo.Context) error {
	month := c.QueryParam("month")
	department := c.QueryParam("department")

	if month == "" {
		month = time.Now().Format("2006-01")
	}

	type HRSummary struct {
		TotalEmployees   int     `json:"total_employees"`
		PresentEmployees int     `json:"present_employees"`
		AbsentEmployees  int     `json:"absent_employees"`
		TotalPromotions  int     `json:"total_promotions"`
		AttendanceRate   float64 `json:"attendance_rate"`
	}

	var summary HRSummary

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	// Get total employees
	empQuery := database.DB.Model(&models.EmployeeBasicInfo{}).Where("status = ?", 1)
	if department != "" {
		empQuery = empQuery.Where("department = ?", department)
	}
	empQuery.Count(&summary.TotalEmployees)

	// Get attendance data
	attendanceQuery := `
		SELECT 
			COUNT(DISTINCT employee_id) as present_employees
		FROM employee_attendances
		WHERE date >= ? AND date <= ? AND status = 'present'
	`
	attendanceParams := []interface{}{startDate, endDate}
	if department != "" {
		attendanceQuery = `
			SELECT 
				COUNT(DISTINCT ea.employee_id) as present_employees
			FROM employee_attendances ea
			JOIN employee_basic_infos ebi ON ea.employee_id = ebi.employee_id
			WHERE ea.date >= ? AND ea.date <= ? AND ea.status = 'present' AND ebi.department = ?
		`
		attendanceParams = append(attendanceParams, department)
	}
	database.DB.Raw(attendanceQuery, attendanceParams...).Scan(&summary.PresentEmployees)

	summary.AbsentEmployees = summary.TotalEmployees - summary.PresentEmployees
	if summary.TotalEmployees > 0 {
		summary.AttendanceRate = (float64(summary.PresentEmployees) / float64(summary.TotalEmployees)) * 100
	}

	// Get promotions
	promotionQuery := database.DB.Model(&models.YearlyEmployeePromotion{}).Where("YEAR(promotion_date) = ?", year)
	if department != "" {
		promotionQuery = promotionQuery.Joins("JOIN employee_basic_infos ON yearly_employee_promotions.employee_id = employee_basic_infos.employee_id").
			Where("employee_basic_infos.department = ?", department)
	}
	promotionQuery.Count(&summary.TotalPromotions)

	return c.JSON(http.StatusOK, summary)
}

func GetEmployees(c echo.Context) error {
	department := c.QueryParam("department")
	status := c.QueryParam("status")
	if status == "" {
		status = "1"
	}

	var employees []models.EmployeeBasicInfo
	query := database.DB.Where("status = ?", status)

	if department != "" {
		query = query.Where("department = ?", department)
	}

	query.Find(&employees)

	return c.JSON(http.StatusOK, employees)
}

func GetEmployeeAttendance(c echo.Context) error {
	fromDate := c.QueryParam("fromdate")
	toDate := c.QueryParam("todate")
	employeeID := c.QueryParam("employee_id")
	department := c.QueryParam("department")

	if fromDate == "" {
		fromDate = time.Now().Format("2006-01-02")
	}
	if toDate == "" {
		toDate = time.Now().Format("2006-01-02")
	}

	from, _ := time.Parse("2006-01-02", fromDate)
	to, _ := time.Parse("2006-01-02", toDate)

	var attendance []models.EmployeeAttendance
	query := database.DB.Where("date >= ? AND date <= ?", from, to)

	if employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}

	if department != "" {
		query = query.Joins("JOIN employee_basic_infos ON employee_attendances.employee_id = employee_basic_infos.employee_id").
			Where("employee_basic_infos.department = ?", department)
	}

	query.Find(&attendance)

	return c.JSON(http.StatusOK, attendance)
}

func GetEmployeePromotions(c echo.Context) error {
	year := c.QueryParam("year")
	department := c.QueryParam("department")
	employeeID := c.QueryParam("employee_id")

	if year == "" {
		year = strconv.Itoa(time.Now().Year())
	}

	yearInt, _ := strconv.Atoi(year)

	var promotions []models.YearlyEmployeePromotion
	query := database.DB.Where("YEAR(promotion_date) = ?", yearInt)

	if employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}

	if department != "" {
		query = query.Joins("JOIN employee_basic_infos ON yearly_employee_promotions.employee_id = employee_basic_infos.employee_id").
			Where("employee_basic_infos.department = ?", department)
	}

	query.Find(&promotions)

	return c.JSON(http.StatusOK, promotions)
}

func GetDepartmentAnalytics(c echo.Context) error {
	month := c.QueryParam("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type DepartmentAnalytics struct {
		Department      string  `json:"department"`
		TotalEmployees  int     `json:"total_employees"`
		PresentCount    int     `json:"present_count"`
		AbsentCount     int     `json:"absent_count"`
		AttendanceRate  float64 `json:"attendance_rate"`
		PromotionsCount int     `json:"promotions_count"`
	}

	var analytics []DepartmentAnalytics
	query := `
		SELECT 
			ebi.department,
			COUNT(DISTINCT ebi.employee_id) as total_employees,
			COUNT(DISTINCT CASE WHEN ea.status = 'present' THEN ea.employee_id END) as present_count,
			COUNT(DISTINCT CASE WHEN ea.status = 'absent' THEN ea.employee_id END) as absent_count,
			COUNT(DISTINCT yep.id) as promotions_count
		FROM employee_basic_infos ebi
		LEFT JOIN employee_attendances ea ON ebi.employee_id = ea.employee_id 
			AND ea.date >= ? AND ea.date <= ?
		LEFT JOIN yearly_employee_promotions yep ON ebi.employee_id = yep.employee_id 
			AND YEAR(yep.promotion_date) = ?
		WHERE ebi.status = 1
		GROUP BY ebi.department
		ORDER BY ebi.department
	`
	database.DB.Raw(query, startDate, endDate, year).Scan(&analytics)

	// Calculate attendance rate
	for i := range analytics {
		if analytics[i].TotalEmployees > 0 {
			analytics[i].AttendanceRate = (float64(analytics[i].PresentCount) / float64(analytics[i].TotalEmployees)) * 100
		}
	}

	return c.JSON(http.StatusOK, analytics)
}

func GetEmployeeTranOver(c echo.Context) error {
	fromDate := c.QueryParam("fromdate")
	toDate := c.QueryParam("todate")
	employeeID := c.QueryParam("employee_id")
	tranType := c.QueryParam("type") // transfer or overtime

	if fromDate == "" {
		fromDate = time.Now().Format("2006-01-02")
	}
	if toDate == "" {
		toDate = time.Now().Format("2006-01-02")
	}

	from, _ := time.Parse("2006-01-02", fromDate)
	to, _ := time.Parse("2006-01-02", toDate)

	var tranOver []models.EmployeeTranOver
	query := database.DB.Where("date >= ? AND date <= ?", from, to)

	if employeeID != "" {
		query = query.Where("employee_id = ?", employeeID)
	}

	if tranType != "" {
		query = query.Where("type = ?", tranType)
	}

	query.Find(&tranOver)

	return c.JSON(http.StatusOK, tranOver)
}

func GetEmployeeTranOverSummary(c echo.Context) error {
	month := c.QueryParam("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(month[:4])
	monthInt, _ := strconv.Atoi(month[5:7])
	startDate := time.Date(year, time.Month(monthInt), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	type TranOverSummary struct {
		Type            string  `json:"type"`
		TotalAmount     float64 `json:"total_amount"`
		TotalCount      int     `json:"total_count"`
	}

	var summary []TranOverSummary
	query := `
		SELECT 
			type,
			COALESCE(SUM(amount), 0) as total_amount,
			COUNT(*) as total_count
		FROM employee_tran_overs
		WHERE date >= ? AND date <= ?
		GROUP BY type
	`
	database.DB.Raw(query, startDate, endDate).Scan(&summary)

	return c.JSON(http.StatusOK, summary)
}

type HRKPI struct {
	Label  string  `json:"label"`
	Value  float64 `json:"value"`
	Change float64 `json:"change"`
}

type DepartmentPerformance struct {
	Department     string  `json:"department"`
	TotalEmployees int     `json:"total_employees"`
	PresentCount   int     `json:"present_count"`
	AttendanceRate float64 `json:"attendance_rate"`
	Promotions     int     `json:"promotions"`
}

type WorkforceMovement struct {
	Type       string  `json:"type"`
	TotalCount int     `json:"total_count"`
	TotalAmount float64 `json:"total_amount"`
}

type HRTrendPoint struct {
	Month          string  `json:"month"`
	Headcount      float64 `json:"headcount"`
	AttritionCount float64 `json:"attrition_count"`
	OvertimeHours  float64 `json:"overtime_hours"`
}

type HRForecastPoint struct {
	Date       string  `json:"date"`
	Forecast   float64 `json:"forecast"`
	UpperBound float64 `json:"upper_bound"`
	LowerBound float64 `json:"lower_bound"`
}

type HRForecastSummary struct {
	TotalForecast        float64           `json:"total_forecast"`
	AverageDailyForecast float64           `json:"average_daily_forecast"`
	ConfidenceLevel      float64           `json:"confidence_level"`
	ModelUsed            string            `json:"model_used"`
	ForecastData         []HRForecastPoint `json:"forecast_data"`
}

type HROverviewResponse struct {
	KPIs        []HRKPI               `json:"kpis"`
	Departments []DepartmentPerformance `json:"departments"`
	Movements   []WorkforceMovement   `json:"movements"`
	Trend       []HRTrendPoint        `json:"trend"`
	Forecast    *HRForecastSummary    `json:"forecast"`
	Alerts      []string              `json:"alerts"`
	LastUpdated time.Time             `json:"last_updated"`
}

type hrSummaryRow struct {
	Headcount float64
	Present   float64
	Attrition float64
	Overtime  float64
}

type hrForecastDetail struct {
	Date       string  `json:"date"`
	Forecast   float64 `json:"forecast"`
	UpperBound float64 `json:"upper_bound"`
	LowerBound float64 `json:"lower_bound"`
}

func GetHROverview(c echo.Context) error {
	yearMonth := c.QueryParam("yearMonth")
	department := c.QueryParam("department")
	if yearMonth == "" {
		yearMonth = time.Now().Format("2006-01")
	}

	year, _ := strconv.Atoi(yearMonth[:4])
	month, _ := strconv.Atoi(yearMonth[5:7])
	startDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

	prevStart := startDate.AddDate(0, -1, 0)
	prevEnd := endDate.AddDate(0, -1, 0)

	resp := HROverviewResponse{
		KPIs:        make([]HRKPI, 0),
		Departments: make([]DepartmentPerformance, 0),
		Movements:   make([]WorkforceMovement, 0),
		Trend:       make([]HRTrendPoint, 0),
		Alerts:      make([]string, 0),
		LastUpdated: time.Now(),
	}

	// Current headcount
	var headcount int64
	empQuery := database.DB.Model(&models.EmployeeBasicInfo{}).Where("status = ?", 1)
	if department != "" {
		empQuery = empQuery.Where("department = ?", department)
	}
	empQuery.Count(&headcount)

	// Attendance current month
	var presentCount int64
	attendanceQuery := database.DB.Model(&models.EmployeeAttendance{}).
		Where("date BETWEEN ? AND ? AND status = ?", startDate, endDate, "present")
	if department != "" {
		attendanceQuery = attendanceQuery.Joins("JOIN employee_basic_infos ebi ON ebi.employee_id = employee_attendances.employee_id").
			Where("ebi.department = ?", department)
	}
	attendanceQuery.Count(&presentCount)

	var previousHeadcount int64
	prevEmpQuery := database.DB.Model(&models.EmployeeBasicInfo{}).Where("status = ?", 1)
	if department != "" {
		prevEmpQuery = prevEmpQuery.Where("department = ?", department)
	}
	prevEmpQuery.Count(&previousHeadcount)

	// Attrition count (assuming employee_tran_overs with type='attrition')
	var attritionCount int64
	attritionQuery := database.DB.Model(&models.EmployeeTranOver{}).
		Where("date BETWEEN ? AND ? AND type = ?", startDate, endDate, "attrition")
	if department != "" {
		attritionQuery = attritionQuery.Joins("JOIN employee_basic_infos ebi ON ebi.employee_id = employee_tran_overs.employee_id").
			Where("ebi.department = ?", department)
	}
	attritionQuery.Count(&attritionCount)

	// Overtime hours
	var overtimeHours float64
	overtimeQuery := `
		SELECT COALESCE(SUM(amount),0)
		FROM employee_tran_overs eto
		` + departmentJoinClause(department) + `
		WHERE eto.date BETWEEN ? AND ? AND eto.type = 'overtime'` + departmentWhereClause(department)
	database.DB.Raw(overtimeQuery, departmentParams(department, startDate, endDate)...).Scan(&overtimeHours)

	delta := func(cur, prev float64) float64 {
		if prev == 0 {
			return 0
		}
		return ((cur - prev) / prev) * 100
	}

	attendanceRate := 0.0
	if headcount > 0 {
		attendanceRate = (float64(presentCount) / float64(headcount)) * 100
	}

	attritionRate := 0.0
	if headcount > 0 {
		attritionRate = (float64(attritionCount) / float64(headcount)) * 100
	}

	resp.KPIs = append(resp.KPIs,
		HRKPI{Label: "Headcount", Value: float64(headcount), Change: delta(float64(headcount), float64(previousHeadcount))},
		HRKPI{Label: "Attendance %", Value: attendanceRate, Change: 0},
		HRKPI{Label: "Attrition %", Value: attritionRate, Change: 0},
		HRKPI{Label: "Present Employees", Value: float64(presentCount), Change: 0},
		HRKPI{Label: "Attrition Count", Value: float64(attritionCount), Change: 0},
		HRKPI{Label: "Overtime Hours", Value: overtimeHours, Change: 0},
	)

	// Department performance reuse existing analytics query
	deptAnalytics := []struct {
		Department string
		TotalEmployees int
		PresentCount int
		AttendanceRate float64
		PromotionsCount int
	}{}
	database.DB.Raw(`
		SELECT 
			ebi.department,
			COUNT(DISTINCT ebi.employee_id) as total_employees,
			COUNT(DISTINCT CASE WHEN ea.status = 'present' THEN ea.employee_id END) as present_count,
			COUNT(DISTINCT CASE WHEN ea.status = 'present' THEN ea.employee_id END) / NULLIF(COUNT(DISTINCT ebi.employee_id),0) * 100 as attendance_rate,
			COUNT(DISTINCT yep.id) as promotions_count
		FROM employee_basic_infos ebi
		LEFT JOIN employee_attendances ea ON ebi.employee_id = ea.employee_id AND ea.date BETWEEN ? AND ?
		LEFT JOIN yearly_employee_promotions yep ON ebi.employee_id = yep.employee_id AND YEAR(yep.promotion_date) = ?
		WHERE ebi.status = 1` + departmentFilterClause(department) + `
		GROUP BY ebi.department
	`, append([]interface{}{startDate, endDate, year}, departmentFilterParams(department)...)...).Scan(&deptAnalytics)

	for _, row := range deptAnalytics {
		resp.Departments = append(resp.Departments, DepartmentPerformance{
			Department:     row.Department,
			TotalEmployees: row.TotalEmployees,
			PresentCount:   row.PresentCount,
			AttendanceRate: row.AttendanceRate,
			Promotions:     row.PromotionsCount,
		})
	}

	sort.Slice(resp.Departments, func(i, j int) bool {
		return resp.Departments[i].AttendanceRate > resp.Departments[j].AttendanceRate
	})

	// Workforce movement summary via existing helper
	movementRows := []WorkforceMovement{}
	database.DB.Raw(`
		SELECT type, COUNT(*) as total_count, COALESCE(SUM(amount),0) as total_amount
		FROM employee_tran_overs eto
		` + departmentJoinClause(department) + `
		WHERE eto.date BETWEEN ? AND ?
		GROUP BY type
	`, departmentParams(department, startDate, endDate)...).Scan(&movementRows)
	resp.Movements = movementRows

	// Trend for last 12 months using employee_basic_infos snapshot and attrition/overtime
	trendStart := startDate.AddDate(0, -11, 0)
	var trendMonths []string
	for i := 0; i < 12; i++ {
		monthTime := trendStart.AddDate(0, i, 0)
		trendMonths = append(trendMonths, monthTime.Format("2006-01"))
	}

	for _, monthStr := range trendMonths {
		trendYear, _ := strconv.Atoi(monthStr[:4])
		trendMonth, _ := strconv.Atoi(monthStr[5:7])
		trendStartDate := time.Date(trendYear, time.Month(trendMonth), 1, 0, 0, 0, 0, time.UTC)
		trendEndDate := trendStartDate.AddDate(0, 1, 0).Add(-24 * time.Hour)

		var monthHeadcount int64
		empTrendQuery := database.DB.Model(&models.EmployeeBasicInfo{}).Where("status = ?", 1)
		if department != "" {
			empTrendQuery = empTrendQuery.Where("department = ?", department)
		}
		empTrendQuery.Count(&monthHeadcount)

		var monthAttrition int64
		attrTrendQuery := database.DB.Model(&models.EmployeeTranOver{}).
			Where("date BETWEEN ? AND ? AND type = ?", trendStartDate, trendEndDate, "attrition")
		if department != "" {
			attrTrendQuery = attrTrendQuery.Joins("JOIN employee_basic_infos ebi ON ebi.employee_id = employee_tran_overs.employee_id").
				Where("ebi.department = ?", department)
		}
		attrTrendQuery.Count(&monthAttrition)

		var monthOvertime float64
		database.DB.Raw(`
			SELECT COALESCE(SUM(amount),0)
			FROM employee_tran_overs eto
			`+ departmentJoinClause(department) +`
			WHERE eto.date BETWEEN ? AND ? AND eto.type = 'overtime'`+ departmentWhereClause(department), departmentParams(department, trendStartDate, trendEndDate)...).Scan(&monthOvertime)

		resp.Trend = append(resp.Trend, HRTrendPoint{
			Month:          monthStr,
			Headcount:      float64(monthHeadcount),
			AttritionCount: float64(monthAttrition),
			OvertimeHours:  monthOvertime,
		})
	}

	// Forecast (attrition risk)
	var forecastModel models.AIForecast
	if err := database.DB.Where("forecast_type = ?", "hr").Order("created_at DESC").First(&forecastModel).Error; err == nil {
		var details []hrForecastDetail
		if forecastModel.ForecastDetails != "" {
			_ = json.Unmarshal([]byte(forecastModel.ForecastDetails), &details)
		}

		forecastSummary := &HRForecastSummary{
			TotalForecast:   forecastModel.ForecastedValue,
			ConfidenceLevel: forecastModel.ConfidenceLevel,
			ModelUsed:       forecastModel.ModelUsed,
			ForecastData:    make([]HRForecastPoint, 0, len(details)),
		}

		if len(details) > 0 {
			total := 0.0
			for _, d := range details {
				forecastSummary.ForecastData = append(forecastSummary.ForecastData, HRForecastPoint{
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
	if attendanceRate < 90 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Attendance rate %.1f%% below threshold", attendanceRate))
	}
	if attritionRate > 5 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Attrition rate %.1f%% exceeds limit", attritionRate))
	}
	if overtimeHours > 500 {
		resp.Alerts = append(resp.Alerts, fmt.Sprintf("Overtime hours at %.0f, review workforce planning", overtimeHours))
	}
	if resp.Forecast != nil && resp.Forecast.TotalForecast > float64(attritionCount)*1.2 {
		resp.Alerts = append(resp.Alerts, "AI forecast indicates rising attrition risk")
	}

	return c.JSON(http.StatusOK, resp)
}

func departmentJoinClause(department string) string {
	if department != "" {
		return "JOIN employee_basic_infos ebi ON ebi.employee_id = eto.employee_id"
	}
	return ""
}

func departmentWhereClause(department string) string {
	if department != "" {
		return " AND ebi.department = ?"
	}
	return ""
}

func departmentParams(department string, start, end time.Time) []interface{} {
	params := []interface{}{start, end}
	if department != "" {
		params = append(params, department)
	}
	return params
}

func departmentFilterClause(department string) string {
	if department != "" {
		return " AND ebi.department = ?"
	}
	return ""
}

func departmentFilterParams(department string) []interface{} {
	if department != "" {
		return []interface{}{department}
	}
	return []interface{}{}
}
