package utils

import (
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // GORM 資料庫連線變數

// 初始化資料庫連線
func InitDatabase() {
	// Render PostgreSQL 連線字串
	connStr := "postgresql://jim:tOLrquYuXhbSNOmLQ6jFlq7sCKPYiPcS@dpg-cvubllk9c44c73fuovig-a.singapore-postgres.render.com:5432/stock_db_xkna?sslmode=require"

	var err error
	// 使用 GORM 來建立連線
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("無法連接到資料庫:", err)
	}

	// 測試資料庫連線
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("無法獲取資料庫連線:", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("資料庫無法連接:", err)
	}

	log.Println("成功連接到 Render PostgreSQL 資料庫")
}

// GetDB 返回 GORM 資料庫實例
func GetDB() *gorm.DB {
	return DB
}
