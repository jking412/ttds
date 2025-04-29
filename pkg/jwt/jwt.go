package jwt

import (
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
	"time"
)

var (
	jwtManagerInstance JWTManager
	once               sync.Once

	_ JWTManager = (*jwtManager)(nil)
)

const (
	// 定义 token 过期时间
	defaultExpirationTime = 24 * time.Hour
	// 定义 refresh token 过期时间
	defaultRefreshExpirationTime = 7 * 24 * time.Hour
	// 定义 token 签名密钥
	defaultSecretKey = "ttds"
	// 定义 refresh token 签名密钥
	defaultRefreshSecretKey = "ttds"
)

type JWTManager interface {
	GenerateAccessToken(userID int) (string, error)
	GenerateRefreshToken(userID int) (string, error)
	ValidateToken(tokenStr string) (*Claims, error)
	RefreshAccessToken(refreshToken string) (string, error)
}

type jwtManager struct {
	secretKey     []byte
	refreshKey    []byte
	secretExpiry  time.Duration
	refreshExpiry time.Duration
}

type Claims struct {
	UserID int `json:"user_id"`
	jwt.RegisteredClaims
}

func NewJWTManager() JWTManager {
	once.Do(func() {
		secretKey := viper.GetString("jwt_secret_key")
		refreshKey := viper.GetString("jwt.refresh_secret")

		if len(secretKey) == 0 {
			logrus.Warn("secret key is empty, setting default key")
			secretKey = defaultSecretKey
		}

		if len(refreshKey) == 0 {
			logrus.Warn("refresh key is empty, setting default key")
			refreshKey = defaultRefreshSecretKey
		}

		secretExpiryStr := viper.GetString("jwt.expires")
		refreshExpiryStr := viper.GetString("jwt.refresh_expires")

		var secretExpiry, refreshExpiry time.Duration
		var err error

		if len(secretExpiryStr) == 0 {
			logrus.Warn("secret expiry is invalid, setting default value")
			secretExpiry = defaultExpirationTime
		}

		if len(refreshExpiryStr) == 0 {
			logrus.Warn("refresh expiry is invalid, setting default value")
			refreshExpiry = defaultRefreshExpirationTime
		}

		secretExpiry, err = time.ParseDuration(secretExpiryStr)
		if err != nil {
			logrus.Warn("secret expiry is invalid, setting default value")
			secretExpiry = defaultExpirationTime
		}

		refreshExpiry, err = time.ParseDuration(refreshExpiryStr)
		if err != nil {
			logrus.Warn("refresh expiry is invalid, setting default value")
			refreshExpiry = defaultRefreshExpirationTime
		}

		jwtManagerInstance = &jwtManager{
			secretKey:     []byte(secretKey),
			refreshKey:    []byte(refreshKey),
			secretExpiry:  secretExpiry,
			refreshExpiry: refreshExpiry,
		}
	})

	return jwtManagerInstance
}

// GenerateAccessToken 生成 JWT Token
func (j *jwtManager) GenerateAccessToken(userID int) (string, error) {

	// 创建自定义的 Claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.secretExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建新的 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名并生成 token 字符串
	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken 生成 JWT Token
func (j *jwtManager) GenerateRefreshToken(userID int) (string, error) {

	// 创建自定义的 Claims
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(j.refreshExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	// 创建新的 Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 使用密钥签名并生成 token 字符串
	tokenString, err := token.SignedString(j.refreshKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken 验证 JWT Token
func (j *jwtManager) ValidateToken(tokenStr string) (*Claims, error) {

	// 解析 Token
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 确保 Token 使用的是正确的签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.secretKey, nil
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
func (j *jwtManager) validateRefreshToken(refreshTokenStr string) (*Claims, error) {

	// func
	token, err := jwt.ParseWithClaims(refreshTokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {

		// 解析Token所用的加密算法，并根据算法返回对应的密钥
		// 在这里我们只使用一种算法，如果不是我们期望的算法，就返回错误
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return j.refreshKey, nil
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
func (j *jwtManager) RefreshAccessToken(refreshTokenStr string) (string, error) {
	// 验证 Refresh Token
	claims, err := j.validateRefreshToken(refreshTokenStr)
	if err != nil {
		return "", err
	}

	// 生成新的 Access Token
	return j.GenerateAccessToken(claims.UserID)
}
