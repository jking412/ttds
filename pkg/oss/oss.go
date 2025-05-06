package oss

type Manager interface {
	ListObjects() error
	UploadObject(objectName, filePath string) error
	GetObjectUrl(objectName string, expireSeconds int64) (string, error)
}
