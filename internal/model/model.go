package model

import (
	"gorm.io/gorm"
	"time"
)

// User 用户模型
type User struct {
	gorm.Model
	Username      string              `gorm:"type:varchar(50);uniqueIndex;not null"`  // 用户名，唯一
	Email         string              `gorm:"type:varchar(100);uniqueIndex;not null"` // 邮箱，唯一
	Password      string              `gorm:"type:varchar(255);not null"`             // 密码，哈希后的
	Avatar        string              `gorm:"type:varchar(255)"`                      // 头像URL，可选
	Bio           string              `gorm:"type:varchar(255)"`                      // 简介，可选
	SectionStatus []UserSectionStatus `gorm:"foreignKey:UserID"`                      // 用户学习状态
}

// UserSectionStatus 用户小节完成状态模型
type UserSectionStatus struct {
	gorm.Model
	UserID    uint   `gorm:"not null;index"`                        // 用户ID
	SectionID uint   `gorm:"not null;index"`                        // 小节ID
	Status    string `gorm:"type:varchar(20);default:'incomplete'"` // 学习状态，例如：incomplete/completed
	Completed bool
}

// Course 课程模型
type Course struct {
	gorm.Model
	Title       string            `gorm:"type:varchar(255);not null"`
	Description string            `gorm:"type:text"`
	CoverImage  string            `gorm:"type:varchar(255)"`
	Category    string            `gorm:"type:varchar(100)"`
	Chapters    []Chapter         `gorm:"foreignKey:CourseID"`
	References  []CourseReference `gorm:"foreignKey:CourseID"` // 课程参考资料
}

// Chapter 章节模型
type Chapter struct {
	gorm.Model
	Title       string    `gorm:"type:varchar(255);not null"` // 章节标题
	Description string    `gorm:"type:text"`                  // 章节简介
	Order       uint      `gorm:"not null"`                   // 章节排序编号
	CourseID    uint      `gorm:"not null;index"`             // 所属课程ID
	Sections    []Section `gorm:"foreignKey:ChapterID"`       // 关联小节
}

// Section 小节模型
type Section struct {
	gorm.Model
	Title      string              `gorm:"type:varchar(255);not null"` // 小节标题
	Content    string              `gorm:"type:text"`                  // 小节内容或描述
	Order      uint                `gorm:"not null"`                   // 小节排序编号
	ChapterID  uint                `gorm:"not null;index"`             // 所属章节ID
	TemplateID uint                `gorm:"index"`                      // （可选）关联的容器模板ID
	UserStatus []UserSectionStatus `gorm:"foreignKey:SectionID"`       // 用户完成状态
}

// CourseReference 课程参考资料模型
type CourseReference struct {
	gorm.Model
	CourseID    uint   `gorm:"not null;index"`             // 关联的课程ID
	Title       string `gorm:"type:varchar(255);not null"` // 资料标题
	Type        string `gorm:"type:varchar(50);not null"`  // 资料类型，例如：pdf、link、video
	URL         string `gorm:"type:varchar(500);not null"` // 资料链接或存储路径
	Description string `gorm:"type:text"`                  // 资料描述，解释资料用途，可选
}

// ContainerTemplate 容器模板模型
type ContainerTemplate struct {
	gorm.Model
	Name        string `gorm:"type:varchar(100);not null"` // 模板名称，例如 "Ubuntu + Docker"
	Description string `gorm:"type:text"`                  // 模板描述
	Image       string `gorm:"type:varchar(255);not null"` // 使用的容器镜像名，如 "ubuntu:20.04"
	DefaultCmd  string `gorm:"type:varchar(255)"`          // 容器默认启动命令，可选
	Volumes     string `gorm:"type:text"`                  // 挂载的卷，可以是JSON数组表示多个卷
	Ports       string `gorm:"type:text"`                  // 暴露的端口，JSON数组表示多个端口
	Envs        string `gorm:"type:text"`                  // 环境变量，JSON数组表示多个环境变量
	SUDOPass    string `gorm:"type:varchar(255)"`          // SUDO密码（如果有的话）
}

// ContainerInstance 容器实例模型
type ContainerInstance struct {
	gorm.Model
	UserID      uint      `gorm:"not null;index"`             // 关联的用户ID
	SectionID   uint      `gorm:"not null;index"`             // 关联的小节ID（在哪一节学习用的）
	TemplateID  uint      `gorm:"not null;index"`             // 使用的模板ID
	ContainerID string    `gorm:"type:varchar(255);not null"` // 容器实际ID（Docker/K8S管理用）
	Status      string    `gorm:"type:varchar(50);not null"`  // 状态：Pending / Running / Stopped / Error
	Name        string    `gorm:"type:varchar(100);not null"` // 容器名称，便于用户识别
	StartAt     time.Time `gorm:"type:timestamp"`             // 启动时间
	EndAt       time.Time `gorm:"type:timestamp"`             // 结束/销毁时间
	IPAddress   string    `gorm:"type:varchar(100)"`          // 容器分配的IP地址（如果有的话）
	Token       string    `gorm:"type:varchar(255)"`          // 容器访问令牌（如果有的话）
}

// ContainerScript 容器脚本模型
type ContainerScript struct {
	gorm.Model
	SectionID      uint   `gorm:"not null;index"`            // 关联的小节ID（在哪一节要检测）
	Content        string `gorm:"type:text;not null"`        // 脚本内容
	ExpectedOutput string `gorm:"type:text"`                 // 期望输出，比如包含某个文件
	MatchType      string `gorm:"type:varchar(50);not null"` // 匹配方式：contains / equals / regex
	Timeout        uint   `gorm:"default:10"`                // 超时时间（秒），默认10秒
	Description    string `gorm:"type:varchar(255)"`         // 检测说明，可选
}
