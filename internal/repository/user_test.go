package repository

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/db"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"testing"
)

func generateUser() *model.User {
	// 采用随机生成的用户名和密码，避免重复造成的测试失败
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// 生成随机用户名
	username := ""
	for i := 0; i < 10; i++ {
		username += string(str[rand.Intn(len(str))])
	}

	// 生成随机密码
	password := ""
	for i := 0; i < 10; i++ {
		password += string(str[rand.Intn(len(str))])
	}

	user := &model.User{
		Username: username,
		Password: password,
		Email:    username + "@example.com",
	}

	return user
}

func TestCreateUser(t *testing.T) {
	// Setup
	user := generateUser()

	// Execute
	err := userRepositoryInstance.CreateUser(user)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	db.DB.Unscoped().Delete(user)

	// Test password length > 72 bytes, which should be rejected by bcrypt
	var longPassword string
	for i := 0; i < 73; i++ {
		longPassword += "a"
	}

	// Expect error
	err = userRepositoryInstance.CreateUser(&model.User{Username: user.Username, Password: longPassword})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Cleanup
	db.DB.Unscoped().Delete(user)
}

func TestGetUserByID(t *testing.T) {
	// Setup
	user := generateUser()
	userRepositoryInstance.CreateUser(user)

	// Execute
	result, err := userRepositoryInstance.GetUserByID(user.ID)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Username != user.Username {
		t.Errorf("Expected username %v, got %v", user.Username, result.Username)
	}

	// Cleanup
	db.DB.Unscoped().Delete(user)
}

func TestGetUserByUsername(t *testing.T) {
	// Setup
	user := generateUser()
	userRepositoryInstance.CreateUser(user)

	// Execute
	result, err := userRepositoryInstance.GetUserByUsername(user.Username)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Username != user.Username {
		t.Errorf("Expected username %v, got %v", user.Username, result.Username)
	}

	// Cleanup
	db.DB.Unscoped().Delete(user)
}

func TestGetUserByEmail(t *testing.T) {
	// Setup
	user := generateUser()
	userRepositoryInstance.CreateUser(user)

	// Execute
	result, err := userRepositoryInstance.GetUserByEmail(user.Email)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result.Email != user.Email {
		t.Errorf("Expected email %v, got %v", user.Email, result.Email)
	}

	// Cleanup
	db.DB.Unscoped().Delete(user)
}

func TestVerifyPassword(t *testing.T) {
	// Setup
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	// Execute
	err := userRepositoryInstance.VerifyPassword(string(hashedPassword), password)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCheckUserExists(t *testing.T) {
	// Setup
	user := generateUser()
	userRepositoryInstance.CreateUser(user)

	// Execute
	exists, err := userRepositoryInstance.CheckUserExists(user.Username, user.Email)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if !exists {
		t.Errorf("Expected user to exist")
	}

	// Cleanup
	db.DB.Unscoped().Delete(user)
}
