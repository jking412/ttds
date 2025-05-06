package oss

import (
	"awesomeProject/pkg/configs"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/url"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	client        *MinioClient
	minioSyncOnce sync.Once
	_             Manager = (*MinioClient)(nil)
)

func NewOssClient() Manager {
	minioSyncOnce.Do(func() {
		var err error
		client, err = NewMinioClient(
			configs.GetConfig().Oss.Endpoint,
			configs.GetConfig().Oss.AccessKeyID,
			configs.GetConfig().Oss.SecretAccessKey,
			configs.GetConfig().Oss.BucketName,
			configs.GetConfig().Oss.UseSSL,
		)
		if err != nil {
			logrus.Errorf("Failed to create Minio client: %v", err)
		}
	})
	return client
}

type MinioClient struct {
	Client *minio.Client
	Bucket string
	Ctx    context.Context
}

func NewMinioClient(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioClient, error) {
	ctx := context.Background()
	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}

	// 确保 bucket 存在
	exists, err := minioClient.BucketExists(ctx, bucket)
	if err != nil {
		return nil, err
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, err
		}
		logrus.Infof("Created bucket: %s\n", bucket)
	}

	return &MinioClient{Client: minioClient, Bucket: bucket, Ctx: ctx}, nil
}

func (oss *MinioClient) UploadObject(objectName, filePath string) error {
	_, err := oss.Client.FPutObject(oss.Ctx, oss.Bucket, objectName, filePath, minio.PutObjectOptions{})
	if err != nil {
		return err
	}
	fmt.Printf("Uploaded: %s\n", objectName)
	return nil
}

func (oss *MinioClient) GetObjectUrl(objectName string, expireSeconds int64) (string, error) {
	// 返回
	reqParams := make(url.Values)
	presignedURL, err := oss.Client.PresignedGetObject(
		oss.Ctx, oss.Bucket, objectName, time.Duration(expireSeconds)*time.Second, reqParams,
	)
	if err != nil {
		return "", err
	}
	return presignedURL.String(), nil
}

func (oss *MinioClient) ListObjects() error {
	for object := range oss.Client.ListObjects(oss.Ctx, oss.Bucket, minio.ListObjectsOptions{Recursive: true}) {
		if object.Err != nil {
			return object.Err
		}
		fmt.Printf("Found object: %s (%d bytes)\n", object.Key, object.Size)
	}
	return nil
}
