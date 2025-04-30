package configs

// jwt:
//
//	secret: ttds
//	refresh_secret: ttds-refresh
//	expires: 1m
//	refresh_expires: 24h
func setJWTConfig(appConfig *AppConfig) {
	appConfig.JWT.Secret = "ttds"
	appConfig.JWT.RefreshSecret = "ttds-refresh"
	appConfig.JWT.Expires = "1m"
	appConfig.JWT.RefreshExpires = "24h"
}
