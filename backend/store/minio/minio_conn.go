package minio

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/2022AA/bytes-linked/backend/pkg/logging"
	"github.com/2022AA/bytes-linked/backend/setting"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/lifecycle"
)

var minioCli *minio.Client

// Client : 创建minio client对象
func Client() (*minio.Client, error) {
	if minioCli != nil {
		return minioCli, nil
	}
	minioCli, err := minio.New(setting.MinioSetting.MinioEndpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(setting.MinioSetting.MinioAccesskeyID, setting.MinioSetting.MinioAccessKeySecret, ""),
		Secure: setting.MinioSetting.MinioUseSSL,
	})
	if err != nil {
		logging.Error(err.Error())
		return nil, err
	}
	return minioCli, nil
}

// MakeBucket : 初始化bucket存储空间
func MakeBucket(bucketName, location string) (err error) {
	minioCli, err = Client()
	if err != nil {
		logging.Error(err)
		return err
	}
	ctx := context.Background()
	err = minioCli.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{Region: location})
	if err != nil {
		// 检查存储桶是否已经存在。
		exists, errBucketExists := minioCli.BucketExists(ctx, bucketName)
		if errBucketExists == nil {
			if exists {
				logging.Info(fmt.Sprintf("We already own bucket： %s\n", bucketName))
				return nil
			} else {
				logging.Error(err)
				return err
			}
		} else {
			logging.Error(errBucketExists)
			return errBucketExists
		}
	}
	return nil
}

// DelObject : 删除minio的文件
func DelObject(bucketName string, path string) (err error) {
	found, err := minioCli.BucketExists(context.Background(), bucketName)
	if err != nil {
		logging.Error(err)
		return
	}
	if !found {
		return errors.New(fmt.Sprintf("Bucket %s is not found! ", bucketName))
	}
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	err = minioCli.RemoveObject(context.Background(), bucketName, path, opts)
	if err != nil {
		logging.Error(err)
	}
	return
}

// PutObject : 提交文件到minio
func PutObject(bucketName string, path string, file *os.File) (err error) {
	fileStat, err := file.Stat()
	if err != nil {
		logging.Error(err)
		return
	}
	uploadInfo, err := minioCli.PutObject(context.Background(), bucketName,
		path, file, fileStat.Size(), minio.PutObjectOptions{ContentType: "application/octet-stream"})
	if err != nil {
		logging.Error(err)
		return
	}
	logging.Info("Successfully uploaded bytes: ", uploadInfo)
	return
}

// GetObject : 从minio获取文件
func GetObject(bucketName, path, fileName string) (*bytes.Buffer, error) {
	object, err := minioCli.GetObject(context.Background(), bucketName, path, minio.GetObjectOptions{})
	if err != nil {
		logging.Error(fmt.Sprintf("GetObject %s from bucket: %s fail. Err: %s", path, fileName, err))
		return nil, err
	}

	buf := bytes.NewBuffer(nil)
	if _, err = io.Copy(buf, object); err != nil {
		logging.Error(err)
		return nil, err
	}
	return buf, nil
}

// BuildLifecycleRule : 针对指定bucket设置生命周期规则
func BuildLifecycleRule(bucketName string) (err error) {
	ctx := context.Background()
	config := lifecycle.NewConfiguration()
	config.Rules = []lifecycle.Rule{{
		ID:     "expire-bucket",
		Status: "Enabled",
		Expiration: lifecycle.Expiration{
			Days: 365,
		},
	},
	}
	err = minioCli.SetBucketLifecycle(ctx, bucketName, config)
	if err != nil {
		logging.Error(err)
	}
	return err
}
