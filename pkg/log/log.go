package log

import (
	"awesomeProject/pkg/configs"
	"github.com/sirupsen/logrus"
	"io"
	"os"
)

func init() {

}

func InitLog() {
	level := configs.GetConfig().Log.Level
	logPath := configs.GetConfig().Log.Path
	fileName := configs.GetConfig().Log.Filename

	switch level {
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "warn":
		logrus.SetLevel(logrus.WarnLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.SetReportCaller(true)

	// 设置日志格式（比如使用 JSON 格式）
	logrus.SetFormatter(&logrus.TextFormatter{})

	logOutputPath := logPath + "/" + fileName

	// 设置日志输出到文件
	file, err := os.OpenFile(logOutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		logrus.Fatalf("Failed to open log file: %v", err)
	}

	multiWriter := io.MultiWriter(os.Stdout, file)

	logrus.SetOutput(multiWriter)
}
