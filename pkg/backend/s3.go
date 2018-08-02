package backend

import (
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/logicmonitor/k8s-release-manager/pkg/config"
)

// S3 implement the Backend interface
type S3 struct {
	BackendConfig *config.BackendConfig
	Opts          *S3Opts
	svc           *s3.S3
}

// S3Opts represents the S3 backend configuration options
type S3Opts struct {
	Auth   *S3Auth
	Bucket string
	Region string
}

// S3Auth represents the S3 backend authentication configuration options
type S3Auth struct {
	AccessKeyID     string
	SecretAccessKey string
	SessionToken    string
}

// Read reads the specified file from the backend
func (b *S3) Read(path string) ([]byte, error) {
	buf := aws.NewWriteAtBuffer([]byte{})
	downloader := s3manager.NewDownloaderWithClient(b.client())
	_, err := downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(b.Opts.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, checkError(err)
	}
	return buf.Bytes(), nil
}

// Config returns the backend's config
func (b *S3) Config() *config.BackendConfig {
	return b.BackendConfig
}

// Writes the contents to the specified path on the backend
func (b *S3) Write(path string, data io.Reader) error {
	uploader := s3manager.NewUploaderWithClient(b.client())
	_, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(b.Opts.Bucket),
		Key:    aws.String(path),
		Body:   data,
	})
	if err != nil {
		return checkError(err)
	}
	return nil
}

// Delete deletes the specified file from the backend
func (b *S3) Delete(path string) error {
	_, err := b.client().DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(b.Opts.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return checkError(err)
	}
	return nil
}

// List lists all files in the specified path on the backend
func (b *S3) List(path string) (ret []string, err error) {
	// if the storage path is /, cleanup path
	if path == b.PathSeparator() {
		path = ""
	}

	result, err := b.client().ListObjects(&s3.ListObjectsInput{
		Bucket: aws.String(b.Opts.Bucket),
		Prefix: aws.String(path),
	})
	if err != nil {
		return nil, checkError(err)
	}

	for _, obj := range result.Contents {
		// trim the leading path from the filename
		ret = append(ret, strings.Replace(*obj.Key, path+b.PathSeparator(), "", 1))
	}
	return ret, err
}

// PathSeparator returns the backend-specific path separator
func (b *S3) PathSeparator() string {
	return "/"
}

func (b *S3) client() *s3.S3 {
	if b.svc == nil {
		sess := session.Must(session.NewSession(&aws.Config{
			Region:      aws.String(b.Opts.Region),
			Credentials: b.getCreds(),
		}))

		svc := s3.New(sess, &aws.Config{
			Region: aws.String(b.Opts.Region),
		})
		b.svc = svc
	}
	return b.svc
}

func (b *S3) getCreds() *credentials.Credentials {
	if b.Opts.Auth.AccessKeyID == "" || b.Opts.Auth.SecretAccessKey == "" {
		return nil
	}

	if b.Opts.Auth.SessionToken != "" {
		return credentials.NewStaticCredentialsFromCreds(credentials.Value{
			AccessKeyID:     b.Opts.Auth.AccessKeyID,
			SecretAccessKey: b.Opts.Auth.SecretAccessKey,
			SessionToken:    b.Opts.Auth.SessionToken,
		})
	}
	return credentials.NewStaticCredentialsFromCreds(credentials.Value{
		AccessKeyID:     b.Opts.Auth.AccessKeyID,
		SecretAccessKey: b.Opts.Auth.SecretAccessKey,
	})
}

func checkError(err error) error {
	if aerr, ok := err.(awserr.Error); ok {
		switch aerr.Code() {
		case s3.ErrCodeNoSuchBucket:
			return fmt.Errorf("%s %s", s3.ErrCodeNoSuchBucket, aerr.Error())
		default:
			return fmt.Errorf(aerr.Error())
		}
	} else {
		// Print the error, cast err to awserr.Error to get the Code and
		// Message from an error.
		return fmt.Errorf(err.Error())
	}
}
