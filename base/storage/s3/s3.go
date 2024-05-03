// Copyright 2024 The Mothership Authors
// SPDX-License-Identifier: Apache-2.0

package storage_s3

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/openela/mothership/base/awsutils"
	"github.com/openela/mothership/base/storage"
	"net/url"
	"os"
	"strings"
)

// S3 is an implementation of the Storage interface for S3.
type S3 struct {
	storage.Storage

	bucket     string
	uploader   *s3manager.Uploader
	downloader *s3manager.Downloader
}

// New creates a new S3 storage backend.
// Supports AWS CLI related environment variables.
func New(bucket string) (*S3, error) {
	awsCfg := &aws.Config{}
	awsutils.FillOutConfig(awsCfg)

	sess, err := session.NewSession(awsCfg)
	if err != nil {
		return nil, err
	}

	uploader := s3manager.NewUploader(sess)
	downloader := s3manager.NewDownloader(sess)

	return &S3{
		bucket:     bucket,
		uploader:   uploader,
		downloader: downloader,
	}, nil
}

// Download downloads a file from the storage backend to the given path.
func (s *S3) Download(object string, toPath string) error {
	f, err := os.OpenFile(toPath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = s.downloader.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(object),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == s3.ErrCodeNoSuchKey {
				return storage.ErrNotFound
			}
		}
	}

	return err
}

// Get returns the contents of a file from the storage backend.
func (s *S3) Get(object string) ([]byte, error) {
	buf := aws.NewWriteAtBuffer([]byte{})

	_, err := s.downloader.Download(buf, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(object),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == s3.ErrCodeNoSuchKey {
				return nil, storage.ErrNotFound
			}
		}
	}

	return buf.Bytes(), err
}

// Put uploads a file to the storage backend.
func (s *S3) Put(object string, fromPath string) (*storage.UploadInfo, error) {
	f, err := os.Open(fromPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	result, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(object),
		Body:   f,
	})

	return &storage.UploadInfo{
		Location:  result.Location,
		VersionID: result.VersionID,
	}, err
}

// PutBytes uploads a file to the storage backend.
func (s *S3) PutBytes(object string, data []byte) (*storage.UploadInfo, error) {
	result, err := s.uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(object),
		Body:   aws.ReadSeekCloser(bytes.NewBuffer(data)),
	})

	return &storage.UploadInfo{
		Location:  result.Location,
		VersionID: result.VersionID,
	}, err
}

// Delete deletes a file from the storage backend.
func (s *S3) Delete(object string) error {
	_, err := s.uploader.S3.DeleteObject(&s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(object),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == s3.ErrCodeNoSuchKey {
				return storage.ErrNotFound
			}
		}
	}
	return err
}

// Exists checks if a file exists in the storage backend.
func (s *S3) Exists(object string) (bool, error) {
	_, err := s.uploader.S3.HeadObject(&s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(object),
	})
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if awsErr.Code() == s3.ErrCodeNoSuchKey {
				return false, nil
			}
		}
		if strings.Contains(err.Error(), "NotFound") {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CanReadURI checks if a URI can be read by the storage backend.
func (s *S3) CanReadURI(uri string) (bool, error) {
	parsed, err := url.Parse(uri)
	if err != nil {
		return false, err
	}

	if parsed.Scheme != "s3" {
		return false, nil
	}

	// Verify that the host matches the bucket.
	if parsed.Host != s.bucket {
		return false, nil
	}

	return true, nil
}
