package app

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/tsatke/verbose-broccoli/internal/app/config"
)

//go:generate mockery --inpackage --testonly --case snake --name s3StorageClientAPI --filename s3_storage_mock_test.go

type s3StorageClientAPI interface {
	PutObject(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
	GetObject(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
	HeadObject(context.Context, *s3.HeadObjectInput, ...func(*s3.Options)) (*s3.HeadObjectOutput, error)
	DeleteObject(context.Context, *s3.DeleteObjectInput, ...func(*s3.Options)) (*s3.DeleteObjectOutput, error)
}

type S3Storage struct {
	bucket string
	client s3StorageClientAPI
}

func NewS3Storage(cfg config.Config) *S3Storage {
	return &S3Storage{
		bucket: cfg.GetString(config.AWSS3Bucket),
		client: s3.NewFromConfig(cfg.AWS),
	}
}

func (s *S3Storage) Create(docID DocID, rd io.Reader) error {
	_, err := s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(string(docID)),
	})
	if err == nil {
		return fmt.Errorf("object already exists")
	} else if !smithyCodeIs(err, "NotFound") {
		return fmt.Errorf("head object: %w", err)
	}

	_, err = s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(string(docID)),
		Body:   rd,
	})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

func (s *S3Storage) Read(docID DocID) (io.ReadCloser, error) {
	res, err := s.client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(string(docID)),
	})
	if err != nil {
		return nil, fmt.Errorf("get object: %w", err)
	}
	return res.Body, nil
}

func (s *S3Storage) Update(docID DocID, rd io.Reader) error {
	_, err := s.client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(string(docID)),
	})
	if err != nil {
		return fmt.Errorf("head object: %w", err)
	}

	_, err = s.client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(string(docID)),
		Body:   rd,
	})
	if err != nil {
		return fmt.Errorf("put object: %w", err)
	}
	return nil
}

func (s *S3Storage) Delete(docID DocID) error {
	_, err := s.client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(string(docID)),
	})
	if err != nil {
		return fmt.Errorf("delete object: %w", err)
	}
	return nil
}

func smithyCodeIs(err error, expected string) bool {
	var ae smithy.APIError
	if errors.As(err, &ae) {
		if ae.ErrorCode() == expected {
			return true
		}
	}
	return false
}
