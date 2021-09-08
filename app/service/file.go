package service

// @Author: 陈健航
// @Date: 2021/2/10 23:35
// @Description:

import (
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"mime/multipart"
	"scnu-coding/app/dao"
)

var File = newFileService()

type fileService struct {
	minio      *minio.Client
	bucketName string
}

const (
	policyReadOnly     = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:GetBucketLocation\",\"s3:ListBucket\"],\"Resource\":[\"arn:aws:s3:::%s\"]},{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:GetObject\"],\"Resource\":[\"arn:aws:s3:::%s/*\"]}]}"
	policyWriteOnly    = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:GetBucketLocation\",\"s3:ListBucketMultipartUploads\"],\"Resource\":[\"arn:aws:s3:::%s\"]},{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:AbortMultipartUpload\",\"s3:DeletePic\",\"s3:ListMultipartUploadParts\",\"s3:PutObject\"],\"Resource\":[\"arn:aws:s3:::%s/*\"]}]}"
	policyWriteAndRead = "{\"Version\":\"2012-10-17\",\"Statement\":[{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:ListBucket\",\"s3:ListBucketMultipartUploads\",\"s3:GetBucketLocation\"],\"Resource\":[\"arn:aws:s3:::%s\"]},{\"Effect\":\"Allow\",\"Principal\":{\"AWS\":[\"*\"]},\"Action\":[\"s3:AbortMultipartUpload\",\"s3:DeletePic\",\"s3:GetObject\",\"s3:ListMultipartUploadParts\",\"s3:PutObject\"],\"Resource\":[\"arn:aws:s3:::%s/*\"]}]}"
)

func newFileService() (f fileService) {
	endpoint := g.Cfg().GetString("minio.endpoint")
	accessKeyId := g.Cfg().GetString("minio.accessKeyId")
	secretAccessKey := g.Cfg().GetString("minio.secretAccessKey")
	// 初使化 minio client对象。
	m, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
	f = fileService{
		minio:      m,
		bucketName: "scnu-coding",
	}
	location := "cn-north-1"
	ctx := context.Background()
	// 创建一个存储桶
	bucketName := f.bucketName
	isExist, err := m.BucketExists(ctx, bucketName)
	if err != nil {
		panic(err)
	}
	if !isExist {
		if err = m.MakeBucket(
			ctx,
			bucketName,
			minio.MakeBucketOptions{Region: location},
		); err != nil {
			panic(bucketName + " 存储桶创建失败" + err.Error())
		}
		//设置该存储桶策略
		if err = m.SetBucketPolicy(ctx, bucketName, fmt.Sprintf(policyWriteAndRead, bucketName, bucketName)); err != nil {
			panic(err)
		}
	}
	return f
}

func (receiver *fileService) uploadToMinio(ctx context.Context, uploadName string, file io.Reader, uploadSize int64, opts minio.PutObjectOptions) (err error) {
	// 上传
	if _, err = receiver.minio.PutObject(ctx, receiver.bucketName, uploadName, file, uploadSize, opts); err != nil {
		return err
	}
	return nil
}

// UploadFile 上传图片
// @receiver s
// @params file 文件
// @params width 像素
// @return url
// @return err
// @date 2021-02-18 22:43:25
func (receiver *fileService) UploadFile(ctx context.Context, uploadFile *ghttp.UploadFile) (fileId int64, err error) {
	// 文件类型仅支持下面几种类型
	// 编码
	file, err := uploadFile.Open()
	defer func(file multipart.File) {
		_ = file.Close()
	}(file)
	if err != nil {
		return 0, err
	}
	fileName := uploadFile.Filename
	contentType := uploadFile.Header.Get("Content-Type")
	size := uploadFile.Size
	if err = receiver.uploadToMinio(ctx, fileName, file, size,
		minio.PutObjectOptions{ContentType: contentType}); err != nil {
		return 0, err
	}
	url := fmt.Sprintf("%s%s/%s/%s",
		"http://",
		g.Cfg().GetString("minio.endpoint"),
		receiver.bucketName,
		uploadFile.Filename)
	fileId, err = dao.LocalFile.Ctx(ctx).Data(g.Map{
		dao.LocalFile.Columns.Filename:    fileName,
		dao.LocalFile.Columns.ContentType: uploadFile.Header.Get("Content-Type"),
		dao.LocalFile.Columns.Size:        size,
		dao.LocalFile.Columns.Url:         url,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	// 返回可直接访问的url
	return fileId, nil
}

// RemoveObject 根据url删除文件
// @receiver s
// @params url
// @return error
// @date 2021-02-28 16:32:12
func (receiver *fileService) RemoveObject(ctx context.Context, fileId string) error {
	filename, err := dao.LocalFile.Ctx(ctx).WherePri(fileId).Value(dao.LocalFile.Columns.Filename)
	if err != nil {
		return err
	}
	if err := receiver.minio.RemoveObject(
		context.Background(),
		receiver.bucketName,
		filename.String(),
		minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	return nil
}
