package app

import (
	"context"
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestS3StorageSuite(t *testing.T) {
	suite.Run(t, new(S3StorageTestSuite))
}

type S3StorageTestSuite struct {
	suite.Suite

	client  *mockS3StorageClientAPI
	storage *S3Storage
	bucket  string
}

func (suite *S3StorageTestSuite) SetupTest() {
	suite.bucket = "test-bucket"
	suite.client = &mockS3StorageClientAPI{}
	suite.storage = &S3Storage{
		bucket: suite.bucket,
		client: suite.client,
	}
}

func (suite *S3StorageTestSuite) TearDownTest() {
	suite.client.AssertExpectations(suite.T())
}

func (suite *S3StorageTestSuite) TestCreateNotExists() {
	suite.client.
		On("HeadObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.HeadObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(
			nil,
			&smithy.GenericAPIError{
				Code:    "NotFound",
				Message: "test message",
			},
		).
		Once()

	suite.client.
		On("PutObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.PutObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket) &&
					suite.NotNil(i.Body)
			}),
		).
		Return(nil, nil).
		Once()

	suite.NoError(suite.storage.Create("abc", strings.NewReader("content")))
}

func (suite *S3StorageTestSuite) TestCreateErrInPut() {
	suite.client.
		On("HeadObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.HeadObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(
			nil,
			&smithy.GenericAPIError{
				Code:    "NotFound",
				Message: "test message",
			},
		).
		Once()

	suite.client.
		On("PutObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.PutObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket) &&
					suite.NotNil(i.Body)
			}),
		).
		Return(
			nil,
			&smithy.GenericAPIError{
				Code:    "SomeError",
				Message: "test message",
			},
		).
		Once()

	suite.Error(suite.storage.Create("abc", strings.NewReader("content")))
}

func (suite *S3StorageTestSuite) TestCreateAlreadyExists() {
	suite.client.
		On("HeadObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.HeadObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(nil, nil).
		Once()

	suite.Error(suite.storage.Create("abc", strings.NewReader("content")))
}

func (suite *S3StorageTestSuite) TestCreateAPIErrInHead() {
	suite.client.
		On("HeadObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.HeadObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(
			nil,
			&smithy.GenericAPIError{
				Code:    "SomethingElse",
				Message: "test message",
			},
		).
		Once()

	suite.Error(suite.storage.Create("abc", strings.NewReader("content")))
}

func (suite *S3StorageTestSuite) TestCreateErrInHead() {
	suite.client.
		On("HeadObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.HeadObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(nil, fmt.Errorf("some error")).
		Once()

	suite.Error(suite.storage.Create("abc", strings.NewReader("content")))
}

func (suite *S3StorageTestSuite) TestRead() {
	content := "content"

	suite.client.
		On("GetObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.GetObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(
			&s3.GetObjectOutput{
				Body: readCloserWrapper{strings.NewReader(content)},
			},
			nil,
		).
		Once()

	rd, err := suite.storage.Read("abc")
	suite.NoError(err)

	data, err := io.ReadAll(rd)
	suite.NoError(err)
	suite.Equal(content, string(data))
}

func (suite *S3StorageTestSuite) TestReadErrInGet() {
	suite.client.
		On("GetObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.GetObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(
			nil,
			&smithy.GenericAPIError{
				Code:    "SomethingElse",
				Message: "test message",
			},
		).
		Once()

	rd, err := suite.storage.Read("abc")
	suite.Error(err)
	suite.Nil(rd)
}

func (suite *S3StorageTestSuite) TestUpdate() {
	suite.client.
		On("HeadObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.HeadObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(nil, nil).
		Once()

	suite.client.
		On("PutObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.PutObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(nil, nil).
		Once()

	suite.NoError(suite.storage.Update("abc", strings.NewReader("content")))
}

func (suite *S3StorageTestSuite) TestUpdateErrInHead() {
	suite.client.
		On("HeadObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.HeadObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(
			nil,
			&smithy.GenericAPIError{
				Code:    "SomethingElse",
				Message: "test message",
			},
		).
		Once()

	suite.Error(suite.storage.Update("abc", strings.NewReader("content")))
}

func (suite *S3StorageTestSuite) TestUpdateErrInPut() {
	suite.client.
		On("HeadObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.HeadObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(nil, nil).
		Once()

	suite.client.
		On("PutObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.PutObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(
			nil,
			&smithy.GenericAPIError{
				Code:    "SomethingElse",
				Message: "test message",
			},
		).
		Once()

	suite.Error(suite.storage.Update("abc", strings.NewReader("content")))
}

func (suite *S3StorageTestSuite) TestDelete() {
	suite.client.
		On("DeleteObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.DeleteObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(nil, nil).
		Once()

	suite.NoError(suite.storage.Delete("abc"))
}

func (suite *S3StorageTestSuite) TestDeleteErrInDelete() {
	suite.client.
		On("DeleteObject",
			mock.IsType(context.Background()),
			mock.MatchedBy(func(i *s3.DeleteObjectInput) bool {
				return suite.Equal("abc", *i.Key) &&
					suite.Equal(suite.bucket, *i.Bucket)
			}),
		).
		Return(
			nil,
			&smithy.GenericAPIError{
				Code:    "SomethingElse",
				Message: "test message",
			},
		).
		Once()

	suite.Error(suite.storage.Delete("abc"))
}
