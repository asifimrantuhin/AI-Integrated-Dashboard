package models

import (
	"time"
	"gorm.io/gorm"
)

// ProductionAnalysis represents production analysis data
type ProductionAnalysis struct {
	ID             uint           `gorm:"primaryKey" json:"id"`
	Month          time.Time      `gorm:"type:date" json:"month"`
	Factory        string         `gorm:"type:varchar(255)" json:"factory"`
	SummaryGroup   string         `gorm:"type:varchar(255)" json:"summary_group"`
	CMonthlyAmount float64        `gorm:"type:decimal(15,2)" json:"cmonthly_amount"`
	CAvgAmount     float64        `gorm:"type:decimal(15,2)" json:"cavg_amount"`
	CMonthlyPer    float64        `gorm:"type:decimal(10,2)" json:"cmonthly_per"`
	CAvgPer        float64        `gorm:"type:decimal(10,2)" json:"cavg_per"`
	PMonthlyAmount float64        `gorm:"type:decimal(15,2)" json:"pmonthly_amount"`
	PAvgAmount     float64        `gorm:"type:decimal(15,2)" json:"pavg_amount"`
	PMonthlyPer    float64        `gorm:"type:decimal(10,2)" json:"pmonthly_per"`
	PAvgPer        float64        `gorm:"type:decimal(10,2)" json:"pavg_per"`
	TMonthlyAmount float64        `gorm:"type:decimal(15,2)" json:"tmonthly_amount"`
	TMonthlyPer    float64        `gorm:"type:decimal(10,2)" json:"tmonthly_per"`
	AMonthlyAmount float64        `gorm:"type:decimal(15,2)" json:"amonthly_amount"`
	AMonthlyPer    float64        `gorm:"type:decimal(10,2)" json:"amonthly_per"`
	AAvgAmount     float64        `gorm:"type:decimal(15,2)" json:"aavg_amount"`
	CreatedAt      time.Time      `json:"created_at"`
	UpdatedAt      time.Time      `json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}

func (ProductionAnalysis) TableName() string {
	return "production_analyses"
}

// WastageData represents wastage data
type WastageData struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Month     time.Time      `gorm:"type:date" json:"month"`
	Factory   string         `gorm:"type:varchar(255)" json:"factory"`
	Wastage   float64        `gorm:"type:decimal(15,2)" json:"wastage"`
	Amount    float64        `gorm:"type:decimal(15,2)" json:"amount"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (WastageData) TableName() string {
	return "wastage_datas"
}

// CostAnalysis represents cost analysis data
type CostAnalysis struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Month     time.Time      `gorm:"type:date" json:"month"`
	Factory   string         `gorm:"type:varchar(255)" json:"factory"`
	CostType  string         `gorm:"type:varchar(255)" json:"cost_type"`
	Amount    float64        `gorm:"type:decimal(15,2)" json:"amount"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CostAnalysis) TableName() string {
	return "cost_analyses"
}

