package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"time"
)

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

// GenerateAccessToken 生成 JWT Token
func GenerateAccessToken(userID int) (string, error) {

	// 获取密钥
	secretKey := viper.Get("jwt.secret").([]byte)
	if len(secretKey) == 0 {
		logrus.Warn("secret key is empty, setting default key")
		secretKey = []byte("ttds")
	}

	// 获取过期时间
	expirationTimeStr := viper.GetString("jwt.expires")
	// 获取最后一个字符来判断时间单位
	if len(expirationTimeStr) == 0 {
		expirationTimeStr = "10m"
	}

	expirationTime, err := time.ParseDuration(expirationTimeStr)
	if err != nil {
		logrus.Warn("expiration time is invalid, setting default value")
		expirationTime = 10 * time.Minute
	}

	// 创建自定义的 Claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建新的 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名并生成 token 字符串
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken 生成 JWT Token
func GenerateRefreshToken(userID int) (string, error) {

	// 获取密钥
	secretKey := viper.Get("jwt.refresh_secret").([]byte)
	if len(secretKey) == 0 {
		logrus.Warn("secret key is empty, setting default key")
		secretKey = []byte("ttds")
	}

	// 获取过期时间
	expirationTimeStr := viper.GetString("jwt.refresh_expires")
	// 获取最后一个字符来判断时间单位
	if len(expirationTimeStr) == 0 {
		expirationTimeStr = "24h"
	}

	expirationTime, err := time.ParseDuration(expirationTimeStr)
	if err != nil {
		logrus.Warn("expiration time is invalid, setting default value")
		expirationTime = 24 * time.Hour
	}

	// 创建自定义的 Claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expirationTime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建新的 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名并生成 token 字符串
	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证 JWT Token
func ValidateToken(tokenStr string) (*Claims, error) {

	// 获取密钥
	secretKey := viper.Get("jwt.secret").([]byte)
	if len(secretKey) == 0 {
		logrus.Warn("secret key is empty, setting default key")
		secretKey = []byte("ttds")
	}
	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 确保 Token 使用的是正确的签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// 验证 Token 是否有效
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid token")
}

// ValidateRefreshToken 解析 Refresh Token
func ValidateRefreshToken(refreshTokenStr string) (*Claims, error) {
	// 获取密钥
	refreshSecretKey := viper.Get("jwt.refresh_secret").([]byte)
	// 解析 Token
	token, err := jwt.ParseWithClaims(refreshTokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return refreshSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}
	return nil, fmt.Errorf("invalid refresh token")
}

// RefreshAccessToken 刷新 Access Token
func RefreshAccessToken(refreshTokenStr string) (string, error) {
	// 验证 Refresh Token
	claims, err := ValidateRefreshToken(refreshTokenStr)
	if err != nil {
		return "", err
	}

	// 生成新的 Access Token
	return GenerateAccessToken(claims.UserID)
}
