package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"idash-backend-api/database"
	"idash-backend-api/models"

	"github.com/labstack/echo/v4"
)

func postToAIService(endpoint string, payload interface{}, target interface{}) error {
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

	respBody, _ := io.ReadAll(resp.Body)

	if resp.StatusCode >= 400 {
		// include response body for better diagnostics
		return fmt.Errorf("ai service returned status %s: %s", resp.Status, strings.TrimSpace(string(respBody)))
	}

	if target != nil {
		if err := json.Unmarshal(respBody, target); err != nil {
			return fmt.Errorf("failed to decode ai response: %w; body=%s", err, strings.TrimSpace(string(respBody)))
		}
	}

	return nil
}

func mustJSON(val interface{}) []byte {
	bytesValue, err := json.Marshal(val)
	if err != nil {
		return []byte("{}")
	}
	return bytesValue
}

type ForecastRequest struct {
	ForecastType   string  `json:"forecast_type"` // sales, production, finance, inventory
	Period         string  `json:"period"`        // daily, weekly, monthly
	StartDate      string  `json:"start_date"`
	EndDate        string  `json:"end_date"`
	Days           int     `json:"days"` // Forecast horizon
	CompanyID      *int    `json:"company_id"`
	FactoryID      *int    `json:"factory_id"`
	ChannelID      *int    `json:"channel_id"`
	ProductID      *int    `json:"product_id"`
	BudgetCategory *string `json:"budget_category"`
}

type ForecastResponse struct {
	ForecastType         string                   `json:"forecast_type"`
	ForecastID           string                   `json:"forecast_id"`
	EntityID             interface{}              `json:"entity_id"`
	ForecastPeriodDays   int                      `json:"forecast_period_days"`
	ConfidenceLevel      float64                  `json:"confidence_level"`
	TotalForecast        float64                  `json:"total_forecast"`
	AverageDailyForecast float64                  `json:"average_daily_forecast"`
	ProjectedGrowthRate  float64                  `json:"projected_growth_rate,omitempty"`
	ForecastData         []map[string]interface{} `json:"forecast_data"`
	ModelUsed            string                   `json:"model_used"`
	CreatedAt            string                   `json:"created_at"`
}

type SalesPredictionResponse struct {
	Granularity     string                   `json:"granularity"`
	Horizon         int                      `json:"horizon"`
	Predictions     []map[string]interface{} `json:"predictions"`
	Recommendations []map[string]interface{} `json:"recommendations"`
	GeneratedAt     string                   `json:"generated_at"`
}

type InventoryPrescriptionResponse struct {
	Forecast      map[string]interface{} `json:"forecast"`
	Prescriptions map[string]interface{} `json:"prescriptions"`
	GeneratedAt   string                 `json:"generated_at"`
}

type FinancialPrescriptionResponse struct {
	Forecast      map[string]interface{} `json:"forecast"`
	Prescriptions map[string]interface{} `json:"prescriptions"`
	GeneratedAt   string                 `json:"generated_at"`
}

type ScenarioSimulationResponse struct {
	Horizon           int                `json:"horizon"`
	Inputs            map[string]float64 `json:"inputs"`
	ProjectedSales    float64            `json:"projected_sales"`
	ProjectedMargin   float64            `json:"projected_margin"`
	IncrementalProfit float64            `json:"incremental_profit"`
	Narrative         string             `json:"narrative"`
}

// RequestSalesForecast generates sales forecast
func RequestSalesForecast(c echo.Context) error {
	var req ForecastRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	req.ForecastType = "sales"
	if req.Days == 0 {
		req.Days = 30
	}
	if req.StartDate == "" || req.EndDate == "" {
		end := time.Now()
		start := end.AddDate(0, 0, -30)
		if req.StartDate == "" {
			req.StartDate = start.Format("2006-01-02")
		}
		if req.EndDate == "" {
			req.EndDate = end.Format("2006-01-02")
		}
	}

	// Get historical data from database
	var salesData []models.ChannelwiseLICdataMonthly
	query := database.DB.Table("channelwise_monthly_report")

	if req.StartDate != "" {
		startDate, _ := time.Parse("2006-01-02", req.StartDate)
		query = query.Where("data_month >= ?", startDate)
	}
	if req.EndDate != "" {
		endDate, _ := time.Parse("2006-01-02", req.EndDate)
		query = query.Where("data_month <= ?", endDate)
	}
	if req.ChannelID != nil {
		query = query.Where("channel_id = ?", *req.ChannelID)
	}

	query.Find(&salesData)

	// Call AI Service
	aiServiceURL := os.Getenv("AI_SERVICE_URL")
	if aiServiceURL == "" {
		aiServiceURL = "http://localhost:8000"
	}

	var forecastResp ForecastResponse
	if err := postToAIService("/api/forecast/sales", req, &forecastResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "AI service error: " + err.Error()})
	}

	// Store forecast in database
	forecast := models.AIForecast{
		ForecastType:    "sales",
		EntityType:      "channel",
		EntityID:        strconv.Itoa(*req.ChannelID),
		ForecastDate:    time.Now(),
		ForecastedValue: forecastResp.TotalForecast,
		ConfidenceLevel: forecastResp.ConfidenceLevel,
		ModelUsed:       forecastResp.ModelUsed,
		Status:          "active",
	}

	detailsJSON, _ := json.Marshal(forecastResp.ForecastData)
	forecast.ForecastDetails = string(detailsJSON)

	database.DB.Create(&forecast)
	forecastResp.ForecastID = strconv.Itoa(int(forecast.ID))

	return c.JSON(http.StatusOK, forecastResp)
}

