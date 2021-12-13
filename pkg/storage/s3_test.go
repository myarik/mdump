package storage

import (
	"context"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/aws/aws-sdk-go/service/s3/s3manager/s3manageriface"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

type uploadTest struct {
	S3URI     string
	Region    string
	Bucket    string
	Prefix    string
	AccessKey string
	SecretKey string
}

var s3TestData = uploadTest{
	S3URI:     "s3://test-bucket/test/prefix",
	Region:    "eu-central-1",
	Bucket:    "test-bucket",
	Prefix:    "prefix",
	AccessKey: "FAKE_KEY",
	SecretKey: "FAKE_SECRET",
}

const (
	key    = "testkey"
	bucket = "test-bucket"
)

func TestS3Storage(t *testing.T) {
	content := strings.NewReader("my request")
	inputUploader := mockS3{}
	s := s3Storage{bucket, key, inputUploader}
	err := s.Save(context.Background(), "test", content)
	assert.NoError(t, err)
	errStorage := s3Storage{"error_bucket", key, inputUploader}
	err = errStorage.Save(context.Background(), "test", content)
	assert.Error(t, err, "something went wrong")
}

type mockS3 struct {
	s3manageriface.UploaderAPI
}

func (m mockS3) Upload(input *s3manager.UploadInput, _ ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	switch *input.Bucket {
	case "error_bucket":
		return nil, errors.New("something went wrong")
	default:
		return nil, nil
	}
}

func (m mockS3) UploadWithContext(ctx aws.Context, input *s3manager.UploadInput, _ ...func(*s3manager.Uploader)) (*s3manager.UploadOutput, error) {
	switch *input.Bucket {
	case "error_bucket":
		return nil, errors.New("something went wrong")
	default:
		return nil, nil
	}
}
