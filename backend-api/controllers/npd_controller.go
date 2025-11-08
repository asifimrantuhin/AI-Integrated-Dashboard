package controllers

import (
	"idash-backend-api/database"
	"idash-backend-api/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func GetNPDProjects(c echo.Context) error {
	status := c.QueryParam("status")

	var projects []models.NpdProjects
	query := database.DB

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Find(&projects)

	return c.JSON(http.StatusOK, projects)
}

func GetProjectDeliverables(c echo.Context) error {
	projectID := c.QueryParam("project_id")
	status := c.QueryParam("status")

	if projectID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "project_id is required"})
	}

	projID, _ := strconv.Atoi(projectID)

	var deliverables []models.ProjectsDeliverables
	query := database.DB.Where("project_id = ?", projID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Find(&deliverables)

	return c.JSON(http.StatusOK, deliverables)
}

func GetProjectSubDeliverables(c echo.Context) error {
	deliverableID := c.QueryParam("deliverable_id")
	status := c.QueryParam("status")

	if deliverableID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "deliverable_id is required"})
	}

	delID, _ := strconv.Atoi(deliverableID)

	var subDeliverables []models.ProjectsSubDeliverables
	query := database.DB.Where("deliverable_id = ?", delID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	query.Find(&subDeliverables)

	return c.JSON(http.StatusOK, subDeliverables)
}

func GetProjectStatus(c echo.Context) error {
	type ProjectStatus struct {
		TotalProjects     int `json:"total_projects"`
		ActiveProjects    int `json:"active_projects"`
		CompletedProjects int `json:"completed_projects"`
		PendingProjects   int `json:"pending_projects"`
	}

	var status ProjectStatus

	database.DB.Model(&models.NpdProjects{}).Count(&status.TotalProjects)
	database.DB.Model(&models.NpdProjects{}).Where("status = ?", "active").Count(&status.ActiveProjects)
	database.DB.Model(&models.NpdProjects{}).Where("status = ?", "completed").Count(&status.CompletedProjects)
	database.DB.Model(&models.NpdProjects{}).Where("status = ?", "pending").Count(&status.PendingProjects)

	return c.JSON(http.StatusOK, status)
}

