package models

import "time"

type InventoryRawData struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CompanyID int       `json:"company_id"`
	GLID      int       `json:"gl_id"`
	Month     time.Time `gorm:"type:date" json:"month"`
	Amount    float64   `json:"amount"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (InventoryRawData) TableName() string {
	return "inventory_raw_datas"
}

type CogsGp struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CompanyID int       `json:"company_id"`
	Month     time.Time `gorm:"type:date" json:"month"`
	COGS      float64   `json:"cogs"`
	GP        float64   `json:"gp"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (CogsGp) TableName() string {
	return "cogs_gps"
}

type InventoryGlAccounts struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	GLCode    string    `json:"gl_code"`
	GLName    string    `json:"gl_name"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (InventoryGlAccounts) TableName() string {
	return "inventory_gl_accounts"
}
