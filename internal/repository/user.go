package repository

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/db"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"sync"
)

var (
	userRepositoryInstance UserRepository
	once                   sync.Once

	_ UserRepository = (*UserRepositoryImpl)(nil)
)

// UserRepository 定义用户仓库接口
type UserRepository interface {
	CreateUser(user *model.User) error
	GetUserByID(id uint) (*model.User, error)
	GetUserByUsername(username string) (*model.User, error)
	GetUserByEmail(email string) (*model.User, error)
	VerifyPassword(hashedPassword, password string) error
	CheckUserExists(username, email string) (bool, error)
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	once.Do(func() {
		userRepositoryInstance = &UserRepositoryImpl{
			DB: db,
		}
	})
	return userRepositoryInstance
}

// CreateUser 创建一个新的用户
func (r *UserRepositoryImpl) CreateUser(user *model.User) error {
	// 对密码进行哈希处理
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	return r.DB.Create(user).Error
}

// GetUserByID 根据用户 ID 获取用户信息
func (r *UserRepositoryImpl) GetUserByID(id uint) (*model.User, error) {
	var user model.User
	result := r.DB.Preload("SectionStatus").First(&user, id)
	return &user, result.Error
}

// GetUserByUsername 根据用户名获取用户信息
func (r *UserRepositoryImpl) GetUserByUsername(username string) (*model.User, error) {
	var user model.User
	result := r.DB.Where("username = ?", username).First(&user)
	return &user, result.Error
}

// GetUserByEmail 根据邮箱获取用户信息
func (r *UserRepositoryImpl) GetUserByEmail(email string) (*model.User, error) {
	var user model.User
	result := r.DB.Where("email = ?", email).First(&user)
	return &user, result.Error
}

// VerifyPassword 验证用户密码
func (r *UserRepositoryImpl) VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

// CheckUserExists 检查用户名或邮箱是否已存在
func (r *UserRepositoryImpl) CheckUserExists(username, email string) (bool, error) {
	var count int64
	result := r.DB.Model(&model.User{}).Where("username = ? OR email = ?", username, email).Count(&count)
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
