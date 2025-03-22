package app

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

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

// CreateUserHandler 创建用户的处理函数
func CreateUserHandler(c *gin.Context) {
	var user model.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := repository.CreateUser(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, user)
}

// GetUserByIDHandler 根据用户 ID 获取用户信息的处理函数
func GetUserByIDHandler(c *gin.Context) {
	id := c.Param("id")
	var userID uint
	fmt.Sscanf(id, "%d", &userID)
	user, err := repository.GetUserByID(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, user)
}

// CreateUserSectionStatusHandler 创建用户小节完成状态记录的处理函数
func CreateUserSectionStatusHandler(c *gin.Context) {
	var status model.UserSectionStatus
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := repository.CreateUserSectionStatus(&status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, status)
}

// GetUserSectionStatusByUserAndSectionIDHandler 根据用户 ID 和小节 ID 获取用户小节完成状态的处理函数
func GetUserSectionStatusByUserAndSectionIDHandler(c *gin.Context) {
	userIDStr := c.Param("userID")
	sectionIDStr := c.Param("sectionID")
	var userID, sectionID uint
	fmt.Sscanf(userIDStr, "%d", &userID)
	fmt.Sscanf(sectionIDStr, "%d", &sectionID)
	status, err := repository.GetUserSectionStatusByUserAndSectionID(userID, sectionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}

// UpdateUserSectionStatusHandler 更新用户小节完成状态的处理函数
func UpdateUserSectionStatusHandler(c *gin.Context) {
	userIDStr := c.Param("userID")
	sectionIDStr := c.Param("sectionID")
	var userID, sectionID uint
	fmt.Sscanf(userIDStr, "%d", &userID)
	fmt.Sscanf(sectionIDStr, "%d", &sectionID)
	var status model.UserSectionStatus
	if err := c.ShouldBindJSON(&status); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	status.UserID = userID
	status.SectionID = sectionID
	if err := repository.UpdateUserSectionStatus(&status); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, status)
}
