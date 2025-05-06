package oss

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

//  access_key_id: Lb21CW94Xi3rdF294zrk
//  secret_access_key: XgL8HOw4Doz2orhiI1T9JkLeKm1HgUGX9Kbrl7WV
//  bucket_name: ttds-bucket
//  use_ssl: true

func TestMinIOIntegration(t *testing.T) {
	// 初始化MinIO客户端
	cli, err := NewMinioClient(
		"127.0.0.1:9000",
		"Lb21CW94Xi3rdF294zrk",
		"XgL8HOw4Doz2orhiI1T9JkLeKm1HgUGX9Kbrl7WV",
		"ttds-bucket",
		false)
	assert.NoError(t, err)
	assert.NotNil(t, cli)

	// 测试上传minio.go文件
	testObjectName := "minio.go"
	testFilePath := "./minio.go"
	err = cli.UploadObject(testObjectName, testFilePath)
	assert.NoError(t, err)

	// 测试获得对象URL
	var objectURL string
	objectURL, err = cli.GetObjectUrl(testObjectName, 600)
	assert.NoError(t, err)
	assert.NotEmpty(t, objectURL)
	t.Logf("Object URL: %s", objectURL)

	// 测试列出对象
	err = cli.ListObjects()
	assert.NoError(t, err)
}
