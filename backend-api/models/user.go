package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID                uint           `gorm:"primaryKey" json:"id"`
	Name              string         `gorm:"type:varchar(255);not null" json:"name"`
	Email             string         `gorm:"type:varchar(255);uniqueIndex;not null" json:"email"`
	Phone             string         `gorm:"type:varchar(50)" json:"phone,omitempty"`
	Password          string         `gorm:"type:varchar(255);not null" json:"-"`
	Status            string         `gorm:"type:enum('active','suspended','inactive');default:'active'" json:"status"`
	DefaultCompanyID  *uint          `json:"default_company_id"`
	DefaultCompany    *Company       `gorm:"foreignKey:DefaultCompanyID" json:"default_company,omitempty"`
	Roles             []Role         `gorm:"many2many:user_roles;joinForeignKey:UserID;joinReferences:RoleID" json:"roles,omitempty"`
	Companies         []Company      `gorm:"many2many:user_companies;joinForeignKey:UserID;joinReferences:CompanyID" json:"companies,omitempty"`
	Permissions       []Permission   `gorm:"-" json:"permissions,omitempty"`
	LastLoginAt       *time.Time     `json:"last_login_at,omitempty"`
	LastLoginIP       *string        `json:"last_login_ip,omitempty"`
	CreatedAt         time.Time      `json:"created_at"`
	UpdatedAt         time.Time      `json:"updated_at"`
	DeletedAt         gorm.DeletedAt `gorm:"index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(100);uniqueIndex" json:"name"`
	DisplayName string         `gorm:"type:varchar(150)" json:"display_name"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	IsSystem    bool           `gorm:"default:false" json:"is_system"`
	Permissions []Permission   `gorm:"many2many:role_permissions;joinForeignKey:RoleID;joinReferences:PermissionID" json:"permissions,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Role) TableName() string {
	return "roles"
}

type Permission struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"type:varchar(150);uniqueIndex" json:"name"`
	DisplayName string         `gorm:"type:varchar(150)" json:"display_name"`
	Module      string         `gorm:"type:varchar(100)" json:"module"`
	Description string         `gorm:"type:text" json:"description,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (Permission) TableName() string {
	return "permissions"
}

type UserCompany struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	CompanyID  uint      `gorm:"index" json:"company_id"`
	IsPrimary  bool      `gorm:"default:false" json:"is_primary"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (UserCompany) TableName() string {
	return "user_companies"
}

type UserRole struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	RoleID    uint      `gorm:"index" json:"role_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

type RolePermission struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	RoleID       uint      `gorm:"index" json:"role_id"`
	PermissionID uint      `gorm:"index" json:"permission_id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

type UserToken struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index" json:"user_id"`
	Token      string    `gorm:"type:varchar(80);uniqueIndex" json:"token"`
	Device     string    `gorm:"type:varchar(150)" json:"device,omitempty"`
	IPAddress  string    `gorm:"type:varchar(45)" json:"ip_address,omitempty"`
	UserAgent  string    `gorm:"type:text" json:"user_agent,omitempty"`
	LastUsedAt time.Time `json:"last_used_at"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (UserToken) TableName() string {
	return "user_tokens"
}

