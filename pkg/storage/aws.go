package storage

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io"
)

type S3Uploader struct {
	s3     *s3.S3
	bucket string
}

func NewS3Uploader(region, bucket string) *S3Uploader {
	sess := session.Must(session.NewSession(&aws.Config{
		Region: aws.String(region),
	}))
	return &S3Uploader{
		s3:     s3.New(sess),
		bucket: bucket,
	}
}

func (u *S3Uploader) UploadImage(file io.ReadSeeker, fileName string) (string, error) {
	return u.UploadFile(file, fileName)
}

func (u *S3Uploader) UploadFile(file io.ReadSeeker, fileName string) (string, error) {

	_, err := u.s3.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(u.bucket),
		Key:    aws.String(fileName),
		Body:   file,
		ACL:    aws.String("public-read"),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", u.bucket, fileName), nil
}
