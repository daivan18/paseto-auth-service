package utils

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB // GORM 資料庫連線變數

// 初始化資料庫連線
func InitDatabase() {
	var connStr string

	// 判斷是否在 Render 環境
	if os.Getenv("RENDER") == "true" {
		// Render 雲端部署：用 Internal DB URL（Render Web UI 設定）
		connStr = os.Getenv("DATABASE_INTERNAL_URL")
	} else {
		// 本地開發：用 External DB URL（.env 設定）
		connStr = os.Getenv("DATABASE_URL")
	}

	if connStr == "" {
		log.Fatal("❌ 資料庫連線字串未設定，請確認環境變數 DATABASE_URL 或 DATABASE_INTERNAL_URL")
	}

	var err error
	DB, err = gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ 無法連接到資料庫:", err)
	}

	// 測試資料庫連線
	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatal("❌ 無法獲取資料庫連線:", err)
	}
	if err := sqlDB.Ping(); err != nil {
		log.Fatal("❌ 資料庫無法連接:", err)
	}

	log.Println("✅ 成功連接到 PostgreSQL 資料庫")
}

// GetDB 返回 GORM 資料庫實例
func GetDB() *gorm.DB {
	return DB
}
