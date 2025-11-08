package models

import (
	"time"
	"gorm.io/gorm"
)

// QualityControlCheck represents quality control data
type QualityControlCheck struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CompanyID       uint           `gorm:"index" json:"company_id"`
	FactoryID       uint           `gorm:"index" json:"factory_id"`
	ProductionLineID *uint         `gorm:"index" json:"production_line_id"`
	ProductCode     string         `gorm:"type:varchar(100);index" json:"product_code"`
	ProductName     string         `gorm:"type:varchar(255)" json:"product_name"`
	CheckDate       time.Time      `gorm:"type:date;index" json:"check_date"`
	CheckType       string         `gorm:"type:varchar(50);index" json:"check_type"` // incoming, in_process, final, sampling
	Status          string         `gorm:"type:varchar(50);index" json:"status"` // passed, failed, pending, rework
	TotalChecked    int            `json:"total_checked"`
	PassedCount     int            `json:"passed_count"`
	FailedCount     int            `json:"failed_count"`
	DefectRate      float64        `gorm:"type:decimal(5,2)" json:"defect_rate"`
	Defects         string         `gorm:"type:text" json:"defects"`
	Remarks         string         `gorm:"type:text" json:"remarks"`
	InspectorName   string         `gorm:"type:varchar(255)" json:"inspector_name"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (QualityControlCheck) TableName() string {
	return "quality_control_checks"
}

// ProductionPlan represents production planning data
type ProductionPlan struct {
	ID            uint           `gorm:"primaryKey" json:"id"`
	CompanyID     uint           `gorm:"index" json:"company_id"`
	FactoryID     uint           `gorm:"index" json:"factory_id"`
	PlanNumber    string         `gorm:"type:varchar(100);uniqueIndex" json:"plan_number"`
	ProductCode   string         `gorm:"type:varchar(100);index" json:"product_code"`
	ProductName   string         `gorm:"type:varchar(255)" json:"product_name"`
	PlanDate      time.Time      `gorm:"type:date;index" json:"plan_date"`
	StartDate     time.Time      `gorm:"type:date;index" json:"start_date"`
	EndDate       time.Time      `gorm:"type:date;index" json:"end_date"`
	PlannedQuantity float64      `gorm:"type:decimal(15,2)" json:"planned_quantity"`
	ActualQuantity   float64      `gorm:"type:decimal(15,2)" json:"actual_quantity"`
	Status        string         `gorm:"type:varchar(50);index" json:"status"` // draft, approved, in_progress, completed, cancelled
	Priority      string         `gorm:"type:varchar(50);index;default:'medium'" json:"priority"` // low, medium, high, urgent
	Notes         string         `gorm:"type:text" json:"notes"`
	CreatedBy     *uint          `json:"created_by"`
	ApprovedBy    *uint          `json:"approved_by"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ProductionPlan) TableName() string {
	return "production_plans"
}

// MachineMaintenance represents machine maintenance data
type MachineMaintenance struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CompanyID       uint           `gorm:"index" json:"company_id"`
	FactoryID       uint           `gorm:"index" json:"factory_id"`
	MachineCode     string         `gorm:"type:varchar(100);index" json:"machine_code"`
	MachineName     string         `gorm:"type:varchar(255)" json:"machine_name"`
	MaintenanceType string         `gorm:"type:varchar(50);index" json:"maintenance_type"` // preventive, corrective, breakdown, scheduled
	MaintenanceDate time.Time      `gorm:"type:date;index" json:"maintenance_date"`
	StartTime       *time.Time     `json:"start_time"`
	EndTime         *time.Time     `json:"end_time"`
	DowntimeMinutes int            `json:"downtime_minutes"`
	Description     string         `gorm:"type:text" json:"description"`
	ActionsTaken    string         `gorm:"type:text" json:"actions_taken"`
	Cost            float64        `gorm:"type:decimal(15,2)" json:"cost"`
	Status          string         `gorm:"type:varchar(50);index" json:"status"` // scheduled, in_progress, completed, cancelled
	TechnicianName  string         `gorm:"type:varchar(255)" json:"technician_name"`
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (MachineMaintenance) TableName() string {
	return "machine_maintenances"
}

