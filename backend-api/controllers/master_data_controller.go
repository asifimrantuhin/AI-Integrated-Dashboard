package controllers

import (
	"idash-backend-api/database"
	"idash-backend-api/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GetChannels(c echo.Context) error {
	var channels []models.Channel
	database.DB.Where("status = ?", 1).Find(&channels)
	return c.JSON(http.StatusOK, channels)
}

func GetPlantGroups(c echo.Context) error {
	type PlantGroup struct {
		ID   uint   `json:"id"`
		Name string `json:"name"`
	}

	var plantGroups []PlantGroup
	database.DB.Raw("SELECT id, name FROM plant_groups WHERE status = 1").Scan(&plantGroups)
	return c.JSON(http.StatusOK, plantGroups)
}

func GetCompanyPlants(c echo.Context) error {
	type CompanyPlant struct {
		CompanyID   int    `json:"company_id"`
		CompanyName string `json:"company_name"`
		PlantID     int    `json:"plant_id"`
		PlantName   string `json:"plant_name"`
	}

	var companyPlants []CompanyPlant
	query := `
		SELECT 
			cf.company_id,
			c.name as company_name,
			cf.plant_id,
			p.name as plant_name
		FROM company_factories cf
		JOIN tbl_company c ON cf.company_id = c.id
		JOIN plant_groups p ON cf.plant_id = p.id
		WHERE cf.status = 1
	`
	database.DB.Raw(query).Scan(&companyPlants)
	return c.JSON(http.StatusOK, companyPlants)
}

func GetSidebarData(c echo.Context) error {
	month := c.QueryParam("month")
	
	type SidebarData struct {
		TopDistributor   interface{} `json:"top_distributor"`
		BestSellingPG    interface{} `json:"best_selling_pg"`
		BestSellingProduct interface{} `json:"best_selling_product"`
	}

	var sidebarData SidebarData

	// Get top distributor
	type TopDistributor struct {
		DBName string  `json:"db_name"`
		Amount float64 `json:"amount"`
	}
	var topDistributor TopDistributor
	database.DB.Raw(`
		SELECT db_name, amount
		FROM top_channel_d_bs
		WHERE type = 0 AND date = ?
		ORDER BY amount DESC
		LIMIT 1
	`, month).Scan(&topDistributor)
	sidebarData.TopDistributor = topDistributor

	// Get best selling PG
	type BestSellingPG struct {
		CategoryName string  `json:"category_name"`
		Value        float64 `json:"value"`
	}
	var bestSellingPG BestSellingPG
	database.DB.Raw(`
		SELECT category_name, value
		FROM best_selling_pgs
		WHERE year_month = ?
		ORDER BY value DESC
		LIMIT 1
	`, month).Scan(&bestSellingPG)
	sidebarData.BestSellingPG = bestSellingPG

	// Get best selling product
	type BestSellingProduct struct {
		ProductName string  `json:"product_name"`
		Value       float64 `json:"value"`
	}
	var bestSellingProduct BestSellingProduct
	database.DB.Raw(`
		SELECT product_name, value
		FROM best_selling_products
		WHERE year_month = ?
		ORDER BY value DESC
		LIMIT 1
	`, month).Scan(&bestSellingProduct)
	sidebarData.BestSellingProduct = bestSellingProduct

	return c.JSON(http.StatusOK, sidebarData)
}

