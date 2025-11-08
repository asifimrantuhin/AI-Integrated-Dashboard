package models

import "time"

type NpdProjects struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectName string    `json:"project_name"`
	StartDate   time.Time `gorm:"type:date" json:"start_date"`
	EndDate     time.Time `gorm:"type:date" json:"end_date"`
	Status      string    `json:"status"` // e.g., Planning, Development, Launched
	CompanyID   int       `json:"company_id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (NpdProjects) TableName() string {
	return "npd_projects"
}

type ProjectsDeliverables struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ProjectID   uint      `json:"project_id"`
	Name        string    `json:"name"`
	DueDate     time.Time `gorm:"type:date" json:"due_date"`
	Status      string    `json:"status"` // e.g., Pending, Completed
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (ProjectsDeliverables) TableName() string {
	return "projects_deliverables"
}

type ProjectsSubDeliverables struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	DeliverableID      uint      `json:"deliverable_id"`
	Name               string    `json:"name"`
	DueDate            time.Time `gorm:"type:date" json:"due_date"`
	Status             string    `json:"status"` // e.g., Pending, Completed
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (ProjectsSubDeliverables) TableName() string {
	return "projects_sub_deliverables"
}
