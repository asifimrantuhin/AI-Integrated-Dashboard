package models

import "time"

type BankLoan struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	YearMonth        time.Time `gorm:"type:date" json:"year_month"`
	CompanyCode      int       `json:"company_code"`
	NetSales         float64   `json:"net_sales"`
	FinancialExpense float64   `json:"financial_expense"`
	ShortCode        string    `json:"short_code"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (BankLoan) TableName() string {
	return "bank_loan"
}

type BankLoanStatusRawData struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	CompanyID int       `json:"company_id"`
	Month     time.Time `gorm:"type:date" json:"month"`
	Head      string    `json:"head"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (BankLoanStatusRawData) TableName() string {
	return "bank_loan_status_raw_data"
}

type BdgtBudget struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Month        time.Time `gorm:"type:date" json:"month"`
	CategoryID   int       `json:"category_id"`
	DepartmentID int       `json:"department_id"`
	Amount       float64   `json:"amount"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (BdgtBudget) TableName() string {
	return "bdgt_budgets"
}

type BdgtExpense struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	BudgetID  int       `json:"budget_id"`
	Month     time.Time `gorm:"type:date" json:"month"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (BdgtExpense) TableName() string {
	return "bdgt_expenses"
}

type BdgtCategory struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (BdgtCategory) TableName() string {
	return "bdgt_categories"
}

type BdgtDepartment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `json:"name"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (BdgtDepartment) TableName() string {
	return "bdgt_departments"
}