// RequestProductionForecast generates production forecast
func RequestProductionForecast(c echo.Context) error {
	var req ForecastRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	req.ForecastType = "production"
	if req.Days == 0 {
		req.Days = 30
	}
	if req.StartDate == "" || req.EndDate == "" {
		end := time.Now()
		start := end.AddDate(0, 0, -30)
		if req.StartDate == "" {
			req.StartDate = start.Format("2006-01-02")
		}
		if req.EndDate == "" {
			req.EndDate = end.Format("2006-01-02")
		}
	}

	aiServiceURL := os.Getenv("AI_SERVICE_URL")
	if aiServiceURL == "" {
		aiServiceURL = "http://localhost:8000"
	}

	var forecastResp ForecastResponse
	if err := postToAIService("/api/forecast/production", req, &forecastResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "AI service error: " + err.Error()})
	}

	// Store forecast
	forecast := models.AIForecast{
		ForecastType:    "production",
		EntityType:      "factory",
		ForecastDate:    time.Now(),
		ForecastedValue: forecastResp.TotalForecast,
		ConfidenceLevel: forecastResp.ConfidenceLevel,
		ModelUsed:       forecastResp.ModelUsed,
		Status:          "active",
	}
	if req.FactoryID != nil {
		forecast.EntityID = strconv.Itoa(*req.FactoryID)
	}

	detailsJSON, _ := json.Marshal(forecastResp.ForecastData)
	forecast.ForecastDetails = string(detailsJSON)

	database.DB.Create(&forecast)
	forecastResp.ForecastID = strconv.Itoa(int(forecast.ID))

	return c.JSON(http.StatusOK, forecastResp)
}

// RequestFinancialForecast generates financial forecast
func RequestFinancialForecast(c echo.Context) error {
	var req ForecastRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	req.ForecastType = "finance"
	if req.Days == 0 {
		req.Days = 30
	}
	if req.StartDate == "" || req.EndDate == "" {
		end := time.Now()
		start := end.AddDate(0, 0, -30)
		if req.StartDate == "" {
			req.StartDate = start.Format("2006-01-02")
		}
		if req.EndDate == "" {
			req.EndDate = end.Format("2006-01-02")
		}
	}

	var forecastResp ForecastResponse
	if err := postToAIService("/api/forecast/finance", req, &forecastResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "AI service error: " + err.Error()})
	}

	// Store forecast
	forecast := models.AIForecast{
		ForecastType:    "finance",
		EntityType:      "budget",
		ForecastDate:    time.Now(),
		ForecastedValue: forecastResp.TotalForecast,
		ConfidenceLevel: forecastResp.ConfidenceLevel,
		ModelUsed:       forecastResp.ModelUsed,
		Status:          "active",
	}
	if req.BudgetCategory != nil {
		forecast.EntityID = *req.BudgetCategory
	}

	detailsJSON, _ := json.Marshal(forecastResp.ForecastData)
	forecast.ForecastDetails = string(detailsJSON)

	database.DB.Create(&forecast)
	forecastResp.ForecastID = strconv.Itoa(int(forecast.ID))

	return c.JSON(http.StatusOK, forecastResp)
}

// RequestInventoryForecast generates inventory forecast
func RequestInventoryForecast(c echo.Context) error {
	var req ForecastRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	req.ForecastType = "inventory"
	if req.Days == 0 {
		req.Days = 30
	}
	if req.StartDate == "" || req.EndDate == "" {
		end := time.Now()
		start := end.AddDate(0, 0, -30)
		if req.StartDate == "" {
			req.StartDate = start.Format("2006-01-02")
		}
		if req.EndDate == "" {
			req.EndDate = end.Format("2006-01-02")
		}
	}

	var forecastResp ForecastResponse
	if err := postToAIService("/api/forecast/inventory", req, &forecastResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "AI service error: " + err.Error()})
	}

	// Store forecast
	forecast := models.AIForecast{
		ForecastType:    "inventory",
		EntityType:      "product",
		ForecastDate:    time.Now(),
		ForecastedValue: forecastResp.TotalForecast,
		ConfidenceLevel: forecastResp.ConfidenceLevel,
		ModelUsed:       forecastResp.ModelUsed,
		Status:          "active",
	}
	if req.ProductID != nil {
		forecast.EntityID = strconv.Itoa(*req.ProductID)
	}

	detailsJSON, _ := json.Marshal(forecastResp.ForecastData)
	forecast.ForecastDetails = string(detailsJSON)

	database.DB.Create(&forecast)
	forecastResp.ForecastID = strconv.Itoa(int(forecast.ID))

	return c.JSON(http.StatusOK, forecastResp)
}

