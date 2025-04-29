package repository

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/db"
	"fmt"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestMain(m *testing.M) {
	// 初始化数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"123456",
		"localhost",
		3306,
		"ttds",
	)
	db.InitDB(dsn)

	// 执行测试
	m.Run()
}

func TestCreateUser(t *testing.T) {
	// Setup
	user := &model.User{Username: "testuser", Password: "password123"}

	// Execute
	err := CreateUser(user)

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
	err = CreateUser(&model.User{Username: "testuser", Password: longPassword})
	if err == nil {
		t.Errorf("Expected error, got nil")
	}

	// Cleanup
	db.DB.Unscoped().Delete(user)
}

func TestGetUserByID(t *testing.T) {
	// Setup
	user := &model.User{Username: "testuser", Password: "password123"}
	CreateUser(user)

	// Execute
	result, err := GetUserByID(user.ID)

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
	user := &model.User{Username: "testuser", Password: "password123"}
	CreateUser(user)

	// Execute
	result, err := GetUserByUsername(user.Username)

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
	user := &model.User{Email: "test@example.com", Password: "password123"}
	CreateUser(user)

	// Execute
	result, err := GetUserByEmail(user.Email)

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
	err := VerifyPassword(string(hashedPassword), password)

	// Verify
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestCheckUserExists(t *testing.T) {
	// Setup
	user := &model.User{Username: "testuser", Email: "test@example.com", Password: "password123"}
	CreateUser(user)

	// Execute
	exists, err := CheckUserExists(user.Username, user.Email)

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
