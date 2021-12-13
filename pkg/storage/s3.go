package storage

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"io"
)

type s3Storage struct {
	bucket, key string
	client      s3manageriface.UploaderAPI
}

func (s s3Storage) Save(ctx context.Context, fileName string, src io.Reader) error {
	_, err := s.client.UploadWithContext(ctx, &s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(fmt.Sprintf("%s/%s", s.key, fileName)),
		Body:   src,
	})
	return err
}

func NewS3Storage(bucket, key string) *s3Storage {
	sess := session.Must(session.NewSession())
	uploader := s3manager.NewUploader(sess, func(u *s3manager.Uploader) {
		// Define a strategy that will buffer 25 MiB in memory
		u.BufferProvider = s3manager.NewBufferedReadSeekerWriteToPool(25 * 1024 * 1024)
	})
	return &s3Storage{bucket, key, uploader}
}