// GetForecast retrieves forecast by ID
func GetForecast(c echo.Context) error {
	forecastID := c.Param("id")

	var forecast models.AIForecast
	if err := database.DB.Where("id = ?", forecastID).First(&forecast).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Forecast not found"})
	}

	var forecastData []map[string]interface{}
	json.Unmarshal([]byte(forecast.ForecastDetails), &forecastData)

	response := ForecastResponse{
		ForecastType:    forecast.ForecastType,
		ForecastID:      forecastID,
		EntityID:        forecast.EntityID,
		ConfidenceLevel: forecast.ConfidenceLevel,
		ForecastData:    forecastData,
		ModelUsed:       forecast.ModelUsed,
		CreatedAt:       forecast.CreatedAt.Format(time.RFC3339),
	}

	return c.JSON(http.StatusOK, response)
}

// GetForecastsByType retrieves all forecasts of a specific type
func GetForecastsByType(c echo.Context) error {
	forecastType := c.Param("type")

	var forecasts []models.AIForecast
	database.DB.Where("forecast_type = ? AND status = ?", forecastType, "active").
		Order("created_at DESC").
		Limit(50).
		Find(&forecasts)

	return c.JSON(http.StatusOK, forecasts)
}

// AnalyzeData performs data analysis
func AnalyzeData(c echo.Context) error {
	type AnalysisRequest struct {
		Metric       string `json:"metric"`
		StartDate    string `json:"start_date"`
		EndDate      string `json:"end_date"`
		AnalysisType string `json:"analysis_type"`
	}

	var req AnalysisRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	aiServiceURL := os.Getenv("AI_SERVICE_URL")
	if aiServiceURL == "" {
		aiServiceURL = "http://localhost:8000"
	}

	reqBody, _ := json.Marshal(req)
	resp, err := http.Post(aiServiceURL+"/api/analyze", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to connect to AI service"})
	}
	defer resp.Body.Close()

	var analysisResp map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&analysisResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to parse AI service response"})
	}

	return c.JSON(http.StatusOK, analysisResp)
}

// PredictSalesSummary returns next-period sales predictions and prescriptions
func PredictSalesSummary(c echo.Context) error {
	granularity := c.QueryParam("granularity")
	if granularity == "" {
		granularity = "channel"
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, -18, 0)

	payload := map[string]interface{}{
		"start_date":  startDate.Format("2006-01-02"),
		"end_date":    endDate.Format("2006-01-02"),
		"horizon":     30,
		"granularity": granularity,
	}

	var aiResp SalesPredictionResponse
	if err := postToAIService("/api/predict/sales/summary", payload, &aiResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "AI prediction failed: " + err.Error()})
	}

	for _, rec := range aiResp.Recommendations {
		metadataBytes, _ := json.Marshal(rec)
		entityID := ""
		if idVal, ok := rec["entity_id"]; ok {
			entityID = fmt.Sprintf("%v", idVal)
		}
		impact := 0.0
		if lift, ok := rec["lift_vs_baseline"].(float64); ok {
			impact = lift
		}
		recommendation := fmt.Sprintf("%v", rec["recommended_action"])

		database.DB.Create(&models.PrescriptiveRecommendation{
			Module:         "sales",
			EntityType:     granularity,
			EntityID:       entityID,
			RiskLevel:      fmt.Sprintf("%v", rec["risk_level"]),
			Recommendation: recommendation,
			ImpactScore:    impact,
			Metadata:       string(metadataBytes),
		})
	}

	return c.JSON(http.StatusOK, aiResp)
}

