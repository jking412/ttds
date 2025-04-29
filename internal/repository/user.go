package repository

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

//func AnswerCheck(sectionID uint) error {
//
//}

// CreateUser 创建一个新的用户
func CreateUser(user *model.User) error {
	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return db.DB.Create(user).Error
}

// GetUserByID 根据用户 ID 获取用户信息
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	result := db.DB.Preload("SectionStatus").First(&user, id)
	return &user, result.Error
}

// GetUserByUsername 根据用户名获取用户信息
func GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	result := db.DB.Where("username = ?", username).First(&user)
	return &user, result.Error
}

// GetUserByEmail 根据邮箱获取用户信息
func GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := db.DB.Where("email = ?", email).First(&user)
	return &user, result.Error
}

// VerifyPassword 验证用户密码
func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// CheckUserExists 检查用户名或邮箱是否已存在
func CheckUserExists(username, email string) (bool, error) {
	var count int64
	result := db.DB.Model(&model.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
	if result.Error != nil {
		return false, result.Error
	}
	return count > 0, nil
}

// UpdateUserSectionStatus 更新用户小节完成状态
func UpdateUserSectionStatus(status *model.UserSectionStatus) error {
	// 先根据userId和sectionId查询到主键，然后更新
	s := model.UserSectionStatus{UserID: status.UserID, SectionID: status.SectionID}
	if err := db.DB.Where("user_id = ? AND section_id = ?", status.UserID, status.SectionID).First(&s).Error; err != nil {
		return err
	}
	status.ID = s.ID
	status.CreatedAt = s.CreatedAt
	return db.DB.Save(status).Error
}
