package models

import (
	"time"
	"gorm.io/gorm"
)

// SupplyChainRawData represents supply chain raw data
type SupplyChainRawData struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Company   int            `gorm:"type:int" json:"company"`
	PONumber  string         `gorm:"type:varchar(255)" json:"po_number"`
	PODate    time.Time      `gorm:"type:date" json:"po_date"`
	POValue   float64        `gorm:"type:decimal(15,2)" json:"po_value"`
	Status    string         `gorm:"type:varchar(50)" json:"status"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SupplyChainRawData) TableName() string {
	return "supply_chain_raw_datas"
}

// SupplyChainMasterData represents supply chain master data
type SupplyChainMasterData struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Company     int            `gorm:"type:int" json:"company"`
	PONumber    string         `gorm:"type:varchar(255)" json:"po_number"`
	PODate      time.Time      `gorm:"type:date" json:"po_date"`
	POValue     float64        `gorm:"type:decimal(15,2)" json:"po_value"`
	PRAmount    float64        `gorm:"type:decimal(15,2)" json:"pr_amount"`
	PurchaseOrg string         `gorm:"type:varchar(255)" json:"purchase_org"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SupplyChainMasterData) TableName() string {
	return "supply_chain_master_datas"
}

// SupplyChainGrnData represents GRN (Goods Receipt Note) data
type SupplyChainGrnData struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Company   int            `gorm:"type:int" json:"company"`
	PONumber  string         `gorm:"type:varchar(255)" json:"po_number"`
	GRNDate   time.Time      `gorm:"type:date" json:"grn_date"`
	GRNAmount float64        `gorm:"type:decimal(15,2)" json:"grn_amount"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SupplyChainGrnData) TableName() string {
	return "supply_chain_grn_datas"
}

// SupplyChainInvoiceData represents invoice data
type SupplyChainInvoiceData struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Company     int            `gorm:"type:int" json:"company"`
	InvoiceNumber string       `gorm:"type:varchar(255)" json:"invoice_number"`
	IVDate      time.Time      `gorm:"type:date" json:"iv_date"`
	TotalInvoice float64       `gorm:"type:decimal(15,2)" json:"total_invoice"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (SupplyChainInvoiceData) TableName() string {
	return "supply_chain_invoice_datas"
}

