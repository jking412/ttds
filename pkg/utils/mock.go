package utils

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/db"
	"github.com/go-faker/faker/v4"
	"github.com/sirupsen/logrus"
)

func InsertMockData() {
	// 查询course表中是否有数据
	var count int64
	if err := db.DB.Model(&model.Course{}).Count(&count).Error; err != nil {
		logrus.Fatalf("failed to count courses: %v", err)
	}

	if count > 0 {
		logrus.Info("course table already has data, skipping mock data insertion")
		return
	}

	// 插入模拟课程数据
	course := model.Course{
		Title:       "操作系统",
		Description: "操作系统是计算机系统的核心，它管理计算机的硬件资源和软件资源，为用户提供一个良好的运行环境。",
	}

	courses := []model.Course{course, {
		Title:       "计算机网络",
		Description: "计算机网络是计算机系统之间的通信系统，它使得计算机能够相互通信和共享资源。",
	}}

	chapter := model.Chapter{
		Title: "添加系统调用",
	}

	chapters := []model.Chapter{chapter, {
		Title: "编译内核",
	}}

	section := model.Section{
		Title:           "编译内核",
		TaskDescription: "在这个小节中，你将学习如何编译内核",
		ContainerInfo:   "os:base",
	}

	sections := []model.Section{section, {
		Title:           "编译内核",
		TaskDescription: "在这个小节中，你将学习如何编译内核",
		ContainerInfo:   "os:base",
	}}

	user := model.User{
		Name: "admin",
	}

	// 插入模拟课程数据
	if err := db.DB.Create(&courses).Error; err != nil {
		logrus.Fatalf("failed to insert courses: %v", err)
	}

	// 插入模拟章节数据
	for i := range chapters {
		chapters[i].CourseID = courses[0].ID
	}
	if err := db.DB.Create(&chapters).Error; err != nil {
		logrus.Fatalf("failed to insert chapters: %v", err)
	}

	// 插入模拟小节数据
	for j := range sections {
		sections[j].ChapterID = chapters[0].ID
	}
	if err := db.DB.Create(&sections).Error; err != nil {
		logrus.Fatalf("failed to insert sections: %v", err)
	}

	// 插入模拟用户数据
	if err := db.DB.Create(&user).Error; err != nil {
		logrus.Fatalf("failed to insert mock user data: %v", err)
	}

	// 插入模拟用户Section完成情况数据
	userSectionStatus := []model.UserSectionStatus{
		{
			UserID:    user.ID,
			SectionID: sections[0].ID,
			Completed: false,
		},
		{
			UserID:    user.ID,
			SectionID: sections[1].ID,
			Completed: false,
		},
	}

	if err := db.DB.Create(&userSectionStatus).Error; err != nil {
		logrus.Fatalf("failed to insert mock user section status data: %v", err)
	}

}

func InsertMockData1() {
	// 插入模拟课程数据
	for i := 0; i < 5; i++ {
		var course model.Course
		if err := faker.FakeData(&course); err != nil {
			logrus.Fatalf("failed to generate mock course data: %v", err)
		}
		if err := db.DB.Create(&course).Error; err != nil {
			logrus.Fatalf("failed to insert mock course data: %v", err)
		}

		// 插入模拟章节数据
		for j := 0; j < 3; j++ {
			var chapter model.Chapter
			if err := faker.FakeData(&chapter); err != nil {
				logrus.Fatalf("failed to generate mock chapter data: %v", err)
			}
			chapter.CourseID = course.ID
			if err := db.DB.Create(&chapter).Error; err != nil {
				logrus.Fatalf("failed to insert mock chapter data: %v", err)
			}

			// 插入模拟小节数据
			for k := 0; k < 2; k++ {
				var section model.Section
				if err := faker.FakeData(&section); err != nil {
					logrus.Fatalf("failed to generate mock section data: %v", err)
				}
				section.ChapterID = chapter.ID
				if err := db.DB.Create(&section).Error; err != nil {
					logrus.Fatalf("failed to insert mock section data: %v", err)
				}
			}
		}

		// 插入模拟参考书籍数据
		for l := 0; l < 2; l++ {
			var book model.CourseReferenceBook
			if err := faker.FakeData(&book); err != nil {
				logrus.Fatalf("failed to generate mock book data: %v", err)
			}
			book.CourseID = course.ID
			if err := db.DB.Create(&book).Error; err != nil {
				logrus.Fatalf("failed to insert mock book data: %v", err)
			}
		}
	}

	// 插入模拟用户数据
	for m := 0; m < 5; m++ {
		var user model.User
		if err := faker.FakeData(&user); err != nil {
			logrus.Fatalf("failed to generate mock user data: %v", err)
		}
		if err := db.DB.Create(&user).Error; err != nil {
			logrus.Fatalf("failed to insert mock user data: %v", err)
		}

		// 插入模拟用户小节状态数据
		sections := []model.Section{}
		db.DB.Find(&sections)
		for _, section := range sections {
			var status model.UserSectionStatus
			if err := faker.FakeData(&status); err != nil {
				logrus.Fatalf("failed to generate mock user section status data: %v", err)
			}
			status.UserID = user.ID
			status.SectionID = section.ID
			if err := db.DB.Create(&status).Error; err != nil {
				logrus.Fatalf("failed to insert mock user section status data: %v", err)
			}
		}
	}

	logrus.Info("mock data inserted successfully")
}
