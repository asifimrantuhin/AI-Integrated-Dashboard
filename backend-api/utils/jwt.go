package utils

import (
	"idash-backend-api/config"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	UserID           uint     `json:"user_id"`
	Roles            []string `json:"roles"`
	Permissions      []string `json:"permissions"`
	DefaultCompanyID *uint    `json:"default_company_id,omitempty"`
	CompanyIDs       []uint   `json:"company_ids,omitempty"`
	jwt.RegisteredClaims
}

func GenerateJWT(userID uint, roles []string, permissions []string, defaultCompanyID *uint, companyIDs []uint) (string, error) {
	expires := time.Now().Add(time.Hour * time.Duration(config.AppConfig.JWTTTL))
	claims := Claims{
		UserID:           userID,
		Roles:            roles,
		Permissions:      permissions,
		DefaultCompanyID: defaultCompanyID,
		CompanyIDs:       companyIDs,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expires),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.AppConfig.JWTSecret))
}

func ValidateJWT(tokenString string) (*Claims, error) {
	claims := &Claims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(config.AppConfig.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}

	return claims, nil
}

