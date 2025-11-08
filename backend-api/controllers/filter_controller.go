package controllers

import (
	"encoding/json"
	"idash-backend-api/database"
	"idash-backend-api/utils"
	"net/http"
	"time"
	"strconv"
	"strings"
	
	"github.com/labstack/echo/v4"
)

// FilterParams represents common filter parameters
type FilterParams struct {
	StartDate      *time.Time `json:"start_date"`
	EndDate        *time.Time `json:"end_date"`
	CompanyID      *uint      `json:"company_id"`
	FactoryID      *uint      `json:"factory_id"`
	ChannelID      *uint      `json:"channel_id"`
	ProductID      *uint      `json:"product_id"`
	Status         *string    `json:"status"`
	Category       *string    `json:"category"`
	Limit          int        `json:"limit"`
	Offset         int        `json:"offset"`
	SortBy         string     `json:"sort_by"`
	SortOrder      string     `json:"sort_order"`
}

// ParseFilterParams parses filter parameters from query string
func ParseFilterParams(c echo.Context) FilterParams {
	params := FilterParams{
		Limit:     100,
		Offset:    0,
		SortBy:    "created_at",
		SortOrder: "desc",
	}
	
	// Parse date range
	if startDateStr := c.QueryParam("start_date"); startDateStr != "" {
		if t, err := time.Parse("2006-01-02", startDateStr); err == nil {
			params.StartDate = &t
		}
	}
	if endDateStr := c.QueryParam("end_date"); endDateStr != "" {
		if t, err := time.Parse("2006-01-02", endDateStr); err == nil {
			params.EndDate = &t
		}
	}
	
	// Parse IDs
	if companyIDStr := c.QueryParam("company_id"); companyIDStr != "" {
		if id, err := strconv.ParseUint(companyIDStr, 10, 32); err == nil {
			companyID := uint(id)
			params.CompanyID = &companyID
		}
	}
	if factoryIDStr := c.QueryParam("factory_id"); factoryIDStr != "" {
		if id, err := strconv.ParseUint(factoryIDStr, 10, 32); err == nil {
			factoryID := uint(id)
			params.FactoryID = &factoryID
		}
	}
	if channelIDStr := c.QueryParam("channel_id"); channelIDStr != "" {
		if id, err := strconv.ParseUint(channelIDStr, 10, 32); err == nil {
			channelID := uint(id)
			params.ChannelID = &channelID
		}
	}
	if productIDStr := c.QueryParam("product_id"); productIDStr != "" {
		if id, err := strconv.ParseUint(productIDStr, 10, 32); err == nil {
			productID := uint(id)
			params.ProductID = &productID
		}
	}
	
	// Parse status and category
	if status := c.QueryParam("status"); status != "" {
		params.Status = &status
	}
	if category := c.QueryParam("category"); category != "" {
		params.Category = &category
	}
	
	// Parse pagination
	if limitStr := c.QueryParam("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err == nil && limit > 0 {
			params.Limit = limit
		}
	}
	if offsetStr := c.QueryParam("offset"); offsetStr != "" {
		if offset, err := strconv.Atoi(offsetStr); err == nil && offset >= 0 {
			params.Offset = offset
		}
	}
	
	// Parse sorting
	if sortBy := c.QueryParam("sort_by"); sortBy != "" {
		params.SortBy = sortBy
	}
	if sortOrder := c.QueryParam("sort_order"); sortOrder != "" {
		sortOrder = strings.ToLower(sortOrder)
		if sortOrder == "asc" || sortOrder == "desc" {
			params.SortOrder = sortOrder
		}
	}
	
	return params
}

