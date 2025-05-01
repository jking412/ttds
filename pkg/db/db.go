package db

import (
	"awesomeProject/internal/model"
	"awesomeProject/pkg/configs"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string, config *gorm.Config) {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), config)
	if err != nil {
		logrus.Fatalf("failed to connect database: %v", err)
	}
	logrus.Info("database connected successfully")

	// 自动迁移模型
	err = DB.AutoMigrate(
		&model.Course{},
		&model.Chapter{},
		&model.Section{},
		&model.User{},
		&model.UserSectionStatus{},
		&model.CourseReference{},
		&model.ContainerTemplate{},
		&model.ContainerInstance{},
		&model.ContainerScript{},
		// old
		&model.CourseReferenceBook{},
		&model.SectionRecord{},
	)
	if err != nil {
		logrus.Fatalf("failed to migrate database: %v", err)
	}
	logrus.Info("database migrated successfully")
}

func GenerateDsnFromConfig() string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		configs.GetConfig().DB.Username,
		configs.GetConfig().DB.Password,
		configs.GetConfig().DB.Host,
		configs.GetConfig().DB.Port,
		configs.GetConfig().DB.DBName,
	)
	return dsn
}
