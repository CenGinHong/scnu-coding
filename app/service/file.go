package service

// @Author: 陈健航
// @Date: 2021/2/10 23:35
// @Description:

import (
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"io"
	"mime/multipart"
	"scnu-coding/app/utils"
)

var File = newFileService()

type fileService struct {
	minio        *minio.Client
	bucketName   string
	notUsedCache *utils.MyCache
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
	location := g.Cfg().GetString("minio.location")
	bucketName := g.Cfg().GetString("minio.bucketName")
	// 初使化 minio client对象。
	m, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyId, secretAccessKey, ""),
		Secure: false,
	})
	if err != nil {
		panic(err)
	}
	f = fileService{
		minio:        m,
		bucketName:   bucketName,
		notUsedCache: utils.NewMyCache(),
	}
	ctx := context.Background()
	// 创建一个存储桶
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
		if err = m.SetBucketPolicy(ctx, bucketName, fmt.Sprintf(policyReadOnly, bucketName, bucketName)); err != nil {
			panic(err)
		}
	}
	return f
}

func (f *fileService) uploadToMinio(ctx context.Context, uploadName string, file io.Reader, uploadSize int64, opts minio.PutObjectOptions) (err error) {
	// 上传
	if _, err = f.minio.PutObject(ctx, f.bucketName, uploadName, file, uploadSize, opts); err != nil {
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
	ext := gfile.Ext(uploadFile.Filename)
	fileName = gstr.StrTillEx(uploadFile.Filename, ext) + gtime.TimestampStr() + ext
	contentType := uploadFile.Header.Get("Content-Type")
	size := uploadFile.Size
	if err = f.uploadToMinio(ctx, fileName, file, size,
		minio.PutObjectOptions{ContentType: contentType}); err != nil {
		return "", err
	}
	// 返回fileName
	return fileName, nil
}

// RemoveObject 根据url删除文件
// @receiver s
// @params url
// @return error
// @date 2021-02-28 16:32:12
func (f *fileService) RemoveObject(ctx context.Context, fileName string) error {
	if err := f.minio.RemoveObject(
		ctx,
		f.bucketName,
		fileName,
		minio.RemoveObjectOptions{}); err != nil {
		return err
	}
	return nil
}

func (f *fileService) GetMinioAddr(_ context.Context, oriAddr string) (addr string) {
	addr = g.Cfg().GetString("minio.protocol") + "://" + g.Cfg().GetString("minio.endpoint") + "/" +
		f.bucketName + "/" + oriAddr
	return addr
}
