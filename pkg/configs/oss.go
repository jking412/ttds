package configs

// oss:
//
//	endpoint: http://127.0.0.1:9000
//	access_key_id: Lb21CW94Xi3rdF294zrk
//	secret_access_key: XgL8HOw4Doz2orhiI1T9JkLeKm1HgUGX9Kbrl7WV
//	bucket_name: ttds-bucket
//	use_ssl: true
func setOssConfig(appConfig *AppConfig) {
	appConfig.Oss.Endpoint = "127.0.0.1:9000"
	appConfig.Oss.AccessKeyID = "Lb21CW94Xi3rdF294zrk"
	appConfig.Oss.SecretAccessKey = "XgL8HOw4Doz2orhiI1T9JkLeKm1HgUGX9Kbrl7WV"
	appConfig.Oss.BucketName = "ttds-bucket"
	appConfig.Oss.UseSSL = false

}
