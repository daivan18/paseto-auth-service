package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"os"
)

func GenerateAndSaveKey(filePath string) error {
	key := make([]byte, 32)
	if _, err := rand.Read(key); err != nil {
		return err
	}

	encoded := base64.StdEncoding.EncodeToString(key)

	// 寫入 key 到檔案
	err := os.WriteFile(filePath, []byte(encoded), 0600)
	if err != nil {
		return err
	}

	fmt.Printf("✅ 金鑰已寫入 %s\n", filePath)
	fmt.Println("🔐 若要使用環境變數部署，請設定以下內容到你的環境：")
	fmt.Printf("export PASETO_SECRET=\"%s\"\n", encoded)

	return nil
}
