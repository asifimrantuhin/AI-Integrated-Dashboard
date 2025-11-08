package controllers

import (
	"encoding/json"
	"idash-backend-api/database"
	"idash-backend-api/models"
	"idash-backend-api/utils"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

// GetQualityControlData returns quality control data with filters
func GetQualityControlData(c echo.Context) error {
	params := ParseFilterParams(c)
	
	cacheKey := utils.GenerateCacheKey("/api/manufacturing/quality", map[string]interface{}{
		"start_date": params.StartDate,
		"end_date":   params.EndDate,
		"factory_id": params.FactoryID,
		"status":     params.Status,
	})
	
	if cached, found := utils.GetFromCache(database.DB, cacheKey); found {
		return c.JSONBlob(http.StatusOK, []byte(cached))
	}
	
	query := database.DB.Table("quality_control_checks")
	
	if params.StartDate != nil {
		query = query.Where("check_date >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("check_date <= ?", *params.EndDate)
	}
	if params.CompanyID != nil {
		query = query.Where("company_id = ?", *params.CompanyID)
	}
	if params.FactoryID != nil {
		query = query.Where("factory_id = ?", *params.FactoryID)
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	
	var total int64
	query.Count(&total)
	
	var results []models.QualityControlCheck
	orderBy := params.SortBy + " " + params.SortOrder
	query.Order(orderBy).Limit(params.Limit).Offset(params.Offset).Find(&results)
	
	// Calculate summary
	var passedCount, failedCount int64
	var totalChecked int64
	database.DB.Table("quality_control_checks").
		Where("check_date >= ? AND check_date <= ?", params.StartDate, params.EndDate).
		Select("SUM(passed_count) as passed, SUM(failed_count) as failed, SUM(total_checked) as total").
		Row().Scan(&passedCount, &failedCount, &totalChecked)
	
	response := map[string]interface{}{
		"data": results,
		"total": total,
		"summary": map[string]interface{}{
			"total_checked": totalChecked,
			"passed_count":  passedCount,
			"failed_count":  failedCount,
			"pass_rate":     float64(passedCount) / float64(totalChecked+1) * 100,
		},
		"limit":  params.Limit,
		"offset": params.Offset,
	}
	
	responseJSON, _ := json.Marshal(response)
	utils.SetCache(database.DB, cacheKey, "/api/manufacturing/quality", nil, string(responseJSON), 5*time.Minute)
	
	return c.JSON(http.StatusOK, response)
}

// GetProductionPlans returns production plans
func GetProductionPlans(c echo.Context) error {
	params := ParseFilterParams(c)
	
	query := database.DB.Table("production_plans")
	
	if params.StartDate != nil {
		query = query.Where("plan_date >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("plan_date <= ?", *params.EndDate)
	}
	if params.CompanyID != nil {
		query = query.Where("company_id = ?", *params.CompanyID)
	}
	if params.FactoryID != nil {
		query = query.Where("factory_id = ?", *params.FactoryID)
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	
	var total int64
	query.Count(&total)
	
	var results []models.ProductionPlan
	orderBy := params.SortBy + " " + params.SortOrder
	query.Order(orderBy).Limit(params.Limit).Offset(params.Offset).Find(&results)
	
	response := map[string]interface{}{
		"data":   results,
		"total":  total,
		"limit":  params.Limit,
		"offset": params.Offset,
	}
	
	return c.JSON(http.StatusOK, response)
}

// GetMachineMaintenance returns machine maintenance data
func GetMachineMaintenance(c echo.Context) error {
	params := ParseFilterParams(c)
	
	query := database.DB.Table("machine_maintenances")
	
	if params.StartDate != nil {
		query = query.Where("maintenance_date >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("maintenance_date <= ?", *params.EndDate)
	}
	if params.FactoryID != nil {
		query = query.Where("factory_id = ?", *params.FactoryID)
	}
	if params.Status != nil {
		query = query.Where("status = ?", *params.Status)
	}
	
	var total int64
	query.Count(&total)
	
	var results []models.MachineMaintenance
	orderBy := params.SortBy + " " + params.SortOrder
	query.Order(orderBy).Limit(params.Limit).Offset(params.Offset).Find(&results)
	
	response := map[string]interface{}{
		"data":   results,
		"total":  total,
		"limit":  params.Limit,
		"offset": params.Offset,
	}
	
	return c.JSON(http.StatusOK, response)
}

// GetProductionEfficiency returns production efficiency metrics
func GetProductionEfficiency(c echo.Context) error {
	params := ParseFilterParams(c)
	
	query := database.DB.Table("production_efficiency")
	
	if params.StartDate != nil {
		query = query.Where("production_date >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("production_date <= ?", *params.EndDate)
	}
	if params.FactoryID != nil {
		query = query.Where("factory_id = ?", *params.FactoryID)
	}
	
	var total int64
	query.Count(&total)
	
	var results []models.ProductionEfficiency
	orderBy := params.SortBy + " " + params.SortOrder
	query.Order(orderBy).Limit(params.Limit).Offset(params.Offset).Find(&results)
	
	// Calculate average efficiency
	var avgEfficiency float64
	database.DB.Table("production_efficiency").
		Where("production_date >= ? AND production_date <= ?", params.StartDate, params.EndDate).
		Select("AVG(efficiency_percentage)").Row().Scan(&avgEfficiency)
	
	response := map[string]interface{}{
		"data":            results,
		"total":           total,
		"average_efficiency": avgEfficiency,
		"limit":           params.Limit,
		"offset":          params.Offset,
	}
	
	return c.JSON(http.StatusOK, response)
}

// GetSupplierPerformance returns supplier performance data
func GetSupplierPerformance(c echo.Context) error {
	params := ParseFilterParams(c)
	
	query := database.DB.Table("supplier_performance")
	
	if params.StartDate != nil {
		query = query.Where("evaluation_date >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("evaluation_date <= ?", *params.EndDate)
	}
	
	var total int64
	query.Count(&total)
	
	var results []models.SupplierPerformance
	orderBy := params.SortBy + " " + params.SortOrder
	query.Order(orderBy).Limit(params.Limit).Offset(params.Offset).Find(&results)
	
	response := map[string]interface{}{
		"data":   results,
		"total":  total,
		"limit":  params.Limit,
		"offset": params.Offset,
	}
	
	return c.JSON(http.StatusOK, response)
}

// GetEnergyConsumption returns energy consumption data
func GetEnergyConsumption(c echo.Context) error {
	params := ParseFilterParams(c)
	
	query := database.DB.Table("energy_consumption")
	
	if params.StartDate != nil {
		query = query.Where("consumption_date >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("consumption_date <= ?", *params.EndDate)
	}
	if params.FactoryID != nil {
		query = query.Where("factory_id = ?", *params.FactoryID)
	}
	
	var total int64
	query.Count(&total)
	
	var results []models.EnergyConsumption
	orderBy := params.SortBy + " " + params.SortOrder
	query.Order(orderBy).Limit(params.Limit).Offset(params.Offset).Find(&results)
	
	// Calculate total consumption and cost
	var totalConsumption, totalCost float64
	database.DB.Table("energy_consumption").
		Where("consumption_date >= ? AND consumption_date <= ?", params.StartDate, params.EndDate).
		Select("SUM(consumption_amount) as total_consumption, SUM(cost) as total_cost").
		Row().Scan(&totalConsumption, &totalCost)
	
	response := map[string]interface{}{
		"data":              results,
		"total":             total,
		"total_consumption": totalConsumption,
		"total_cost":        totalCost,
		"limit":             params.Limit,
		"offset":            params.Offset,
	}
	
	return c.JSON(http.StatusOK, response)
}

