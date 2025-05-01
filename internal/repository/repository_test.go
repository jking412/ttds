package repository

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/db"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

func TestMain(m *testing.M) {
	// 初始化数据库
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		"root",
		"123456",
		"localhost",
		3306,
		"ttds",
	)
	db.InitDB(dsn, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 👈 打印所有 SQL
	})
	NewUserRepository(db.DB)
	NewCourseRepository(db.DB)
	// 执行测试
	m.Run()
}

// Helper function to generate random string
func randomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func generateMockUser() *model.User {
	randSuffix := strconv.Itoa(rand.Intn(10000))
	user := &model.User{
		Username: "testuser_" + randSuffix,
		Email:    fmt.Sprintf("test_%s@example.com", randSuffix),
		Password: "password123", // Note: In real tests, hash the password
	}
	db.DB.Create(user)
	return user
}

func generateMockCourse(withDetails bool) *model.Course {
	course := &model.Course{
		Title:       "测试课程 " + randomString(5),
		Description: "这是一个测试课程",
		CoverImage:  "test.jpg",
		Category:    "测试",
		Chapters:    make([]model.Chapter, 0),
		References:  make([]model.CourseReference, 0),
	}
	db.DB.Create(course)

	if withDetails {
		// 创建参考资料
		reference := generateMockReference(course.ID)
		course.References = append(course.References, *reference)

		// 创建章节和小节
		chapter := generateMockChapter(course.ID, 1, true)
		course.Chapters = append(course.Chapters, *chapter)
	}

	return course
}

func generateMockReference(courseID uint) *model.CourseReference {
	reference := &model.CourseReference{
		CourseID:    courseID,
		Title:       "测试参考资料 " + randomString(5),
		Type:        "pdf",
		URL:         randomString(10) + ".pdf",
		Description: "这是一个测试参考资料",
	}
	db.DB.Create(reference)
	return reference
}

func generateMockChapter(courseID uint, order uint, withSections bool) *model.Chapter {
	chapter := &model.Chapter{
		Title:       fmt.Sprintf("测试章节 %d", order),
		Description: "这是一个测试章节",
		Order:       order,
		CourseID:    courseID,
		Sections:    make([]model.Section, 0),
	}
	db.DB.Create(chapter)

	if withSections {
		section1 := generateMockSection(chapter.ID, 1)
		section2 := generateMockSection(chapter.ID, 2)
		chapter.Sections = append(chapter.Sections, *section1, *section2)
	}
	return chapter
}

func generateMockSection(chapterID uint, order uint) *model.Section {
	section := &model.Section{
		Title:     fmt.Sprintf("测试小节 %d", order),
		Content:   "这是测试小节内容",
		Order:     order,
		ChapterID: chapterID,
	}
	db.DB.Create(section)
	return section
}

func generateMockUserSectionStatus(userID, sectionID uint, completed bool) *model.UserSectionStatus {
	status := "incomplete"
	if completed {
		status = "completed"
	}
	userStatus := &model.UserSectionStatus{
		UserID:    userID,
		SectionID: sectionID,
		Status:    status,
		Completed: completed,
	}
	db.DB.Create(userStatus)
	return userStatus
}

func cleanMockUser(user *model.User) {
	db.DB.Unscoped().Delete(user)
}

func cleanMockCourse(course *model.Course) {
	// Need to reload to get associated data if not passed in
	db.DB.Preload("Chapters.Sections").Preload("References").First(course, course.ID)

	for _, chapter := range course.Chapters {
		cleanMockChapter(&chapter)
	}
	for _, reference := range course.References {
		cleanMockCourseReference(&reference)
	}
	db.DB.Unscoped().Delete(course)
}

func cleanMockCourseReference(reference *model.CourseReference) {
	db.DB.Unscoped().Delete(reference)
}

func cleanMockChapter(chapter *model.Chapter) {
	// Need to reload to get associated data if not passed in
	db.DB.Preload("Sections").First(chapter, chapter.ID)
	for _, section := range chapter.Sections {
		cleanMockSection(&section)
	}
	db.DB.Unscoped().Delete(chapter)
}

func cleanMockSection(section *model.Section) {
	// Clean related statuses first if any
	db.DB.Unscoped().Where("section_id = ?", section.ID).Delete(&model.UserSectionStatus{})
	db.DB.Unscoped().Delete(section)
}

func cleanMockUserSectionStatus(status *model.UserSectionStatus) {
	db.DB.Unscoped().Delete(status)
}
