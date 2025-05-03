package utils

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/db"
	"fmt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"math/rand"
	"time"
)

func MockCourseData() {
	// 插入模拟用户数据
	password, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := model.User{
		Username: "apitestuser",
		Email:    randomString(10, 50) + "@example.com",
		Password: string(password),
	}

	if err := db.DB.Create(&user).Error; err != nil {
		logrus.Fatalf("failed to insert mock user data: %v", err)
	}

	// 插入模拟课程数据
	for i := 0; i < 5; i++ {
		course := model.Course{
			Title:       randomString(10, 255),
			Description: randomText(50, 500),
			CoverImage:  fmt.Sprintf("images/%s.jpg", randomString(5, 50)),
			Category:    randomString(5, 100),
		}
		if err := db.DB.Create(&course).Error; err != nil {
			logrus.Fatalf("failed to insert mock course data: %v", err)
		}

		// 插入模拟章节数据
		for j := 0; j < 3; j++ {
			chapter := model.Chapter{
				Title:       randomString(10, 255),
				Description: randomText(50, 500),
				Order:       uint(rand.Uint32()),
				CourseID:    course.ID,
			}
			if err := db.DB.Create(&chapter).Error; err != nil {
				logrus.Fatalf("failed to insert mock chapter data: %v", err)
			}

			// 插入模拟小节数据
			for k := 0; k < 2; k++ {
				section := model.Section{
					Title:     randomString(10, 255),
					Content:   randomText(100, 2000),
					Order:     uint(rand.Uint32()),
					ChapterID: chapter.ID,
				}
				if err := db.DB.Create(&section).Error; err != nil {
					logrus.Fatalf("failed to insert mock section data: %v", err)
				}
			}
		}

		// 插入模拟课程参考资料数据
		for l := 0; l < 2; l++ {
			reference := model.CourseReference{
				CourseID:    course.ID,
				Title:       randomString(10, 255),
				Type:        randomReferenceType(),
				URL:         fmt.Sprintf("https://example.com/%s", randomString(5, 50)),
				Description: randomText(50, 500),
			}
			if err := db.DB.Create(&reference).Error; err != nil {
				logrus.Fatalf("failed to insert mock reference data: %v", err)
			}
		}
	}

	// 插入模拟用户小节状态数据
	users := []model.User{}
	db.DB.Find(&users)
	sections := []model.Section{}
	db.DB.Find(&sections)
	for _, user := range users {
		for _, section := range sections {
			status := model.UserSectionStatus{
				UserID:    user.ID,
				SectionID: section.ID,
				Completed: rand.Intn(2) == 1,
			}
			if err := db.DB.Create(&status).Error; err != nil {
				logrus.Fatalf("failed to insert mock user section status data: %v", err)
			}
		}
	}

	logrus.Info("mock data inserted successfully")
}

// Mock data generators for each table

func mockCourses() map[string]interface{} {
	return map[string]interface{}{
		"title":       randomString(10, 255),
		"description": randomText(50, 500),
		"cover_image": fmt.Sprintf("images/%s.jpg", randomString(5, 50)),
		"category":    randomString(5, 100),
	}
}

func mockChapters(courseID uint64) map[string]interface{} {
	return map[string]interface{}{
		"title":       randomString(10, 255),
		"description": randomText(50, 500),
		"order":       rand.Uint64(),
		"course_id":   courseID,
	}
}

func mockSections(chapterID uint64) map[string]interface{} {
	return map[string]interface{}{
		"title":      randomString(10, 255),
		"content":    randomText(100, 2000),
		"order":      rand.Uint64(),
		"chapter_id": chapterID,
	}
}

func mockCourseReferences(courseID uint64) map[string]interface{} {
	return map[string]interface{}{
		"course_id":   courseID,
		"title":       randomString(10, 255),
		"type":        randomReferenceType(),
		"url":         fmt.Sprintf("https://example.com/%s", randomString(5, 500)),
		"description": randomText(50, 500),
	}
}

func mockCourseReferenceBooks(courseID uint64) map[string]interface{} {
	return map[string]interface{}{
		"course_id":  courseID,
		"book_title": randomString(10, 255),
		"author":     randomString(5, 255),
	}
}

// Helper functions

func randomString(minLen, maxLen int) string {
	length := minLen + rand.Intn(maxLen-minLen+1)
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 "
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func randomText(minLen, maxLen int) string {
	length := minLen + rand.Intn(maxLen-minLen+1)
	chars := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 ,.\n"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func randomReferenceType() string {
	types := []string{"video", "article", "paper", "website", "book"}
	return types[rand.Intn(len(types))]
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
