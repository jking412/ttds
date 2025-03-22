package model

import (
	"gorm.io/gorm"
)

// Course 课程模型
type Course struct {
	gorm.Model
	Title          string `gorm:"type:varchar(255)"`
	Description    string `gorm:"type:text"`
	Chapters       []Chapter
	ReferenceBooks []CourseReferenceBook
}

// Chapter 章节模型
type Chapter struct {
	gorm.Model
	CourseID uint   `gorm:"index"`
	Title    string `gorm:"type:varchar(255)"`
	Sections []Section
}

// Section 小节模型
type Section struct {
	gorm.Model
	ChapterID       uint   `gorm:"index"`
	Title           string `gorm:"type:varchar(255)"`
	TaskDescription string `gorm:"type:varchar(1023)"`
	ContainerInfo   string `gorm:"type:text"`
	UserStatus      []UserSectionStatus
}

// User 用户模型
type User struct {
	gorm.Model
	Name          string `gorm:"type:varchar(255)"`
	SectionStatus []UserSectionStatus
}

// UserSectionStatus 用户小节完成状态模型
type UserSectionStatus struct {
	gorm.Model
	UserID    uint `gorm:"index"`
	SectionID uint `gorm:"index"`
	Completed bool
}

// CourseReferenceBook 课程参考书籍模型
type CourseReferenceBook struct {
	gorm.Model
	CourseID  uint   `gorm:"index"`
	BookTitle string `gorm:"type:varchar(255)"`
	Author    string `gorm:"type:varchar(255)"`
}
