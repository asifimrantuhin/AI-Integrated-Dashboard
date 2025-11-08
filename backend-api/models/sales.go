package models

import "time"

// ChannelwiseLICdataMonthly represents the monthly sales report data
type ChannelwiseLICdataMonthly struct {
	DataMonth         time.Time `gorm:"type:date" json:"data_month"`
	ChannelID         int       `json:"channel_id"`
	ChannelName       string    `json:"channel_name"`
	LiftingTarget     float64   `json:"lifting_target"`
	Billed            float64   `json:"billed"`
	Delivered         float64   `json:"delivered"`
	PrimaryCollection float64   `json:"primary_collection"`
	IMSTarget         float64   `json:"ims_target"`
	IMS               float64   `json:"ims"`
	MarketCollection  float64   `json:"market_collection"`
	MemoTarget        float64   `json:"memo_target"`
	MemoQty           float64   `json:"memo_qty"`
	PGTarget          float64   `json:"pg_target"`
	PGCover           float64   `json:"pg_cover"`
	TotalRetailer     int       `json:"total_retailer"`
	BusinessRetailer  int       `json:"business_retailer"`
}

// TableName specifies the table name for GORM
func (ChannelwiseLICdataMonthly) TableName() string {
	return "channelwise_monthly_report"
}

// ChannelwiseLICdata represents the daily sales report data
type ChannelwiseLICdata struct {
	DataDate          time.Time `gorm:"type:date" json:"data_date"`
	ChannelID         int       `json:"channel_id"`
	ChannelName       string    `json:"channel_name"`
	LiftingTarget     float64   `json:"lifting_target"`
	Billed            float64   `json:"billed"`
	Delivery          float64   `json:"delivery"`
	LiftingCollection float64   `json:"lifting_collection"`
	IMSTarget         float64   `json:"ims_target"`
	IMS               float64   `json:"ims"`
	IMSCollection     float64   `json:"ims_collection"`
}

// TableName specifies the table name for GORM
func (ChannelwiseLICdata) TableName() string {
	return "channelwise_lic_data"
}

// BestSellingProduct model
type BestSellingProduct struct {
	YearMonth   time.Time `gorm:"type:date" json:"year_month"`
	ChannelID   int       `json:"channel_id"`
	ProductID   int       `json:"product_id"`
	ProductName string    `json:"product_name"`
	Qty         float64   `json:"qty"`
	Value       float64   `json:"value"`
	CatID       int       `json:"cat_id"`
}

func (BestSellingProduct) TableName() string {
	return "best_selling_products"
}

// BestSellingPg model
type BestSellingPg struct {
	YearMonth    time.Time `gorm:"type:date" json:"year_month"`
	ChannelID    int       `json:"channel_id"`
	CategoryID   int       `json:"category_id"`
	CategoryName string    `json:"category_name"`
	Qty          float64   `json:"qty"`
	Value        float64   `json:"value"`
}

func (BestSellingPg) TableName() string {
	return "best_selling_pgs"
}

// TopChannelDB model
type TopChannelDB struct {
	DBName string    `json:"db_name"`
	Amount float64   `json:"amount"`
	Type   int       `json:"type"` // 0 for distributor, 1 for retailer
	Date   time.Time `gorm:"type:date" json:"date"`
}

func (TopChannelDB) TableName() string {
	return "top_channel_d_bs"
}

// OrderDeliverySummary model
type OrderDeliverySummary struct {
	Months    time.Time `gorm:"type:date" json:"months"`
	ChannelID int       `json:"channel_id"`
	Amounts   float64   `json:"amounts"`
	Types     int       `json:"types"` // 0 for order, 1 for delivery
}

func (OrderDeliverySummary) TableName() string {
	return "order_delivery_summaries"
}

// TopRetailer model
type TopRetailer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Date      time.Time `gorm:"type:date" json:"date"`
	DBName    string    `json:"db_name"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TopRetailer) TableName() string {
	return "top_retailers"
}

// Channel model
type Channel struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Channel) TableName() string {
	return "channels"
}
