package configs

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"sync"
)

func init() {
	defaultConfigFuncList = make([]setDefaultConfigFunc, 0)
	defaultConfigFuncList = append(defaultConfigFuncList, setEnvConfig)
	defaultConfigFuncList = append(defaultConfigFuncList, setServerConfig)
	defaultConfigFuncList = append(defaultConfigFuncList, setDBConfig)
	defaultConfigFuncList = append(defaultConfigFuncList, setRedisConfig)
	defaultConfigFuncList = append(defaultConfigFuncList, setJWTConfig)
	defaultConfigFuncList = append(defaultConfigFuncList, setLogConfig)
}

// AppConfig
// viper使用mapstructure来解析配置文件
type AppConfig struct {
	Env string `mapstructure:"env"`

	Server struct {
		Address string `mapstructure:"address"`
	} `mapstructure:"server"`

	DB struct {
		Driver   string `mapstructure:"driver"`
		Username string `mapstructure:"username"`
		Password string `mapstructure:"password"`
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		DBName   string `mapstructure:"dbname"`
	} `mapstructure:"db"`

	Redis struct {
		Host     string `mapstructure:"host"`
		Port     string `mapstructure:"port"`
		Password string `mapstructure:"password"`
		DB       int    `mapstructure:"db"`
	} `mapstructure:"redis"`

	Oss struct {
		Endpoint        string `mapstructure:"endpoint"`
		AccessKeyID     string `mapstructure:"access_key_id"`
		SecretAccessKey string `mapstructure:"secret_access_key"`
		UseSSL          bool   `mapstructure:"use_ssl"`
		BucketName      string `mapstructure:"bucket_name"`
	} `mapstructure:"oss"`

	JWT struct {
		Secret         string `mapstructure:"secret"`
		RefreshSecret  string `mapstructure:"refresh_secret"`
		Expires        string `mapstructure:"expires"`         // 建议在使用时转换为 time.Duration
		RefreshExpires string `mapstructure:"refresh_expires"` // 同上
	} `mapstructure:"jwt"`

	Log struct {
		Level    string `mapstructure:"level"`
		Path     string `mapstructure:"path"`
		Filename string `mapstructure:"filename"`
	} `mapstructure:"log"`
}

var (
	cfg  *AppConfig
	once sync.Once
)

func Init() error {
	var err error
	once.Do(func() {
		newViper := viper.New()
		newViper.SetConfigName("config")
		newViper.SetConfigType("yaml")
		newViper.AddConfigPath("config")
		newViper.AddConfigPath(".")

		err = newViper.ReadInConfig()
		if err != nil {
			logrus.Warn("Error reading config: ", err)
		}

		// bind 到结构体
		appConfig := defaultConfig()
		err = newViper.Unmarshal(appConfig)
		if err != nil {
			return
		}

		cfg = appConfig
	})
	return err
}

// GetConfig 全局获取配置
func GetConfig() *AppConfig {
	return cfg
}

type setDefaultConfigFunc func(appConfig *AppConfig)

var defaultConfigFuncList []setDefaultConfigFunc

func defaultConfig() *AppConfig {

	appConfig := new(AppConfig)

	for _, fn := range defaultConfigFuncList {
		fn(appConfig)
	}

	return appConfig
}