// GetInventoryPrescription returns optimisation suggestions for stock
func GetInventoryPrescription(c echo.Context) error {
	productIDParam := c.QueryParam("productId")
	var productID *int
	if productIDParam != "" {
		if id, err := strconv.Atoi(productIDParam); err == nil {
			productID = &id
		}
	}

	endDate := time.Now()
	startDate := endDate.AddDate(0, -12, 0)

	payload := map[string]interface{}{
		"module":     "inventory",
		"start_date": startDate.Format("2006-01-02"),
		"end_date":   endDate.Format("2006-01-02"),
		"horizon":    60,
		"product_id": productID,
	}

	var aiResp InventoryPrescriptionResponse
	if err := postToAIService("/api/prescribe/inventory", payload, &aiResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "AI inventory prescription failed: " + err.Error()})
	}

	if actions, ok := aiResp.Prescriptions["inventory_actions"].([]interface{}); ok {
		for _, raw := range actions {
			recMap, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			metadataBytes, _ := json.Marshal(recMap)
			entity := fmt.Sprintf("%v", recMap["entity_id"])
			impact := 0.0
			if qty, ok := recMap["recommended_order_qty"].(float64); ok {
				impact = qty
			}
			action := fmt.Sprintf("%v", recMap["action"])

			database.DB.Create(&models.PrescriptiveRecommendation{
				Module:         "inventory",
				EntityType:     "product",
				EntityID:       entity,
				RiskLevel:      "medium",
				Recommendation: action,
				ImpactScore:    impact,
				Metadata:       string(metadataBytes),
			})
		}
	}

	return c.JSON(http.StatusOK, aiResp)
}

// GetFinancialPrescription surfaces budget optimisation actions
func GetFinancialPrescription(c echo.Context) error {
	category := c.QueryParam("category")

	endDate := time.Now()
	startDate := endDate.AddDate(0, -12, 0)

	payload := map[string]interface{}{
		"module":          "finance",
		"start_date":      startDate.Format("2006-01-02"),
		"end_date":        endDate.Format("2006-01-02"),
		"horizon":         90,
		"budget_category": category,
	}

	var aiResp FinancialPrescriptionResponse
	if err := postToAIService("/api/prescribe/finance", payload, &aiResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "AI financial prescription failed: " + err.Error()})
	}

	if actions, ok := aiResp.Prescriptions["financial_actions"].([]interface{}); ok {
		for _, raw := range actions {
			recMap, ok := raw.(map[string]interface{})
			if !ok {
				continue
			}
			metadataBytes, _ := json.Marshal(recMap)
			entity := fmt.Sprintf("%v", recMap["entity_id"])
			impact := 0.0
			if variance, ok := recMap["average_variance"].(float64); ok {
				impact = variance
			}

			database.DB.Create(&models.PrescriptiveRecommendation{
				Module:         "finance",
				EntityType:     "category",
				EntityID:       entity,
				RiskLevel:      "medium",
				Recommendation: fmt.Sprintf("%v", recMap["recommended_action"]),
				ImpactScore:    impact,
				Metadata:       string(metadataBytes),
			})
		}
	}

	return c.JSON(http.StatusOK, aiResp)
}

// AnalyzeAnomaliesWithActions proxies enriched anomaly detection
func AnalyzeAnomaliesWithActions(c echo.Context) error {
	metric := c.QueryParam("metric")
	if metric == "" {
		metric = "sales"
	}
	endDate := time.Now()
	startDate := endDate.AddDate(0, -6, 0)

	payload := map[string]interface{}{
		"metric":        metric,
		"analysis_type": "anomaly",
		"start_date":    startDate.Format("2006-01-02"),
		"end_date":      endDate.Format("2006-01-02"),
	}

	var aiResp map[string]interface{}
	if err := postToAIService("/api/analyze/anomaly/enriched", payload, &aiResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "AI anomaly analysis failed: " + err.Error()})
	}

	return c.JSON(http.StatusOK, aiResp)
}

// RunScenarioSimulation executes a what-if scenario via AI service
func RunScenarioSimulation(c echo.Context) error {
	type ScenarioRequestBody struct {
		Horizon     int                `json:"horizon"`
		BaseMetrics map[string]float64 `json:"base_metrics"`
		Adjustments map[string]float64 `json:"adjustments"`
	}

	var body ScenarioRequestBody
	if err := c.Bind(&body); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid scenario payload"})
	}

	if body.BaseMetrics == nil {
		body.BaseMetrics = map[string]float64{}
	}
	if body.Adjustments == nil {
		body.Adjustments = map[string]float64{}
	}
	if body.Horizon == 0 {
		body.Horizon = 90
	}

	var aiResp ScenarioSimulationResponse
	if err := postToAIService("/api/scenario/whatif", body, &aiResp); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Scenario simulation failed: " + err.Error()})
	}

	resultsBytes, _ := json.Marshal(aiResp)
	database.DB.Create(&models.ScenarioSimulation{
		Module:      "executive",
		ScenarioKey: "default",
		BaseMetrics: string(mustJSON(body.BaseMetrics)),
		Adjustments: string(mustJSON(body.Adjustments)),
		Results:     string(resultsBytes),
	})

	return c.JSON(http.StatusOK, aiResp)
}
