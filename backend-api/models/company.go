package models

import (
	"time"
	"gorm.io/gorm"
)

type Company struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CompanyCode string       `gorm:"type:varchar(50);uniqueIndex;not null" json:"company_code"`
	CompanyName string       `gorm:"type:varchar(255);not null" json:"company_name"`
	Status    int            `gorm:"type:int;default:1" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Company) TableName() string {
	return "companies"
}

