package io

import (
	"context"
	"fmt"
	"io"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

// Implements the io.ReaderAt interface, reading specific parts of a single S3 file.
// Since we need to perform API calls on-demand depending on the code that uses this
// reader, the API calls for S3 have to be performed from here and not the repository.
// This class should only be used through the S3Repository class.
type S3FileReaderAt struct {
	service *awsS3.Client
	bucket  string
	path    string

	// The context to use for the API calls
	// since the S3FileReaderAt will implement with io.Reader
	// but method ReadAt has not context parameter
	ctx context.Context
}

func NewS3FileReaderAt(
	ctx context.Context,
	service *awsS3.Client,
	bucket string,
	path string,
) *S3FileReaderAt {
	return &S3FileReaderAt{
		service,
		bucket,
		path,
		ctx,
	}
}

func (reader *S3FileReaderAt) ReadAt(
	outputBytes []byte,
	offsetFrom int64,
) (n int, err error) {
	expectedLength := len(outputBytes)
	offsetTo := offsetFrom + int64(expectedLength)

	// Fetching the required bytes range from S3 in a single call
	// It is considered to be the responsibility of the code that
	// calls ReadAt to handle any required chunk splitting,
	// so this methods always performs a single API call.
	object, err := reader.service.GetObject(reader.ctx, &awsS3.GetObjectInput{
		Bucket: aws.String(reader.bucket),
		Key:    aws.String(reader.path),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", offsetFrom, offsetTo-1)),
	})
	if err != nil {
		return 0, err
	}

	// Now that we got the required response, we copy the returned
	// bytes to the buffer we got as an input.
	bytesCount, err := io.ReadFull(object.Body, outputBytes)
	if err == io.ErrUnexpectedEOF {
		// We must not expose underlying implementation errors,
		// and ErrUnexpectedEOF is not expected by io.ReaderAt.ReadAt
		return bytesCount, io.EOF
	}

	return bytesCount, err
}
