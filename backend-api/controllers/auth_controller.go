package controllers

import (
	"errors"
	"net/http"
	"strings"
	"time"

	"idash-backend-api/database"
	"idash-backend-api/models"
	"idash-backend-api/utils"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type LoginRequest struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required"`
	CompanyID *uint  `json:"company_id"`
}

type RegisterRequest struct {
	Name             string   `json:"name" validate:"required"`
	Email            string   `json:"email" validate:"required,email"`
	Phone            string   `json:"phone"`
	Password         string   `json:"password" validate:"required,min=6"`
	Roles            []string `json:"roles"`
	CompanyIDs       []uint   `json:"company_ids"`
	DefaultCompanyID *uint    `json:"default_company_id"`
	Status           string   `json:"status"`
}

type AuthResponse struct {
	Token            string      `json:"token"`
	User             models.User `json:"user"`
	Roles            []string    `json:"roles"`
	Permissions      []string    `json:"permissions"`
	DefaultCompanyID *uint       `json:"default_company_id"`
	CompanyIDs       []uint      `json:"company_ids"`
}

func Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	var user models.User
	if err := preloadUserByEmail(&user, req.Email); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
		}
		// Log the underlying DB error for diagnostics and return a clearer message
		c.Logger().Errorf("preloadUserByEmail error: %v", err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load user: " + err.Error()})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid credentials"})
	}

	roles, permissions := aggregateUserAccess(&user)
	companyIDs := extractCompanyIDs(user.Companies)

	// Validate requested company context
	defaultCompanyID := user.DefaultCompanyID
	if req.CompanyID != nil {
		if !containsUint(companyIDs, *req.CompanyID) {
			return c.JSON(http.StatusForbidden, map[string]string{"error": "Access to requested company denied"})
		}
		defaultCompanyID = req.CompanyID
	}

	token, err := utils.GenerateJWT(user.ID, roles, permissions, defaultCompanyID, companyIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	now := time.Now()
	ip := c.RealIP()
	database.DB.Model(&user).Updates(map[string]interface{}{
		"last_login_at": now,
		"last_login_ip": ip,
	})

	user.Password = ""
	user.Permissions = make([]models.Permission, 0)

	return c.JSON(http.StatusOK, AuthResponse{
		Token:            token,
		User:             user,
		Roles:            roles,
		Permissions:      permissions,
		DefaultCompanyID: defaultCompanyID,
		CompanyIDs:       companyIDs,
	})
}

func Register(c echo.Context) error {
	var req RegisterRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	user, rolesList, permissions, companyIDs, err := createUserAccount(req)
	if err != nil {
		if httpErr, ok := err.(*echo.HTTPError); ok {
			return c.JSON(httpErr.Code, map[string]string{"error": httpErr.Message.(string)})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	token, err := utils.GenerateJWT(user.ID, rolesList, permissions, user.DefaultCompanyID, companyIDs)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to generate token"})
	}

	return c.JSON(http.StatusCreated, AuthResponse{
		Token:            token,
		User:             user,
		Roles:            rolesList,
		Permissions:      permissions,
		DefaultCompanyID: user.DefaultCompanyID,
		CompanyIDs:       companyIDs,
	})
}

func GetUser(c echo.Context) error {
	userIDAny := c.Get("user_id")
	userID, ok := userIDAny.(uint)
	if !ok {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "Invalid user context"})
	}

	var user models.User
	if err := preloadUserByID(&user, userID); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load user"})
	}

	roles, permissions := aggregateUserAccess(&user)
	companyIDs := extractCompanyIDs(user.Companies)

	user.Password = ""

	return c.JSON(http.StatusOK, AuthResponse{
		Token:            "",
		User:             user,
		Roles:            roles,
		Permissions:      permissions,
		DefaultCompanyID: user.DefaultCompanyID,
		CompanyIDs:       companyIDs,
	})
}

func preloadUserByEmail(user *models.User, email string) error {
	// Use Unscoped() to avoid GORM adding soft-delete filters (deleted_at IS NULL)
	// because some legacy schemas may not include deleted_at columns on join tables.
	return database.DB.Unscoped().
		Preload("Roles.Permissions").
		Preload("Companies").
		Preload("DefaultCompany").
		Where("email = ?", email).
		First(user).Error
}

