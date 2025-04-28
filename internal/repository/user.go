package repository

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/db"
)

//func AnswerCheck(sectionID uint) error {
//
//}

// CreateUser 创建一个新的用户
func CreateUser(user *model.User) error {
	return db.DB.Create(user).Error
}

// GetUserByID 根据用户 ID 获取用户信息
func GetUserByID(id uint) (*model.User, error) {
	var user model.User
	result := db.DB.Preload("SectionStatus").First(&user, id)
	return &user, result.Error
}

// CreateUserSectionStatus 创建用户小节完成状态记录
func CreateUserSectionStatus(status *model.UserSectionStatus) error {
	return db.DB.Create(status).Error
}

// GetUserSectionStatusByUserAndSectionID 根据用户 ID 和小节 ID 获取用户小节完成状态
func GetUserSectionStatusByUserAndSectionID(userID, sectionID uint) (*model.UserSectionStatus, error) {
	var status model.UserSectionStatus
	result := db.DB.Where("user_id = ? AND section_id = ?", userID, sectionID).First(&status)
	return &status, result.Error
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
