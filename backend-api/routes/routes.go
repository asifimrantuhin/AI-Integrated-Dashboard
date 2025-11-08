package routes

import (
	"idash-backend-api/controllers"
	"idash-backend-api/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo) {
	// Public routes
	e.POST("/api/auth/login", controllers.Login)
	e.POST("/api/auth/register", controllers.Register)

	// Protected routes
	api := e.Group("/api", middleware.AuthMiddleware())
	{
		// Auth routes
		api.GET("/auth/user", controllers.GetUser)

		// Companies
		api.GET("/companies", controllers.GetCompanies)
		api.GET("/companies/:id", controllers.GetCompany)
		api.POST("/companies", controllers.CreateCompany)
		api.PUT("/companies/:id", controllers.UpdateCompany)
		api.DELETE("/companies/:id", controllers.DeleteCompany)

		// Master Data
		api.GET("/channels", controllers.GetChannels)
		api.GET("/plant-groups", controllers.GetPlantGroups)
		api.GET("/company-plants", controllers.GetCompanyPlants)
		api.GET("/sidebar-data", controllers.GetSidebarData)

		// Sales routes
		sales := api.Group("/sales")
		{
			sales.GET("/overview", controllers.GetSalesOverview)
			sales.GET("/summary", controllers.GetSalesSummary)
			sales.GET("/cumulative", controllers.GetSalesCumulative)
			sales.GET("/channelwise", controllers.GetChannelwiseSales)
			sales.GET("/daily", controllers.GetDailySales)
			sales.GET("/best-selling-products", controllers.GetBestSellingProducts)
			sales.GET("/best-selling-pgs", controllers.GetBestSellingPGs)
			sales.GET("/slow-moving-products", controllers.GetSlowMovingProducts)
			sales.GET("/top-distributors", controllers.GetTopDistributors)
			sales.GET("/top-retailers", controllers.GetTopRetailers)
			sales.GET("/order-vs-delivery", controllers.GetOrderVsDelivery)
		}

		// Production routes
		production := api.Group("/production")
		{
			production.GET("/overview", controllers.GetProductionOverview)
			production.GET("/summary", controllers.GetProductionSummary)
			production.GET("/analysis", controllers.GetProductionAnalysis)
			production.GET("/last-months", controllers.GetProductionLastMonths)
			production.GET("/wastage", controllers.GetWastageData)
			production.GET("/cost", controllers.GetCostAnalysis)
			production.GET("/trends", controllers.GetProductionTrends)
		}

		// Finance routes
		finance := api.Group("/finance")
		{
			finance.GET("/overview", controllers.GetFinanceOverview)
			finance.GET("/summary", controllers.GetFinanceSummary)
			finance.GET("/budget", controllers.GetBudget)
			finance.GET("/budget-vs-actual", controllers.GetBudgetVsActual)
			finance.GET("/bank-loan", controllers.GetBankLoan)
			finance.GET("/bank-loan-status", controllers.GetBankLoanStatus)
			finance.GET("/company-bank-loan", controllers.GetCompanyBankLoanSummary)
			finance.GET("/expenses", controllers.GetExpenses)
			finance.GET("/budget-categories", controllers.GetBudgetCategories)
			finance.GET("/budget-departments", controllers.GetBudgetDepartments)
		}

		// Inventory routes
		inventory := api.Group("/inventory")
		{
			inventory.GET("/overview", controllers.GetInventoryOverview)
			inventory.GET("/summary", controllers.GetInventorySummary)
			inventory.GET("/ratio", controllers.GetInventoryRatio)
			inventory.GET("/valuation", controllers.GetInventoryValuation)
			inventory.GET("/cogs-gp", controllers.GetCOGSGP)
			inventory.GET("/company-summary", controllers.GetCompanyInventorySummary)
			inventory.GET("/trends", controllers.GetInventoryTrends)
			inventory.GET("/gl-accounts", controllers.GetInventoryGLAccounts)
		}

		// HR routes
		hr := api.Group("/hr")
		{
			hr.GET("/overview", controllers.GetHROverview)
			hr.GET("/summary", controllers.GetHRSummary)
			hr.GET("/employees", controllers.GetEmployees)
			hr.GET("/attendance", controllers.GetEmployeeAttendance)
			hr.GET("/promotions", controllers.GetEmployeePromotions)
			hr.GET("/department", controllers.GetDepartmentAnalytics)
			hr.GET("/tranover", controllers.GetEmployeeTranOver)
			hr.GET("/tranover-summary", controllers.GetEmployeeTranOverSummary)
		}

		// Supply Chain routes
		supplychain := api.Group("/supplychain")
		{
			supplychain.GET("/summary", controllers.GetSupplyChainSummary)
			supplychain.GET("/grn", controllers.GetGRNData)
			supplychain.GET("/invoice", controllers.GetInvoiceData)
			supplychain.GET("/po", controllers.GetPOManagement)
			supplychain.GET("/pending", controllers.GetPendingItems)
			supplychain.GET("/company-po-monthly", controllers.GetCompanyPOMonthly)
			supplychain.GET("/company-grn-monthly", controllers.GetCompanyGRNMonthly)
			supplychain.GET("/company-invoice-monthly", controllers.GetCompanyInvoiceMonthly)
		}

		// NPD routes
		npd := api.Group("/npd")
		{
			npd.GET("/projects", controllers.GetNPDProjects)
			npd.GET("/deliverables", controllers.GetProjectDeliverables)
			npd.GET("/sub-deliverables", controllers.GetProjectSubDeliverables)
			npd.GET("/status", controllers.GetProjectStatus)
		}

		// Dashboard routes
		dashboard := api.Group("/dashboard")
		{
			dashboard.GET("/executive", controllers.GetExecutiveDashboard, middleware.RequireRoles("executive"))
			dashboard.GET("/finance", controllers.GetFinanceDashboard, middleware.RequireRoles("executive", "finance_manager"))
			dashboard.GET("/sales", controllers.GetSalesDashboard, middleware.RequireRoles("executive", "sales_manager"))
			dashboard.GET("/production", controllers.GetProductionDashboard, middleware.RequireRoles("executive", "production_manager"))
			dashboard.GET("/hr", controllers.GetHRDashboard, middleware.RequireRoles("executive", "hr_manager"))
			dashboard.GET("/supplychain", controllers.GetSupplyChainDashboard, middleware.RequireRoles("executive", "supply_chain_manager"))
		}

		// AI/Forecasting routes
		ai := api.Group("/ai")
		{
			ai.POST("/forecast/sales", controllers.RequestSalesForecast)
			ai.POST("/forecast/production", controllers.RequestProductionForecast)
			ai.POST("/forecast/finance", controllers.RequestFinancialForecast)
			ai.POST("/forecast/inventory", controllers.RequestInventoryForecast)
			ai.GET("/forecast/:id", controllers.GetForecast)
			ai.GET("/forecasts/:type", controllers.GetForecastsByType)
			ai.POST("/analyze", controllers.AnalyzeData)
			ai.GET("/predict/sales", controllers.PredictSalesSummary)
			ai.GET("/prescribe/inventory", controllers.GetInventoryPrescription)
			ai.GET("/prescribe/finance", controllers.GetFinancialPrescription)
			ai.GET("/anomaly", controllers.AnalyzeAnomaliesWithActions)
			ai.POST("/scenario", controllers.RunScenarioSimulation)
		}

		// BI / Executive Intelligence routes
		bi := api.Group("/bi", middleware.RequireRoles("executive", "analyst"))
		{
			bi.GET("/overview", controllers.GetBIOverview)
			bi.POST("/assistant", controllers.PostBIAssistant)
		}

		// Admin & user management routes
		admin := api.Group("/admin", middleware.RequirePermissions("manage_users"))
		{
			admin.GET("/users", controllers.ListUsers)
			admin.GET("/users/:id", controllers.GetUserDetail)
			admin.POST("/users", controllers.CreateUser)
			admin.PUT("/users/:id", controllers.UpdateUser)
			admin.DELETE("/users/:id", controllers.DeleteUser)
			admin.GET("/roles", controllers.ListRoles)
			admin.GET("/permissions", controllers.ListPermissions)
		}

		// Manufacturing routes
		manufacturing := api.Group("/manufacturing")
		{
			manufacturing.GET("/quality", controllers.GetQualityControlData)
			manufacturing.GET("/production-plans", controllers.GetProductionPlans)
			manufacturing.GET("/machine-maintenance", controllers.GetMachineMaintenance)
			manufacturing.GET("/efficiency", controllers.GetProductionEfficiency)
			manufacturing.GET("/supplier-performance", controllers.GetSupplierPerformance)
			manufacturing.GET("/energy-consumption", controllers.GetEnergyConsumption)
		}

		// Filter routes
		filter := api.Group("/filter")
		{
			filter.GET("/sales", controllers.GetFilteredSalesData)
			filter.GET("/production", controllers.GetFilteredProductionData)
			filter.GET("/filters/:report_type", controllers.GetAdvancedFilters)
		}

		// Business Calculator routes
		business := api.Group("/business")
		{
			business.GET("/calculation", controllers.GetBusinessCalculation)
			business.GET("/primary-calculation", controllers.GetPrimaryBusinessCalculation)
			business.GET("/secondary-calculation", controllers.GetSecondaryBusinessCalculation)
			business.GET("/primary-intelligence", controllers.GetPrimaryIntelligenceDetails)
			business.GET("/secondary-intelligence", controllers.GetSecondaryIntelligenceDetails)
			business.GET("/top-retailer-list", controllers.GetTopRetailerList)
		}

		// Integration / ETL routes
		integration := api.Group("/integration", middleware.RequireRoles("executive"))
		{
			integration.POST("/sync", controllers.TriggerIntegrationSync)
			integration.GET("/status", controllers.GetIntegrationStatus)
		}

		// Reporting routes
		reports := api.Group("/reports", middleware.RequireRoles("executive", "analyst"))
		{
			reports.GET("", controllers.ListReports)
			reports.GET("/:id", controllers.GenerateReport)
		}
	}
}
