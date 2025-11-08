package models

import (
	"time"
	"gorm.io/gorm"
)

// PlantGroup represents plant groups
type PlantGroup struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Name      string         `gorm:"type:varchar(255)" json:"name"`
	Status    int            `gorm:"type:int;default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (PlantGroup) TableName() string {
	return "plant_group"
}

// CompanyFactory represents company-factory relationships
type CompanyFactory struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CompanyID int            `gorm:"type:int" json:"company_id"`
	PlantID   int            `gorm:"type:int" json:"plant_id"`
	Status    int            `gorm:"type:int;default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (CompanyFactory) TableName() string {
	return "company_factories"
}