// GetFilteredSalesData returns filtered sales data with caching
func GetFilteredSalesData(c echo.Context) error {
	params := ParseFilterParams(c)
	
	// Generate cache key
	cacheKey := utils.GenerateCacheKey("/api/sales/filtered", map[string]interface{}{
		"start_date": params.StartDate,
		"end_date":   params.EndDate,
		"channel_id": params.ChannelID,
		"limit":      params.Limit,
		"offset":     params.Offset,
	})
	
	// Try to get from cache
	if cached, found := utils.GetFromCache(database.DB, cacheKey); found {
		return c.JSONBlob(http.StatusOK, []byte(cached))
	}
	
	// Build query
	query := database.DB.Table("channelwise_monthly_report")
	
	if params.StartDate != nil {
		query = query.Where("data_month >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("data_month <= ?", *params.EndDate)
	}
	if params.ChannelID != nil {
		query = query.Where("channel_id = ?", *params.ChannelID)
	}
	
	// Get total count
	var total int64
	query.Count(&total)
	
	// Apply sorting and pagination
	orderBy := params.SortBy + " " + params.SortOrder
	var results []map[string]interface{}
	query.Order(orderBy).Limit(params.Limit).Offset(params.Offset).Find(&results)
	
	response := map[string]interface{}{
		"data":   results,
		"total":  total,
		"limit":  params.Limit,
		"offset": params.Offset,
	}
	
	responseJSON, _ := json.Marshal(response)
	
	// Cache for 5 minutes
	utils.SetCache(database.DB, cacheKey, "/api/sales/filtered", map[string]interface{}{
		"start_date": params.StartDate,
		"end_date":   params.EndDate,
		"channel_id": params.ChannelID,
	}, string(responseJSON), 5*time.Minute)
	
	return c.JSON(http.StatusOK, response)
}

// GetFilteredProductionData returns filtered production data
func GetFilteredProductionData(c echo.Context) error {
	params := ParseFilterParams(c)
	
	cacheKey := utils.GenerateCacheKey("/api/production/filtered", map[string]interface{}{
		"start_date": params.StartDate,
		"end_date":   params.EndDate,
		"factory_id": params.FactoryID,
		"limit":      params.Limit,
		"offset":     params.Offset,
	})
	
	if cached, found := utils.GetFromCache(database.DB, cacheKey); found {
		return c.JSONBlob(http.StatusOK, []byte(cached))
	}
	
	query := database.DB.Table("production_analyses")
	
	if params.StartDate != nil {
		query = query.Where("month >= ?", *params.StartDate)
	}
	if params.EndDate != nil {
		query = query.Where("month <= ?", *params.EndDate)
	}
	if params.FactoryID != nil {
		query = query.Where("factory_id = ?", *params.FactoryID)
	}
	
	var total int64
	query.Count(&total)
	
	orderBy := params.SortBy + " " + params.SortOrder
	var results []map[string]interface{}
	query.Order(orderBy).Limit(params.Limit).Offset(params.Offset).Find(&results)
	
	response := map[string]interface{}{
		"data":   results,
		"total":  total,
		"limit":  params.Limit,
		"offset": params.Offset,
	}
	
	responseJSON, _ := json.Marshal(response)
	utils.SetCache(database.DB, cacheKey, "/api/production/filtered", map[string]interface{}{
		"start_date": params.StartDate,
		"end_date":   params.EndDate,
		"factory_id": params.FactoryID,
	}, string(responseJSON), 5*time.Minute)
	
	return c.JSON(http.StatusOK, response)
}

// GetAdvancedFilters returns available filter options for a report type
func GetAdvancedFilters(c echo.Context) error {
	reportType := c.Param("report_type")
	
	filters := map[string]interface{}{
		"date_range": map[string]interface{}{
			"type":        "date_range",
			"label":       "Date Range",
			"required":    true,
			"default":     "last_30_days",
			"options":     []string{"last_7_days", "last_30_days", "last_90_days", "this_month", "this_year", "custom"},
		},
		"company": map[string]interface{}{
			"type":     "multi_select",
			"label":    "Company",
			"required": false,
			"source":   "/api/companies",
		},
		"factory": map[string]interface{}{
			"type":     "multi_select",
			"label":    "Factory",
			"required": false,
			"source":   "/api/factories",
			"depends_on": "company",
		},
		"channel": map[string]interface{}{
			"type":     "multi_select",
			"label":    "Channel",
			"required": false,
			"source":   "/api/channels",
		},
		"product": map[string]interface{}{
			"type":     "multi_select",
			"label":    "Product",
			"required": false,
			"source":   "/api/products",
		},
		"status": map[string]interface{}{
			"type":     "multi_select",
			"label":    "Status",
			"required": false,
			"options":  []string{"active", "inactive", "pending", "completed"},
		},
	}
	
	// Customize filters based on report type
	switch reportType {
	case "sales":
		filters["channel"] = map[string]interface{}{
			"type":     "multi_select",
			"label":    "Sales Channel",
			"required": false,
			"source":   "/api/channels",
		}
	case "production":
		filters["factory"] = map[string]interface{}{
			"type":     "multi_select",
			"label":    "Factory",
			"required": false,
			"source":   "/api/factories",
		}
		filters["production_line"] = map[string]interface{}{
			"type":     "multi_select",
			"label":    "Production Line",
			"required": false,
			"source":   "/api/production-lines",
		}
	case "finance":
		filters["budget_category"] = map[string]interface{}{
			"type":     "multi_select",
			"label":    "Budget Category",
			"required": false,
			"source":   "/api/budget-categories",
		}
		filters["department"] = map[string]interface{}{
			"type":     "multi_select",
			"label":    "Department",
			"required": false,
			"source":   "/api/departments",
		}
	}
	
	return c.JSON(http.StatusOK, map[string]interface{}{
		"report_type": reportType,
		"filters":     filters,
	})
}

