package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"
	
	"gorm.io/gorm"
)

// APICache represents cache table structure
type APICache struct {
	CacheKey   string    `gorm:"primaryKey"`
	CacheValue string    `gorm:"type:text"`
	Endpoint   string    `gorm:"index"`
	Parameters string    `gorm:"type:text"`
	ExpiresAt  time.Time `gorm:"index"`
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

func (APICache) TableName() string {
	return "api_cache"
}

// GenerateCacheKey generates a unique cache key from endpoint and parameters
func GenerateCacheKey(endpoint string, params map[string]interface{}) string {
	paramBytes, _ := json.Marshal(params)
	hash := md5.Sum([]byte(endpoint + string(paramBytes)))
	return fmt.Sprintf("%s:%s", endpoint, hex.EncodeToString(hash[:]))
}

// GetFromCache retrieves cached data
func GetFromCache(db *gorm.DB, cacheKey string) (string, bool) {
	var cache APICache
	result := db.Where("cache_key = ? AND expires_at > ?", cacheKey, time.Now()).First(&cache)
	if result.Error != nil {
		return "", false
	}
	return cache.CacheValue, true
}

// SetCache stores data in cache
func SetCache(db *gorm.DB, cacheKey, endpoint string, params map[string]interface{}, value string, ttl time.Duration) error {
	paramStr := ""
	if params != nil {
		paramBytes, _ := json.Marshal(params)
		paramStr = string(paramBytes)
	}
	
	cache := APICache{
		CacheKey:   cacheKey,
		CacheValue: value,
		Endpoint:   endpoint,
		Parameters: paramStr,
		ExpiresAt:  time.Now().Add(ttl),
	}
	
	return db.Save(&cache).Error
}

// ClearCache removes expired cache entries
func ClearCache(db *gorm.DB) error {
	return db.Where("expires_at < ?", time.Now()).Delete(&APICache{}).Error
}

// ClearCacheByEndpoint removes cache for a specific endpoint
func ClearCacheByEndpoint(db *gorm.DB, endpoint string) error {
	return db.Where("endpoint = ?", endpoint).Delete(&APICache{}).Error
}

