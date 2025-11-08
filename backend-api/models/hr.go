package models

import "time"

type EmployeeBasicInfo struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EmployeeID string    `json:"employee_id"`
	Name      string    `json:"name"`
	Department string    `json:"department"`
	CompanyID int       `json:"company_id"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (EmployeeBasicInfo) TableName() string {
	return "employee_basic_infos"
}

type EmployeeAttendance struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EmployeeID string    `json:"employee_id"`
	Date      time.Time `gorm:"type:date" json:"date"`
	Status    string    `json:"status"` // e.g., Present, Absent, Leave
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (EmployeeAttendance) TableName() string {
	return "employee_attendances"
}

type YearlyEmployeePromotion struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	EmployeeID    string    `json:"employee_id"`
	Year          int       `json:"year"`
	PromotionDate time.Time `gorm:"type:date" json:"promotion_date"`
	NewDesignation string   `json:"new_designation"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (YearlyEmployeePromotion) TableName() string {
	return "yearly_employee_promotions"
}

type EmployeeTranOver struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	EmployeeID string    `json:"employee_id"`
	Date      time.Time `gorm:"type:date" json:"date"`
	Type      string    `json:"type"` // transfer or overtime
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (EmployeeTranOver) TableName() string {
	return "employee_tran_overs"
}
