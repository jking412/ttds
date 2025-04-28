package app

import (
	"awesomeProject/internal/model"
	"awesomeProject/internal/repository"
	"awesomeProject/pkg/db"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// GetAllCoursesHandler 获取所有课程信息的处理函数
func GetAllCoursesHandler(c *gin.Context) {

	var courses []model.Course
	// 从缓存中获取课程信息
	coursesStr, err := db.RDB.Get(context.Background(), "courses").Result()
	if err == nil {
		err := json.Unmarshal([]byte(coursesStr), &courses)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(http.StatusOK, courses)
		return
	}

	courses, err = repository.GetAllCourses()

	// 获取课程信息中所有章节的小节信息
	for i := range courses {
		chapters, err := repository.GetChaptersByCourseID(courses[i].ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		for j := range chapters {
			sections, err := repository.GetSectionsByChapterID(chapters[j].ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			}
			chapters[j].Sections = sections
		}
		courses[i].Chapters = chapters
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 先json.Marshal
	coursesJSON, err := json.Marshal(courses)

	// 设置缓存
	cmd := db.RDB.Set(context.Background(), "courses", coursesJSON, time.Minute)
	if err := cmd.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, courses)
}

// CreateCourseHandler 创建课程的处理函数
func CreateCourseHandler(c *gin.Context) {
	var course model.Course
	if err := c.ShouldBindJSON(&course); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := repository.CreateCourse(&course); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, course)
}

// GetCourseByIDHandler 根据课程 ID 获取课程信息的处理函数
func GetCourseByIDHandler(c *gin.Context) {
	id := c.Param("id")
	var courseID uint
	fmt.Sscanf(id, "%d", &courseID)
	course, err := repository.GetCourseByID(courseID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, course)
}

// CreateCourseReferenceBookHandler 创建课程参考书籍记录的处理函数
func CreateCourseReferenceBookHandler(c *gin.Context) {
	var book model.CourseReferenceBook
	if err := c.ShouldBindJSON(&book); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := repository.CreateCourseReferenceBook(&book); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, book)
}

// GetCourseReferenceBooksByCourseIDHandler 根据课程 ID 获取所有参考书籍信息的处理函数
func GetCourseReferenceBooksByCourseIDHandler(c *gin.Context) {
	courseIDStr := c.Param("courseID")
	var courseID uint
	fmt.Sscanf(courseIDStr, "%d", &courseID)
	books, err := repository.GetCourseReferenceBooksByCourseID(courseID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, books)
}
