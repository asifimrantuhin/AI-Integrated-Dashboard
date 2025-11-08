package database

import (
	"idash-backend-api/config"
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	dsn := config.AppConfig.DBUser + ":" + config.AppConfig.DBPassword + "@tcp(" + config.AppConfig.DBHost + ":" + config.AppConfig.DBPort + ")/" + config.AppConfig.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	log.Println("Database connected successfully")
}

func GetDB() *gorm.DB {
	return DB
}

