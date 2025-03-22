package app

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateSectionHandler 创建小节的处理函数
func CreateSectionHandler(c *gin.Context) {
	var section model.Section
	if err := c.ShouldBindJSON(&section); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := repository.CreateSection(&section); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, section)
}

// GetSectionByIDHandler 根据小节 ID 获取小节信息的处理函数
func GetSectionByIDHandler(c *gin.Context) {
	id := c.Param("id")
	var sectionID uint
	fmt.Sscanf(id, "%d", &sectionID)
	section, err := repository.GetSectionByID(sectionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, section)
}

// GetSectionsByChapterIDHandler 根据章节 ID 获取所有小节信息的处理函数
func GetSectionsByChapterIDHandler(c *gin.Context) {
	chapterIDStr := c.Param("chapterID")
	var chapterID uint
	fmt.Sscanf(chapterIDStr, "%d", &chapterID)
	sections, err := repository.GetSectionsByChapterID(chapterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, sections)
}
