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

func InitDBData() {
	// 开始之前先查询是否存在名为 "admin" 的用户
	var count int64
	db.DB.Model(&model.User{}).Where("username = ?", "admin").Count(&count)
	if count > 0 {
		logrus.Info("admin user already exists, skipping mock data generation")
		return
	}

	password, _ := bcrypt.GenerateFromPassword([]byte("123456"), bcrypt.DefaultCost)
	// 创建admin用户
	admin := model.User{
		Username: "admin",
		Email:    "admin@example.com",
		Password: string(password),
	}

	if err := db.DB.Create(&admin).Error; err != nil {
		logrus.Fatalf("failed to insert admin user: %v", err)
	}

	// 创建两个课程
	for i := 1; i <= 2; i++ {
		course := model.Course{
			Title:       fmt.Sprintf("course-%d", i),
			Description: fmt.Sprintf("课程%d的描述", i),
			CoverImage:  "",
			Category:    "default",
		}
		if err := db.DB.Create(&course).Error; err != nil {
			logrus.Fatalf("failed to insert course data: %v", err)
		}

		// 每个课程3个章节
		for j := 1; j <= 3; j++ {
			chapter := model.Chapter{
				Title:       fmt.Sprintf("chapter-%d", j),
				Description: fmt.Sprintf("课程%d的第%d章", i, j),
				Order:       uint(j),
				CourseID:    course.ID,
			}
			if err := db.DB.Create(&chapter).Error; err != nil {
				logrus.Fatalf("failed to insert chapter data: %v", err)
			}

			// 每个章节3个小节
			for k := 1; k <= 3; k++ {
				section := model.Section{
					Title:     fmt.Sprintf("section-%d", k),
					Content:   fmt.Sprintf("课程%d第%d章的第%d节内容", i, j, k),
					Order:     uint(k),
					ChapterID: chapter.ID,
				}

				// 只有第一个章节的第一个小节有模板
				if i == 1 && j == 1 && k == 1 {
					template := model.ContainerTemplate{
						Name:        "test-container",
						Description: "测试容器",
						Image:       "os:test",
						Ports:       "3001:3000",
					}
					if err := db.DB.Create(&template).Error; err != nil {
						logrus.Fatalf("failed to insert container template: %v", err)
					}
					section.TemplateID = template.ID

					// 添加4个脚本
					for s := 1; s <= 4; s++ {
						script := model.ContainerScript{
							Order:      uint(s),
							TemplateID: template.ID,
							Content:    "sleep 5",
							Timeout:    10,
						}
						if s == 4 {
							script.Content = "exit 1"
						}
						if err := db.DB.Create(&script).Error; err != nil {
							logrus.Fatalf("failed to insert container script: %v", err)
						}
					}
				}

				if err := db.DB.Create(&section).Error; err != nil {
					logrus.Fatalf("failed to insert section data: %v", err)
				}

				// 设置小节状态：每个章节的第一个小节完成
				status := "incomplete"
				completed := false
				if k == 1 {
					status = "completed"
					completed = true
				}

				sectionStatus := model.UserSectionStatus{
					UserID:    admin.ID,
					SectionID: section.ID,
					Status:    status,
					Completed: completed,
				}
				if err := db.DB.Create(&sectionStatus).Error; err != nil {
					logrus.Fatalf("failed to insert section status: %v", err)
				}
			}
		}

		// 每个课程2个参考资料
		for r := 1; r <= 2; r++ {
			reference := model.CourseReference{
				CourseID:    course.ID,
				Title:       fmt.Sprintf("reference-%d", r),
				Type:        "book",
				Description: fmt.Sprintf("课程%d的第%d个参考资料", i, r),
			}
			if err := db.DB.Create(&reference).Error; err != nil {
				logrus.Fatalf("failed to insert course reference: %v", err)
			}
		}
	}
}

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
