package configs

// redis:
//
//	host: 127.0.0.1
//	port: 6379
//	password:
//	db: 0
func setRedisConfig(appConfig *AppConfig) {
	appConfig.Redis.Host = "127.0.0.1"
	appConfig.Redis.Port = "6379"
	appConfig.Redis.Password = ""
	appConfig.Redis.DB = 0
}
