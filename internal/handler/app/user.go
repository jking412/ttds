package app

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"awesomeProject/internal/usecase"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"strconv"
)

var userService usecase.UserService

func initUserService() {
	userService = usecase.NewUserService()
}

const bzImage = "/home/skynesser/code/system_call/linux-6.2/arch/x86/boot/bzImage"

// CheckAnswer 检查bzImage是否存在
func CheckAnswer(c *gin.Context) {

	// 从前端取出sectionID
	sectionID := c.Query("section_id")
	var sectionIDUint uint
	fmt.Sscanf(sectionID, "%d", &sectionIDUint)

	// 检查bzImage是否存在
	_, err := os.Stat(bzImage)
	if err != nil {
		// 不存在则告知用户答案错误，注意请求是合法的
		c.JSON(http.StatusOK, gin.H{"answer": "false"})
		return
	}

	// 存在则告知用户答案正确，注意请求是合法的，并修改userSectionStatus
	// 修改userSectionStatus，userID为1，sectionID为sectionIDUint
	err = repository.UpdateUserSectionStatus(&model.UserSectionStatus{
		UserID:    1,
		SectionID: sectionIDUint,
		Completed: true,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"answer": "true"})
}

// RegisterHandler 用户注册处理函数
func RegisterHandler(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据", "details": err.Error()})
		return
	}

	accessToken, refreshToken, err := userService.Register(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	user.Password = ""
	c.JSON(http.StatusCreated, gin.H{
		"user":          user,
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"message":       "注册成功",
	})
}

// LoginHandler 用户登录处理函数
func LoginHandler(c *gin.Context) {
	type LoginRequest struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "无效的请求数据", "details": err.Error()})
		return
	}

	accessToken, refreshToken, err := userService.Login(req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
		"message":       "登录成功",
	})
}

// LogoutHandler 用户注销处理函数
func LogoutHandler(c *gin.Context) {
	if err := userService.Logout(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "注销成功",
	})
}

// GetCurrentUserHandler 获取当前登录用户信息的处理函数
func GetCurrentUserHandler(c *gin.Context) {
	userIDInterface, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "未授权"})
		return
	}

	userID, err := strconv.ParseUint(fmt.Sprintf("%v", userIDInterface), 10, 64)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "服务器错误", "details": err.Error()})
		return
	}

	user, err := userService.GetCurrentUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "用户不存在", "details": err.Error()})
		return
	}

	user.Password = ""
	c.JSON(http.StatusOK, user)
}
