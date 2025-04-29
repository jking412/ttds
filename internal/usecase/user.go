package usecase

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"awesomeProject/pkg/db"
	"awesomeProject/pkg/jwt"
	"errors"
	"sync"
)

var (
	userServiceInstance UserService
	once                sync.Once

	_ UserService = (*UserServiceImpl)(nil)
)

// UserService 用户服务接口
type UserService interface {
	Register(user *model.User) (string, string, error)
	Login(username, password string) (string, string, error)
	Logout() error
	GetCurrentUser(userID uint) (*model.User, error)
}

// UserServiceImpl 用户服务实现
type UserServiceImpl struct {
	UserRepository repository.UserRepository
	jwtManager     jwt.JWTManager
}

func NewUserService() UserService {
	once.Do(func() {
		userServiceInstance = &UserServiceImpl{
			UserRepository: repository.NewUserRepository(db.DB),
			jwtManager:     jwt.NewJWTManager(),
		}
	})
	return userServiceInstance
}

// Register 用户注册
func (s *UserServiceImpl) Register(user *model.User) (string, string, error) {
	// 检查用户名和邮箱是否已存在
	exists, err := s.UserRepository.CheckUserExists(user.Username, user.Email)
	if err != nil {
		return "", "", err
	}

	if exists {
		return "", "", errors.New("用户名或邮箱已存在")
	}

	// 创建新用户
	if err := s.UserRepository.CreateUser(user); err != nil {
		return "", "", err
	}

	// 生成JWT令牌
	userID := int(user.ID)
	accessToken, err := s.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Login 用户登录
func (s *UserServiceImpl) Login(username, password string) (string, string, error) {
	// 根据用户名查找用户
	user, err := s.UserRepository.GetUserByUsername(username)
	if err != nil {
		return "", "", errors.New("用户名或密码错误")
	}

	// 验证密码
	if err := s.UserRepository.VerifyPassword(user.Password, password); err != nil {
		return "", "", errors.New("用户名或密码错误")
	}

	// 生成JWT令牌
	userID := int(user.ID)
	accessToken, err := s.jwtManager.GenerateAccessToken(userID)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := s.jwtManager.GenerateRefreshToken(userID)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// Logout 用户注销
func (s *UserServiceImpl) Logout() error {
	// 这里只返回成功消息，实际的注销逻辑应该在客户端完成
	return nil
}

func (s *UserServiceImpl) GetCurrentUser(userID uint) (*model.User, error) {
	user, err := s.UserRepository.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	user.Password = ""
	return user, nil
}