// MaterialRequirement represents MRP data
type MaterialRequirement struct {
	ID              uint           `gorm:"primaryKey" json:"id"`
	CompanyID       uint           `gorm:"index" json:"company_id"`
	ProductionPlanID *uint         `gorm:"index" json:"production_plan_id"`
	MaterialCode    string         `gorm:"type:varchar(100);index" json:"material_code"`
	MaterialName    string         `gorm:"type:varchar(255)" json:"material_name"`
	MaterialType    string         `gorm:"type:varchar(100);index" json:"material_type"`
	Unit            string         `gorm:"type:varchar(50)" json:"unit"`
	RequiredQuantity float64       `gorm:"type:decimal(15,2)" json:"required_quantity"`
	AvailableQuantity float64      `gorm:"type:decimal(15,2)" json:"available_quantity"`
	ShortageQuantity  float64      `gorm:"type:decimal(15,2)" json:"shortage_quantity"`
	RequiredDate    time.Time      `gorm:"type:date;index" json:"required_date"`
	Status          string         `gorm:"type:varchar(50);index" json:"status"` // pending, ordered, received, shortage
	CreatedAt       time.Time      `json:"created_at"`
	UpdatedAt       time.Time      `json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
}

func (MaterialRequirement) TableName() string {
	return "material_requirements"
}

// ProductionEfficiency represents production efficiency metrics
type ProductionEfficiency struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	CompanyID         uint      `gorm:"index" json:"company_id"`
	FactoryID         uint      `gorm:"index" json:"factory_id"`
	ProductionLineID  *uint     `gorm:"index" json:"production_line_id"`
	ProductionDate    time.Time `gorm:"type:date;index" json:"production_date"`
	Shift             string    `gorm:"type:varchar(50);index" json:"shift"`
	PlannedOutput     float64   `gorm:"type:decimal(15,2)" json:"planned_output"`
	ActualOutput      float64   `gorm:"type:decimal(15,2)" json:"actual_output"`
	EfficiencyPercentage float64 `gorm:"type:decimal(5,2);index" json:"efficiency_percentage"`
	PlannedHours      int       `json:"planned_hours"`
	ActualHours       int       `json:"actual_hours"`
	DowntimeMinutes   int       `json:"downtime_minutes"`
	OEE               float64   `gorm:"type:decimal(5,2)" json:"oee"` // Overall Equipment Effectiveness
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (ProductionEfficiency) TableName() string {
	return "production_efficiency"
}

// SupplierPerformance represents supplier performance data
type SupplierPerformance struct {
	ID                uint      `gorm:"primaryKey" json:"id"`
	CompanyID         uint      `gorm:"index" json:"company_id"`
	SupplierCode      string    `gorm:"type:varchar(100);index" json:"supplier_code"`
	SupplierName      string    `gorm:"type:varchar(255)" json:"supplier_name"`
	EvaluationDate    time.Time `gorm:"type:date;index" json:"evaluation_date"`
	TotalOrders       int       `json:"total_orders"`
	OnTimeDeliveries  int       `json:"on_time_deliveries"`
	OnTimePercentage  float64   `gorm:"type:decimal(5,2);index" json:"on_time_percentage"`
	QualityIssues     int       `json:"quality_issues"`
	QualityScore      float64   `gorm:"type:decimal(5,2);index" json:"quality_score"`
	CostScore         float64   `gorm:"type:decimal(5,2)" json:"cost_score"`
	OverallScore      float64   `gorm:"type:decimal(5,2);index" json:"overall_score"`
	Rating            string    `gorm:"type:varchar(50);index" json:"rating"` // excellent, good, average, poor
	Comments          string    `gorm:"type:text" json:"comments"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (SupplierPerformance) TableName() string {
	return "supplier_performance"
}

// EnergyConsumption represents energy consumption data
type EnergyConsumption struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	CompanyID       uint      `gorm:"index" json:"company_id"`
	FactoryID       uint      `gorm:"index" json:"factory_id"`
	ConsumptionDate time.Time `gorm:"type:date;index" json:"consumption_date"`
	EnergyType      string    `gorm:"type:varchar(50);index" json:"energy_type"` // electricity, gas, water, steam, compressed_air
	ConsumptionAmount float64 `gorm:"type:decimal(15,2)" json:"consumption_amount"`
	Unit            string    `gorm:"type:varchar(50)" json:"unit"`
	Cost            float64   `gorm:"type:decimal(15,2)" json:"cost"`
	MeterReading    string    `gorm:"type:varchar(100)" json:"meter_reading"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (EnergyConsumption) TableName() string {
	return "energy_consumption"
}

// AIForecast represents AI forecast results
type AIForecast struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	ForecastType   string    `gorm:"type:varchar(50);index" json:"forecast_type"` // sales, production, finance, inventory
	EntityType     string    `gorm:"type:varchar(50);index" json:"entity_type"` // product, channel, factory, etc.
	EntityID       string    `gorm:"type:varchar(100);index" json:"entity_id"`
	ForecastDate   time.Time `gorm:"type:date;index" json:"forecast_date"`
	ForecastedValue float64  `gorm:"type:decimal(15,2)" json:"forecasted_value"`
	ConfidenceLevel float64  `gorm:"type:decimal(5,2)" json:"confidence_level"`
	UpperBound     *float64  `gorm:"type:decimal(15,2)" json:"upper_bound"`
	LowerBound     *float64  `gorm:"type:decimal(15,2)" json:"lower_bound"`
	ModelUsed      string    `gorm:"type:varchar(100)" json:"model_used"`
	ForecastDetails string   `gorm:"type:json" json:"forecast_details"`
	Status         string    `gorm:"type:varchar(50);index;default:'active'" json:"status"` // pending, active, expired
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (AIForecast) TableName() string {
	return "ai_forecasts"
}

