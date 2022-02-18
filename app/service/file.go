package service

// @Author: 陈健航
// @Date: 2021/2/10 23:35
// @Description:

import (
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"mime/multipart"
	"time"
)

var File = newFileService()

type fileService struct {
	minio  *minio.Client
	bucket string
}

const (
	policyPublic  = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:GetBucketLocation\",\"s3:ListBucket\",\"s3:ListBucketMultipartUploads\"],\"Resource\":[\"arn:aws:s3:::%s\"]},{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:ListMultipartUploadParts\",\"s3:PutObject\",\"s3:AbortMultipartUpload\",\"s3:DeleteObject\",\"s3:GetObject\"],\"Resource\":[\"arn:aws:s3:::%s/*\"]}]}\n"
	policyPrivate = "{\"Version\":\"2012-10-17\",\"Statement\":[]}"
)

func newFileService() (f fileService) {
	endpoint := g.Cfg().GetString("minio.endpoint")
	accessKeyId := g.Cfg().GetString("minio.accessKeyId")
	secretAccessKey := g.Cfg().GetString("minio.secretAccessKey")
	location := g.Cfg().GetString("minio.location")
	bucket := g.Cfg().GetString("minio.bucket")
	// 初使化 minio client对象。
	m, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
	f = fileService{
		minio:  m,
		bucket: bucket,
	}

	ctx := context.Background()
	// 创建一个存储桶
	isExist, err := m.BucketExists(ctx, bucket)
	if err != nil {
		panic(err)
	}
	if !isExist {
		if err = m.MakeBucket(
			ctx,
			bucket,
			minio.MakeBucketOptions{Region: location},
		); err != nil {
			panic(bucket + " 存储桶创建失败" + err.Error())
		}
		//设置该存储桶策略
		if err = m.SetBucketPolicy(ctx, bucket, fmt.Sprintf(policyPublic, bucket, bucket)); err != nil {
			panic(err)
		}
	}
	return f
}

// UploadFile 上传图片
// @receiver s
// @params file 文件
// @params width 像素
// @return url
// @return err
// @date 2021-02-18 22:43:25
func (f *fileService) UploadFile(ctx context.Context, uploadFile *ghttp.UploadFile) (fileName string, err error) {
	// 打开文件
	file, err := uploadFile.Open()
	if err != nil {
		return "", err
	}
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)
	// 在文件名加入时间戳
	fileName = fmt.Sprintf("%s_%s", gtime.TimestampStr(), uploadFile.Filename)
	opt := minio.PutObjectOptions{ContentType: uploadFile.Header.Get("Content-Type")}
	size := uploadFile.Size
	// 上传
	if _, err = f.minio.PutObject(ctx, f.bucket, fileName, file, size, opt); err != nil {
		return "", err
	}
	// 返回fileName
	return fileName, nil
}

func (f *fileService) UploadFileAndGetUrl(ctx context.Context, uploadFile *ghttp.UploadFile) (url string, err error) {
	// 打开文件
	file, err := f.UploadFile(ctx, uploadFile)
	if err != nil {
		return "", err
	}
	url = fmt.Sprintf("%s/%s", f.minio.EndpointURL().Host, file)
	// 返回fileName
	return url, nil
}

// RemoveObject 根据url删除文件
// @receiver s
// @params url
// @return error
// @date 2021-02-28 16:32:12
func (f *fileService) RemoveObject(ctx context.Context, fileName string) error {
	if err := f.minio.RemoveObject(
		ctx,
		f.bucket,
		fileName,
		minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	return nil
}

func (f *fileService) GetObjectPresignedUrl(ctx context.Context, objectName string, expires time.Duration) (string, error) {
	url, err := f.minio.PresignedGetObject(ctx, f.bucket, objectName, expires, nil)
	if err != nil {
		return "", nil
	}
	return url.String(), err
}

func (f fileService) GetObjectUrl(_ context.Context, objectName string) string {
	if objectName == "" {
		return ""
	}
	return fmt.Sprintf("%s/%s/%s", f.minio.EndpointURL(), f.bucket, objectName)
}

func (f fileService) Get() {
	acl, err := f.minio.GetBucketPolicy(context.Background(), f.bucket)
	if err != nil {
		return
	}
	fmt.Println(acl)
}
