package handler

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/daivan18/paseto-auth-service/utils"
	"github.com/gin-gonic/gin"
	"github.com/o1egl/paseto"
	"golang.org/x/crypto/bcrypt"
)

var pasetoKey []byte

// LoadKey 讀取 private.key 檔案並初始化 pasetoKey
func LoadKey() error {

	// 優先從環境變數 PASETO_SECRET 讀取
	envKey := os.Getenv("PASETO_SECRET")
	if envKey != "" {
		decoded, err := base64.StdEncoding.DecodeString(envKey)
		if err != nil {
			return fmt.Errorf("failed to decode PASETO_SECRET from env: %v", err)
		}
		if len(decoded) != 32 {
			return fmt.Errorf("env key length invalid: expected 32 bytes, got %d", len(decoded))
		}
		pasetoKey = decoded
		fmt.Println("Loaded PASETO key from environment variable.")
		return nil
	}

	// 若環境變數未設，改從檔案讀取
	filePath := "keys/secret.key"

	encodedKey, err := os.ReadFile(filePath) // 讀取檔案中的 Base64 編碼字串
	if err != nil {
		return fmt.Errorf("unable to read key file: %v", err)
	}

	// 將 Base64 編碼的字串解碼為原始的 byte slice
	decodedKey, err := base64.StdEncoding.DecodeString(string(encodedKey))
	if err != nil {
		return fmt.Errorf("unable to decode base64 key: %v", err)
	}

	// 驗證解碼後的金鑰長度是否為 32 bytes
	if len(decodedKey) != 32 {
		return fmt.Errorf("invalid key size after decoding, expected 32 bytes, got %d", len(decodedKey))
	}

	pasetoKey = decodedKey
	return nil
}

// Login 檢查帳密並產生 token
func Login(c *gin.Context) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input",
		})
		return
	}

	var hashedPwd string
	result := utils.GetDB().Table("users").Select("password_hash").Where("username = ?", req.Username).Scan(&hashedPwd)
	if result.Error != nil || result.RowsAffected == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Invalid credentials (user not found)",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(req.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "密碼錯誤",
		})
		return
	}

	now := time.Now()
	exp := now.Add(15 * time.Minute)

	jsonToken := paseto.JSONToken{
		Expiration: exp,
		IssuedAt:   now,
	}
	jsonToken.Set("username", req.Username)

	token, err := paseto.NewV2().Encrypt(pasetoKey, jsonToken, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Token generation failed",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":   token,
		"message": "登入成功",
	})
}

// Verify 給 Python 專案驗證 token 用
func Verify(c *gin.Context) {
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing or invalid Authorization header"})
		return
	}

	token := strings.TrimPrefix(authHeader, "Bearer ")

	var jsonToken paseto.JSONToken
	var footer string
	err := paseto.NewV2().Decrypt(token, pasetoKey, &jsonToken, &footer)
	if err != nil || jsonToken.Expiration.Before(time.Now()) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}

	username := jsonToken.Get("username")
	if username == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or missing username in token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"valid":    true,
		"username": username,
	})
}
