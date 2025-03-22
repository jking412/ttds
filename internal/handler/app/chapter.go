package app

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

// CreateChapterHandler 创建章节的处理函数
func CreateChapterHandler(c *gin.Context) {
	var chapter model.Chapter
	if err := c.ShouldBindJSON(&chapter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := repository.CreateChapter(&chapter); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, chapter)
}

// GetChapterByIDHandler 根据章节 ID 获取章节信息的处理函数
func GetChapterByIDHandler(c *gin.Context) {
	id := c.Param("id")
	var chapterID uint
	fmt.Sscanf(id, "%d", &chapterID)
	chapter, err := repository.GetChapterByID(chapterID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chapter)
}

// GetChaptersByCourseIDHandler 根据课程 ID 获取所有章节信息的处理函数
func GetChaptersByCourseIDHandler(c *gin.Context) {
	courseIDStr := c.Param("courseID")
	var courseID uint
	fmt.Sscanf(courseIDStr, "%d", &courseID)
	chapters, err := repository.GetChaptersByCourseID(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, chapters)
}