package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"idash-backend-api/database"
	"idash-backend-api/models"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserManagementRequest struct {
	Name             string  `json:"name"`
	Email            string  `json:"email"`
	Phone            string  `json:"phone"`
	Status           string  `json:"status"`
	Password         *string `json:"password"`
	Roles            []string `json:"roles"`
	CompanyIDs       []uint   `json:"company_ids"`
	DefaultCompanyID *uint    `json:"default_company_id"`
}

func ListUsers(c echo.Context) error {
	var users []models.User
	if err := database.DB.Preload("Roles").Preload("Companies").Preload("DefaultCompany").Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load users"})
	}

	result := make([]AuthResponse, 0, len(users))
	for _, user := range users {
		roles, permissions := aggregateUserAccess(&user)
		companyIDs := extractCompanyIDs(user.Companies)
		user.Password = ""
		result = append(result, AuthResponse{
			Token:            "",
			User:             user,
			Roles:            roles,
			Permissions:      permissions,
			DefaultCompanyID: user.DefaultCompanyID,
			CompanyIDs:       companyIDs,
		})
	}

	return c.JSON(http.StatusOK, result)
}

func GetUserDetail(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var user models.User
	if err := preloadUserByID(&user, id); err != nil {
		if err == gorm.ErrRecordNotFound {
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

func CreateUser(c echo.Context) error {
	var req UserManagementRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	if strings.TrimSpace(req.Email) == "" || req.Password == nil || strings.TrimSpace(*req.Password) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Email and password are required"})
	}

	registerReq := RegisterRequest{
		Name:             req.Name,
		Email:            strings.TrimSpace(req.Email),
		Phone:            req.Phone,
		Password:         strings.TrimSpace(*req.Password),
		Roles:            req.Roles,
		CompanyIDs:       req.CompanyIDs,
		DefaultCompanyID: req.DefaultCompanyID,
		Status:           req.Status,
	}

	user, rolesList, permissions, companyIDs, err := createUserAccount(registerReq)
	if err != nil {
		if httpErr, ok := err.(*echo.HTTPError); ok {
			return c.JSON(httpErr.Code, map[string]string{"error": httpErr.Message.(string)})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to create user"})
	}

	return c.JSON(http.StatusCreated, AuthResponse{
		Token:            "",
		User:             user,
		Roles:            rolesList,
		Permissions:      permissions,
		DefaultCompanyID: user.DefaultCompanyID,
		CompanyIDs:       companyIDs,
	})
}

func UpdateUser(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	var req UserManagementRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request"})
	}

	tx := database.DB.Begin()
	if tx.Error != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to initialize transaction"})
	}

	var user models.User
	if err := tx.Where("id = ?", id).First(&user).Error; err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return c.JSON(http.StatusNotFound, map[string]string{"error": "User not found"})
		}
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load user"})
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Phone != "" {
		updates["phone"] = req.Phone
	}
	if req.Status != "" {
		updates["status"] = strings.ToLower(req.Status)
	}
	if req.DefaultCompanyID != nil {
		updates["default_company_id"] = req.DefaultCompanyID
	}

	if len(updates) > 0 {
		if err := tx.Model(&user).Updates(updates).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
		}
	}

	if req.Password != nil && strings.TrimSpace(*req.Password) != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to hash password"})
		}
		if err := tx.Model(&user).Update("password", string(hashed)).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update password"})
		}
	}

	if len(req.Roles) > 0 {
		var roles []models.Role
		if err := tx.Where("name IN ?", req.Roles).Find(&roles).Error; err != nil {
			tx.Rollback()
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid roles specified"})
		}
		if err := tx.Model(&user).Association("Roles").Replace(roles); err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to assign roles"})
		}
	}

	if req.CompanyIDs != nil {
		var companies []models.Company
		if len(req.CompanyIDs) > 0 {
			if err := tx.Where("id IN ?", req.CompanyIDs).Find(&companies).Error; err != nil {
				tx.Rollback()
				return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid companies specified"})
			}
		}
		if err := tx.Model(&user).Association("Companies").Replace(companies); err != nil {
			tx.Rollback()
			return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to assign companies"})
		}

		if req.DefaultCompanyID != nil && !containsUint(req.CompanyIDs, *req.DefaultCompanyID) {
			tx.Rollback()
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "Default company must be within assigned companies"})
		}
	}

	if err := tx.Commit().Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to update user"})
	}

	if err := preloadUserByID(&user, user.ID); err != nil {
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

func DeleteUser(c echo.Context) error {
	id, err := parseUintParam(c, "id")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid user ID"})
	}

	if err := database.DB.Delete(&models.User{}, id).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to delete user"})
	}

	return c.NoContent(http.StatusNoContent)
}

func ListRoles(c echo.Context) error {
	var roles []models.Role
	if err := database.DB.Preload("Permissions").Find(&roles).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load roles"})
	}
	return c.JSON(http.StatusOK, roles)
}

func ListPermissions(c echo.Context) error {
	var permissions []models.Permission
	if err := database.DB.Find(&permissions).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Failed to load permissions"})
	}
	return c.JSON(http.StatusOK, permissions)
}

func parseUintParam(c echo.Context, name string) (uint, error) {
	value := c.Param(name)
	id64, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, err
	}
	return uint(id64), nil
}
