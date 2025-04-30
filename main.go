package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/daivan18/paseto-auth-service/handler"
	"github.com/daivan18/paseto-auth-service/utils"
	"github.com/gin-gonic/gin"
)

// 設定金鑰路徑
const privateKeyPath = "keys/secret.key"

func main() {
	// 初始化資料庫
	utils.InitDatabase()

	// 若環境變數未設，才檢查並產生檔案金鑰
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

	// 提供給其他專案的 API
	r.POST("/api/login", handler.Login)   // 提供 Token 產生
	r.POST("/api/verify", handler.Verify) // 提供 Token 驗證

	// 啟動服務
	r.Run(":8080")
}
