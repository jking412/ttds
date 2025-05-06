package usecase

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/jwt"
	"errors"
	"gorm.io/gorm"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) CreateUser(user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetUserByID(id uint) (*model.User, error) {
	args := m.Called(id)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByUsername(username string) (*model.User, error) {
	args := m.Called(username)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*model.User, error) {
	args := m.Called(email)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) VerifyPassword(hashedPassword, password string) error {
	args := m.Called(hashedPassword, password)
	return args.Error(0)
}

func (m *MockUserRepository) CheckUserExists(username, email string) (bool, error) {
	args := m.Called(username, email)
	return args.Bool(0), args.Error(1)
}

type MockJWTGenerator struct {
	mock.Mock
}

//type Manager interface {
//	GenerateAccessToken(userID int) (string, error)
//	GenerateRefreshToken(userID int) (string, error)
//	ValidateToken(tokenStr string) (*Claims, error)
//	RefreshAccessToken(refreshToken string) (string, error)
//}

func (m *MockJWTGenerator) GenerateAccessToken(userID uint) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTGenerator) GenerateRefreshToken(userID uint) (string, error) {
	args := m.Called(userID)
	return args.String(0), args.Error(1)
}

func (m *MockJWTGenerator) ValidateToken(tokenStr string) (*jwt.Claims, error) {
	args := m.Called(tokenStr)
	return args.Get(0).(*jwt.Claims), args.Error(1)
}

func (m *MockJWTGenerator) RefreshAccessToken(refreshToken string) (string, error) {
	args := m.Called(refreshToken)
	return args.String(0), args.Error(1)
}

func TestUserServiceImpl_Register(t *testing.T) {
	tests := []struct {
		name           string
		user           *model.User
		mockSetup      func(*MockUserRepository, *MockJWTGenerator)
		expectedTokens []string
		expectedError  error
	}{
		{
			name: "successful registration",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
			},
			mockSetup: func(mur *MockUserRepository, mjg *MockJWTGenerator) {
				mur.On("CheckUserExists", "testuser", "test@example.com").Return(false, nil)
				mur.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)
				mjg.On("GenerateAccessToken", mock.AnythingOfType("int")).Return("access_token", nil)
				mjg.On("GenerateRefreshToken", mock.AnythingOfType("int")).Return("refresh_token", nil)
			},
			expectedTokens: []string{"access_token", "refresh_token"},
			expectedError:  nil,
		},
		{
			name: "user already exists",
			user: &model.User{
				Username: "existinguser",
				Email:    "existing@example.com",
				Password: "password",
			},
			mockSetup: func(mur *MockUserRepository, mjg *MockJWTGenerator) {
				mur.On("CheckUserExists", "existinguser", "existing@example.com").Return(true, nil)
			},
			expectedTokens: []string{"", ""},
			expectedError:  errors.New("用户名或邮箱已存在"),
		},
		{
			name: "error checking user existence",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
			},
			mockSetup: func(mur *MockUserRepository, mjg *MockJWTGenerator) {
				mur.On("CheckUserExists", "testuser", "test@example.com").Return(false, errors.New("database error"))
			},
			expectedTokens: []string{"", ""},
			expectedError:  errors.New("database error"),
		},
		{
			name: "error creating user",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
			},
			mockSetup: func(mur *MockUserRepository, mjg *MockJWTGenerator) {
				mur.On("CheckUserExists", "testuser", "test@example.com").Return(false, nil)
				mur.On("CreateUser", mock.AnythingOfType("*model.User")).Return(errors.New("creation error"))
			},
			expectedTokens: []string{"", ""},
			expectedError:  errors.New("creation error"),
		},
		{
			name: "error generating access token",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
			},
			mockSetup: func(mur *MockUserRepository, mjg *MockJWTGenerator) {
				mur.On("CheckUserExists", "testuser", "test@example.com").Return(false, nil)
				mur.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)
				mjg.On("GenerateAccessToken", mock.AnythingOfType("int")).Return("", errors.New("token error"))
			},
			expectedTokens: []string{"", ""},
			expectedError:  errors.New("token error"),
		},
		{
			name: "error generating refresh token",
			user: &model.User{
				Username: "testuser",
				Email:    "test@example.com",
				Password: "password",
			},
			mockSetup: func(mur *MockUserRepository, mjg *MockJWTGenerator) {
				mur.On("CheckUserExists", "testuser", "test@example.com").Return(false, nil)
				mur.On("CreateUser", mock.AnythingOfType("*model.User")).Return(nil)
				mjg.On("GenerateAccessToken", mock.AnythingOfType("int")).Return("access_token", nil)
				mjg.On("GenerateRefreshToken", mock.AnythingOfType("int")).Return("", errors.New("token error"))
			},
			expectedTokens: []string{"", ""},
			expectedError:  errors.New("token error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := new(MockUserRepository)
			mockJWT := new(MockJWTGenerator)
			tt.mockSetup(mockRepo, mockJWT)

			service := &UserServiceImpl{
				UserRepository: mockRepo,
				jwtManager:     mockJWT,
			}

			accessToken, refreshToken, err := service.Register(tt.user)

			if tt.expectedError != nil {
				assert.EqualError(t, err, tt.expectedError.Error())
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedTokens[0], accessToken)
			assert.Equal(t, tt.expectedTokens[1], refreshToken)

			mockRepo.AssertExpectations(t)
			mockJWT.AssertExpectations(t)
		})
	}
}

func TestUserServiceImpl_Login(t *testing.T) {
	mockRepo := new(MockUserRepository)
	mockJWT := new(MockJWTGenerator)

	mockRepo.On("GetUserByUsername", "testuser").Return(&model.User{Model: gorm.Model{ID: 1}, Username: "testuser", Password: "hashedpassword"}, nil)
	mockRepo.On("VerifyPassword", "hashedpassword", "password123").Return(nil)
	mockJWT.On("GenerateAccessToken", 1).Return("access_token", nil)
	mockJWT.On("GenerateRefreshToken", 1).Return("refresh_token", nil)

	service := &UserServiceImpl{
		UserRepository: mockRepo,
		jwtManager:     mockJWT,
	}

	accessToken, refreshToken, err := service.Login("testuser", "password123")

	assert.NoError(t, err)
	assert.Equal(t, "access_token", accessToken)
	assert.Equal(t, "refresh_token", refreshToken)

	mockRepo.AssertExpectations(t)
	mockJWT.AssertExpectations(t)
}

func TestUserServiceImpl_GetCurrentUser(t *testing.T) {
	mockRepo := new(MockUserRepository)

	mockRepo.On("GetUserByID", uint(1)).Return(&model.User{Model: gorm.Model{ID: 1}, Username: "testuser", Password: ""}, nil)

	service := &UserServiceImpl{
		UserRepository: mockRepo,
	}

	user, err := service.GetCurrentUser(1)

	assert.NoError(t, err)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "", user.Password)

	mockRepo.AssertExpectations(t)
}
