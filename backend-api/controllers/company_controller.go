package controllers

import (
	"idash-backend-api/models"
	"idash-backend-api/database"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetCompanies(c echo.Context) error {
	var companies []models.Company
	if err := database.DB.Find(&companies).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to fetch companies"})
	}
	return c.JSON(http.StatusOK, companies)
}

func GetCompany(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var company models.Company
	if err := database.DB.First(&company, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Company not found"})
	}
	return c.JSON(http.StatusOK, company)
}

func CreateCompany(c echo.Context) error {
	var company models.Company
	if err := c.Bind(&company); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := database.DB.Create(&company).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create company"})
	}

	return c.JSON(http.StatusCreated, company)
}

func UpdateCompany(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	var company models.Company
	if err := database.DB.First(&company, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Company not found"})
	}

	if err := c.Bind(&company); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if err := database.DB.Save(&company).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update company"})
	}

	return c.JSON(http.StatusOK, company)
}

func DeleteCompany(c echo.Context) error {
	id, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	if err := database.DB.Delete(&models.Company{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete company"})
	}
	return c.JSON(http.StatusOK, map[string]string{"message": "Company deleted successfully"})
}

