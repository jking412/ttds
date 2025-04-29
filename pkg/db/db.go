package db

import (
	"awesomeProject/internal/model"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB(dsn string) {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
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
		&model.CourseReferenceBook{},
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
		viper.GetString("db.username"),
		viper.GetString("db.password"),
		viper.GetString("db.host"),
		viper.GetString("db.port"),
		viper.GetString("db.dbname"),
	)
	return dsn
}
