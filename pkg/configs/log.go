package configs

// log:
//
//	level: info
//	path: ./logs
//	filename: ttds.log
func setLogConfig(appConfig *AppConfig) {
	appConfig.Log.Level = "info"
	appConfig.Log.Path = "./logs"
	appConfig.Log.Filename = "ttds.log"
}
