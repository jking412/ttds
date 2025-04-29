package jwt

import (
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/assert"
)

func TestGenerateAccessToken(t *testing.T) {
	manager := NewJWTManager()

	tests := []struct {
		name   string
		userID int
		err    error
	}{
		{"正常生成令牌", 1, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := manager.GenerateAccessToken(tt.userID)
			assert.Equal(t, tt.err, err)
			if err == nil {
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestGenerateRefreshToken(t *testing.T) {
	manager := NewJWTManager()

	tests := []struct {
		name   string
		userID int
		err    error
	}{
		{"正常生成刷新令牌", 1, nil},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := manager.GenerateRefreshToken(tt.userID)
			assert.Equal(t, tt.err, err)
			if err == nil {
				assert.NotEmpty(t, token)
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	manager := NewJWTManager()
	userID := 1

	token, err := manager.GenerateAccessToken(userID)
	assert.NoError(t, err)

	tests := []struct {
		name  string
		token string
		err   error
	}{
		{"验证有效令牌", token, nil},
		{"验证无效令牌", "invalid.token", jwt.ErrSignatureInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := manager.ValidateToken(tt.token)
			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, userID, claims.UserID)
			}
		})
	}
}

func TestRefreshAccessToken(t *testing.T) {
	manager := NewJWTManager()
	userID := 1

	refreshToken, err := manager.GenerateRefreshToken(userID)
	assert.NoError(t, err)

	tests := []struct {
		name  string
		token string
		err   error
	}{
		{"刷新有效令牌", refreshToken, nil},
		{"刷新无效令牌", "invalid.token", jwt.ErrSignatureInvalid},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := manager.RefreshAccessToken(tt.token)
			if tt.err != nil {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, token)
			}
		})
	}
}
