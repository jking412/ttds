package configs

// server:
//
//	address: :8080
func setServerConfig(appConfig *AppConfig) {
	appConfig.Server.Address = ":8080"
}
