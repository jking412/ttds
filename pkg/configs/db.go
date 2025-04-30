package configs

// db:
//
//	driver: mysql
//	username: root
//	password: 123456
//	host: 127.0.0.1
//	port: 3306
//	dbname: ttds
func setDBConfig(appConfig *AppConfig) {
	appConfig.DB.Driver = "mysql"
	appConfig.DB.Username = "root"
	appConfig.DB.Password = "123456"
	appConfig.DB.Host = "127.0.0.1"
	appConfig.DB.Port = "3306"
	appConfig.DB.DBName = "ttds"

}
