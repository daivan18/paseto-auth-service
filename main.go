package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/daivan18/paseto-auth-service/handler"
	"github.com/daivan18/paseto-auth-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

const privateKeyPath = "keys/secret.key"

func main() {
	// 若不是在 Render 環境中，就載入本地 .env
	if os.Getenv("RENDER") != "true" {
		_ = godotenv.Load()
	}

	// 初始化資料庫
	utils.InitDatabase()

	// 若未設 PASETO_SECRET，檢查是否要產生本地金鑰
	if os.Getenv("PASETO_SECRET") == "" {
		if _, err := os.Stat(privateKeyPath); os.IsNotExist(err) {
			fmt.Println("🔐 金鑰不存在，正在自動產生...")

			// 建立金鑰資料夾（若不存在）
			keyDir := filepath.Dir(privateKeyPath)
			if err := os.MkdirAll(keyDir, 0700); err != nil {
				log.Fatal("❌ 無法建立金鑰資料夾:", err)
			}

			// 產生金鑰
			if err := GenerateAndSaveKey(privateKeyPath); err != nil {
				log.Fatal("❌ 無法生成金鑰:", err)
			}
			fmt.Println("✅ 金鑰產生成功")
		}
	}

	// 載入金鑰（優先從環境變數讀取，否則從檔案）
	if err := handler.LoadKey(); err != nil {
		log.Fatal("❌ 無法讀取金鑰:", err)
	}
	fmt.Println("✅ 金鑰成功載入")

	// 初始化 Gin 路由
	r := gin.Default()

	// API 路由
	r.POST("/api/login", handler.Login)
	r.POST("/api/verify", handler.Verify)

	// ➕ 加入 health check endpoint
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 啟動服務
	port := os.Getenv("PORT") // Render 預設會使用 port 10000
	if port == "" {
		port = "8080" // 預設給本地使用
	}
	log.Println("🚀 Paseto Auth Service is running on port", port)
	r.Run(":" + port) // ⚠️ 修正：使用 `port` 而非重讀 os.Getenv("PORT")
}