func preloadUserByID(user *models.User, id uint) error {
	// See note in preloadUserByEmail about Unscoped usage
	return database.DB.Unscoped().
		Preload("Roles.Permissions").
		Preload("Companies").
		Preload("DefaultCompany").
		Where("id = ?", id).
		First(user).Error
}

func aggregateUserAccess(user *models.User) ([]string, []string) {
	roleSet := make(map[string]struct{})
	permissionSet := make(map[string]struct{})

	for _, role := range user.Roles {
		roleSet[strings.ToLower(role.Name)] = struct{}{}
		for _, permission := range role.Permissions {
			permissionSet[strings.ToLower(permission.Name)] = struct{}{}
		}
	}

	roles := make([]string, 0, len(roleSet))
	for role := range roleSet {
		roles = append(roles, role)
	}

	permissions := make([]string, 0, len(permissionSet))
	for permission := range permissionSet {
		permissions = append(permissions, permission)
	}

	return roles, permissions
}

func extractCompanyIDs(companies []models.Company) []uint {
	ids := make([]uint, 0, len(companies))
	for _, company := range companies {
		ids = append(ids, company.ID)
	}
	return ids
}

func containsUint(list []uint, value uint) bool {
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func createUserAccount(req RegisterRequest) (models.User, []string, []string, []uint, error) {
	rolesRequested := req.Roles
	if len(rolesRequested) == 0 {
		rolesRequested = []string{"analyst"}
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return models.User{}, nil, nil, nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to initialize transaction")
	}

	var existing models.User
	if err := tx.Where("email = ?", req.Email).First(&existing).Error; err == nil {
		tx.Rollback()
		return models.User{}, nil, nil, nil, echo.NewHTTPError(http.StatusConflict, "User already exists")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		return models.User{}, nil, nil, nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to hash password")
	}

	status := strings.ToLower(strings.TrimSpace(req.Status))
	if status == "" {
		status = "active"
	}

	user := models.User{
		Name:             req.Name,
		Email:            req.Email,
		Phone:            req.Phone,
		Password:         string(hashedPassword),
		Status:           status,
		DefaultCompanyID: req.DefaultCompanyID,
	}

	if err := tx.Create(&user).Error; err != nil {
		tx.Rollback()
		return models.User{}, nil, nil, nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to create user")
	}

	var roles []models.Role
	if err := tx.Where("name IN ?", rolesRequested).Find(&roles).Error; err != nil || len(roles) == 0 {
		tx.Rollback()
		return models.User{}, nil, nil, nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid roles specified")
	}

	if err := tx.Model(&user).Association("Roles").Replace(roles); err != nil {
		tx.Rollback()
		return models.User{}, nil, nil, nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to assign roles")
	}

	if len(req.CompanyIDs) > 0 {
		var companies []models.Company
		if err := tx.Where("id IN ?", req.CompanyIDs).Find(&companies).Error; err != nil {
			tx.Rollback()
			return models.User{}, nil, nil, nil, echo.NewHTTPError(http.StatusBadRequest, "Invalid companies specified")
		}
		if err := tx.Model(&user).Association("Companies").Replace(companies); err != nil {
			tx.Rollback()
			return models.User{}, nil, nil, nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to assign companies")
		}

		if req.DefaultCompanyID != nil && !containsUint(req.CompanyIDs, *req.DefaultCompanyID) {
			tx.Rollback()
			return models.User{}, nil, nil, nil, echo.NewHTTPError(http.StatusBadRequest, "Default company must be within assigned companies")
		}
	}

	if err := tx.Commit().Error; err != nil {
		return models.User{}, nil, nil, nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to save user")
	}

	if err := preloadUserByID(&user, user.ID); err != nil {
		return models.User{}, nil, nil, nil, echo.NewHTTPError(http.StatusInternalServerError, "Failed to load user")
	}

	rolesList, permissions := aggregateUserAccess(&user)
	companyIDs := extractCompanyIDs(user.Companies)
	user.Password = ""

	return user, rolesList, permissions, companyIDs, nil
}
