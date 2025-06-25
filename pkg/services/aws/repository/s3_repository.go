package repository

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"math"
	"mime"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"g-management/pkg/log"
	"g-management/pkg/services/aws/entity"
	customIo "g-management/pkg/services/aws/io"
	customLogger "g-management/pkg/services/aws/log"
	"g-management/pkg/shared/utils"
	"g-management/pkg/shared/wraperror"

	aws "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	awsS3Manager "github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
)

const (
	PublicPresignedUrlExpiryTime  = 7 * 24 * time.Hour
	PrivatePresignedUrlExpiryTime = 7 * 24 * time.Hour
)

const (
	MaxS3UploadPartsCount  = 10_000            // Hard limit of the S3 API
	DefaultS3PartSize      = int64(30_000_000) // Defaulting to 30Mb
	ProgressReportInterval = time.Second * 5
	prefixCDN              = "cdn."

	defaultPartSize   = 5 * 1024 * 1024 // 5 MB per download.
	defaultConcurrent = 2               // number of concurrent download s3 part.

	MaxDurationWaiter time.Duration = 30 * time.Second
	codeNotFound                    = "NotFound"
)

type CustomFile struct {
	File        *os.File
	ContentType *string
}

type FileInfo struct {
	Bucket   string
	Path     string
	MimeType string
	Stream   io.ReaderAt
	FileSize int64
}

type S3RepositoryInterface interface {
	BucketExists(ctx context.Context, bucket string) (bool, error)
	PublicUrl(ctx context.Context, bucket string, path string) (string, error)
	UploadUrl(ctx context.Context, bucket string, path string) (*entity.FileUploadInformation, error)
	NormalizePath(path string) string
	NormalizeDirPath(path string) string
	ListFiles(ctx context.Context, bucket string, path string) (map[string]types.Object, error)
	PutMultipartFromStream(
		ctx context.Context,
		logger *log.Logger,
		fileInfo FileInfo,
		onProgress func(processedBytes int64),
	) error
	GetFileReaderAt(
		ctx context.Context,
		bucket string,
		path string,
	) io.ReaderAt
	GetFileReader(
		ctx context.Context,
		bucket string,
		path string,
	) (
		reader io.ReadCloser,
		fileSize int64,
		mimeType *string,
		err error,
	)
	RenameImagePath(ctx context.Context, bucket, srcPath, desPath string) error
	MultipartUploadUrl(ctx context.Context, bucket, path, filename, mimeType string, fileSize int64) (*entity.MultipartUploadInformation, *string, error)
	CompleteMultipartUpload(ctx context.Context, bucket, path, uploadID string, fileSize int64) error
	GenerateMultipartUploadUrl(ctx context.Context, bucket, path, uploadID string, partNumber, partSize, fileSize int64) ([]entity.PartUploadInformation, error)
	AbortMultipartUpload(ctx context.Context, bucket, path, uploadID string) error
	DeleteAllFilesInPath(ctx context.Context, bucket, path string) error
	DeleteFile(ctx context.Context, bucket, path string) error
	ObjectExists(ctx context.Context, bucket, path string) (bool, error)
	DeleteUrl(ctx context.Context, bucket, path string) (*entity.FileUploadInformation, error)
	UploadFileToInternal(ctx context.Context, filename, path string, body []byte, contentType string) error
	NonExistPublicUrl(ctx context.Context, bucket string, path string) (string, error)
	GetHeadObject(ctx context.Context, bucket, path string) (*awsS3.HeadObjectOutput, error)
	GenerateCloudFrontUrlWithTimeStamp(ctx context.Context, bucket string, domain, path *string) (*string, error)
	DeleteMultipleFile(ctx context.Context, bucket string, imagePaths []string) error
	GenerateCloudFrontUrlWithCDN(ctx context.Context, bucket string, domain, path *string) (*string, error)
	GetContentTypeAndUrlWithCDN(ctx context.Context, bucket string, domain, path *string) (*string, *string, error)
	GetContentTypeAndUrl(ctx context.Context, bucket string, domain, path *string) (*string, *string, error)
	DownloadFile(ctx context.Context, bucket string, key string) (*os.File, error)
	UploadFile(ctx context.Context, bucket string, key string, file *os.File) error
	DeleteFilesInPath(ctx context.Context, bucket string, path string) error
	GenerateCloudFrontUrlForAudio(domain, path *string) (*string, error)
	UploadFileMultipart(ctx context.Context, bucket string, key string, file *os.File) error
	UploadFileWithContext(ctx context.Context, bucket string, key string, customFile CustomFile) error
	ObjectExistsAndLastModified(ctx context.Context, bucket, path string) (bool, *time.Time, error)
	GenerateCloudFrontUrlCDNWithoutTime(domain, path *string) (*string, error)
	UploadFileStream(ctx context.Context, bucket string, key string, streamData *io.PipeReader, contentType string, filename string) error
	DownloadFileToLocal(ctx context.Context, bucket string, key string, outputFile string) (string, error)
}

type S3Repository struct {
	client                    *awsS3.Client
	clientWithOnlyErrorLogger *awsS3.Client
	presignClient             *awsS3.PresignClient
	s3Downloader              *awsS3Manager.Downloader
	s3Uploader                *awsS3Manager.Uploader
}

// newS3Client comment
// en: create a new s3 client
func newS3Client(cfg aws.Config, s3Ops awsS3.Options, logger customLogger.AwsLogger, clientLogMode aws.ClientLogMode) *awsS3.Client {
	client := awsS3.NewFromConfig(cfg, func(options *awsS3.Options) {
		options.BaseEndpoint = s3Ops.BaseEndpoint
		options.UsePathStyle = s3Ops.UsePathStyle
		options.EndpointOptions.DisableHTTPS = s3Ops.EndpointOptions.DisableHTTPS
		options.Logger = logger
		options.ClientLogMode = clientLogMode
	})
	return client
}

func NewS3Repository() (S3RepositoryInterface, error) {
	var s3Ops awsS3.Options
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
	)
	s3Ops = awsS3.Options{
		BaseEndpoint: aws.String(os.Getenv("ENDPOINT")),
	}
	s3Ops.EndpointOptions.DisableHTTPS = os.Getenv("S3_DISABLE_SSL") == utils.TrueStatusString

	if err != nil {
		return nil, err
	}

	client := newS3Client(
		cfg,
		s3Ops,
		customLogger.AwsLogger{Logger: log.With("service", "aws-s3-api")},
		aws.LogRequest|aws.LogResponseWithBody|aws.LogRetries,
	)

	clientWithOnlyErrorLogger := newS3Client(
		cfg,
		s3Ops,
		customLogger.AwsLogger{Logger: log.With("service", "aws-s3-api-only-error-logger")},
		aws.LogRequest|aws.LogResponse|aws.LogRetries,
	)

	presignClient := awsS3.NewPresignClient(newS3Client(
		cfg,
		s3Ops,
		customLogger.AwsLogger{Logger: log.With("service", "aws-s3-api-presign")},
		aws.LogRequest|aws.LogResponseWithBody|aws.LogRetries,
	))

	return &S3Repository{
		client:                    client,
		clientWithOnlyErrorLogger: clientWithOnlyErrorLogger,
		presignClient:             presignClient,
		s3Downloader: awsS3Manager.NewDownloader(clientWithOnlyErrorLogger, func(downloader *awsS3Manager.Downloader) {
			downloader.PartSize = defaultPartSize
			downloader.Concurrency = defaultConcurrent
		}),
		s3Uploader: awsS3Manager.NewUploader(clientWithOnlyErrorLogger, func(uploader *awsS3Manager.Uploader) {
			uploader.PartSize = defaultPartSize
			uploader.Concurrency = defaultConcurrent
		}),
	}, nil
}

func (fur *S3Repository) BucketExists(ctx context.Context, bucket string) (bool, error) {
	_, err := fur.client.HeadBucket(ctx, &awsS3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return false, nil
		}

		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			code := apiErr.ErrorCode()
			if code == codeNotFound {
				return false, nil
			}
			return false, err
		}

		return false, err
	}

	return true, nil
}

func (fur *S3Repository) ObjectExists(ctx context.Context, bucket, path string) (bool, error) {
	_, err := fur.client.HeadObject(ctx, &awsS3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return false, nil
		}

		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			code := apiErr.ErrorCode()
			if code == codeNotFound {
				return false, nil
			}
			return false, err
		}

		return false, err
	}

	return true, nil
}

// ObjectExistsAndLastModified comment
// en: ObjectExistsAndLastModified func data to check file exist in S3 and return time last modified of file
func (fur *S3Repository) ObjectExistsAndLastModified(ctx context.Context, bucket, path string) (bool, *time.Time, error) {
	data, err := fur.client.HeadObject(ctx, &awsS3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return false, nil, nil
		}

		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			code := apiErr.ErrorCode()
			if code == codeNotFound {
				return false, nil, nil
			}
			return false, nil, err
		}

		return false, nil, err
	}
	lastModified := data.LastModified
	return true, lastModified, nil
}

func (fur *S3Repository) PublicUrl(ctx context.Context, bucket, path string) (string, error) {
	isExist, err := fur.ObjectExists(ctx, bucket, path)
	if err != nil {
		return "", err
	}
	if !isExist {
		return "", utils.ErrS3ObjectNotFound
	}

	req, err := fur.presignClient.PresignGetObject(ctx, &awsS3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	}, func(po *awsS3.PresignOptions) {
		po.Expires = PublicPresignedUrlExpiryTime
	})

	return req.URL, err
}

func (fur *S3Repository) UploadUrl(ctx context.Context, bucket, path string) (*entity.FileUploadInformation, error) {
	req, err := fur.presignClient.PresignPutObject(ctx, &awsS3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	}, func(po *awsS3.PresignOptions) {
		po.Expires = PrivatePresignedUrlExpiryTime
	})

	return &entity.FileUploadInformation{
		Method: req.Method,
		Url:    req.URL,
	}, err
}

// Normalizes the path (removing any leading slash)
func (fur *S3Repository) NormalizePath(path string) string {
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	return path
}

// Normalizes the path with a slash at the end, but no leading slash.
// S3 directories needs a slash at the end, otherwise it can match
// partial file names, since S3 only have keys (not directories)
func (fur *S3Repository) NormalizeDirPath(path string) string {
	path = fur.NormalizePath(path)
	if len(path) > 0 && path[len(path)-1] != '/' {
		path += "/"
	}

	return path
}

// Returns a list of files in the given S3 directory
// Returned keys are a path relative to that directory
func (fur *S3Repository) ListFiles(ctx context.Context, bucket, path string) (map[string]types.Object, error) {
	path = fur.NormalizeDirPath(path)

	response, err := fur.client.ListObjects(ctx, &awsS3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	list := map[string]types.Object{}
	if response.Contents != nil {
		for _, object := range response.Contents {
			if object.Key == nil {
				return nil, fmt.Errorf("unexpected nil file key: %+v", object)
			}

			fileName := (*object.Key)[len(path):]
			list[fileName] = object
		}
	}

	return list, nil
}

// Lists the list of already-uploaded parts for a specific multipart upload in S3
func (fur *S3Repository) retrieveMultipartUploadStatus(
	ctx context.Context,
	bucket string,
	path string,
	uploadId string,
) (
	lastUploadedPartNumber int64,
	alreadyUploadedBytes int64,
	completedParts []types.CompletedPart,
	err error,
) {
	lastUploadedPartNumber = 0
	alreadyUploadedBytes = 0
	completedParts = make([]types.CompletedPart, 0)

	// Listing parts in S3 requires handling the paging of
	// the S3 API, via the given callback below
	var partCallbackErr error
	paginator := awsS3.NewListPartsPaginator(fur.client, &awsS3.ListPartsInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(path),
		UploadId: aws.String(uploadId),
	})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return 0, 0, nil, err
		}

		for _, part := range output.Parts {
			if part.PartNumber == nil || part.Size == nil {
				partCallbackErr = fmt.Errorf("the upload part contained nil information and cannot be processed: %+v", part)
				break
			}

			alreadyUploadedBytes += *part.Size
			if int64(*part.PartNumber) > lastUploadedPartNumber {
				lastUploadedPartNumber = int64(*part.PartNumber)
			}

			completedParts = append(completedParts, types.CompletedPart{
				ETag:       part.ETag,
				PartNumber: part.PartNumber,
			})
		}
	}
	if partCallbackErr != nil {
		return 0, 0, nil, partCallbackErr
	}

	return lastUploadedPartNumber, alreadyUploadedBytes, completedParts, nil
}

// Determines the part size to use for the next uploads.
// We default to a part size of 10Mb, but since we need
// to finish the file with a maximum of 10k parts, this
// size cannot be constant and is determined by this function.
// For example, for a 200Gb file, we would need an average part
// size of 200Gb/10k = 20Mb.
// This also handles edge cases, for example if we have only one
// remaining part but 1Gb of data to transfer, we don't have any
// other choice than handling a single part of 1Gb.
func (fur *S3Repository) getMultipartUploadPartSize(
	lastUploadedPartNumber int64,
	alreadyUploadedBytes int64,
	fileSize int64,
) int64 {
	remainingPossibleParts := MaxS3UploadPartsCount - lastUploadedPartNumber
	remainingBytes := fileSize - alreadyUploadedBytes
	partSize := DefaultS3PartSize
	if remainingBytes > partSize*remainingPossibleParts {
		partSize = int64(
			math.Ceil(
				float64(remainingBytes) / float64(remainingPossibleParts),
			),
		)
	}

	return partSize
}

// Creates a file in S3 using a multipart transfer from a given stream.
// If a previous partial transfer already exists in S3, it will retrieve
// the status and continue the transfer from where is stopped.
// This creates a MultipartUpload object in S3, and does not delete it
// unless the transfer is completely finished.
// Partial unfinished/abandoned transfers can take space and be billed by S3.
// This function is not thread-safe
func (fur *S3Repository) PutMultipartFromStream(
	ctx context.Context,
	logger *log.Logger,
	fileInfo FileInfo,
	onProgress func(processedBytes int64),
) error {
	path := fur.NormalizePath(fileInfo.Path)

	uploadsList, err := fur.client.ListMultipartUploads(ctx, &awsS3.ListMultipartUploadsInput{
		Bucket: aws.String(fileInfo.Bucket),
		Prefix: aws.String(path),
	})
	if err != nil {
		return err
	}

	// Retrieving or initializing the upload status depending on the S3 information.
	var uploadId string
	var lastUploadedPartNumber int64
	var alreadyUploadedBytes int64
	var completedParts []types.CompletedPart
	if len(uploadsList.Uploads) > 0 && uploadsList.Uploads[0].UploadId != nil {
		// A MultipartUpload object already exists in S3 for this file.
		// In this case, we use it and continue the transfer using this information.
		// This assumes no other process is trying to upload the same file at the same time.
		// Doing so would cause a corrupted file.

		uploadId = *uploadsList.Uploads[0].UploadId
		logger = logger.With("uploadId", uploadId)
		logger.Info(ctx, fmt.Sprintf("Found %v upload(s) in S3. Using the first one.\n", len(uploadsList.Uploads)))

		lastUploadedPartNumber, alreadyUploadedBytes, completedParts, err = fur.retrieveMultipartUploadStatus(ctx, fileInfo.Bucket, path, uploadId)
		if err != nil {
			return err
		}

		logger.Info(ctx, fmt.Sprintf("Found %v parts already uploaded (%v bytes)\n", lastUploadedPartNumber, alreadyUploadedBytes))
	} else {
		// No existing pending uploads: creating an empty one and starting the transfer from zero.

		createdUpload, err := fur.client.CreateMultipartUpload(ctx, &awsS3.CreateMultipartUploadInput{
			Bucket:      aws.String(fileInfo.Bucket),
			ContentType: aws.String(fileInfo.MimeType),
			Key:         aws.String(path),
		})
		if err != nil {
			return err
		}

		uploadId = *createdUpload.UploadId
		lastUploadedPartNumber = 0
		alreadyUploadedBytes = 0
		completedParts = make([]types.CompletedPart, 0)

		logger = logger.With("uploadId", uploadId)
		logger.Info(ctx, "Created a multipart upload in S3")
	}

	// Getting the part size we have to use
	partSize := fur.getMultipartUploadPartSize(lastUploadedPartNumber, alreadyUploadedBytes, fileInfo.FileSize)

	nextProgressReport := time.Now().Add(ProgressReportInterval)
	for alreadyUploadedBytes < fileInfo.FileSize {
		partNumber := lastUploadedPartNumber + 1
		partReader := io.NewSectionReader(fileInfo.Stream, alreadyUploadedBytes, partSize)

		progressPercent := (float64(alreadyUploadedBytes) / float64(fileInfo.FileSize)) * float64(100)
		logger.
			With("fileSize", fileInfo.FileSize).
			With("currentOffset", alreadyUploadedBytes).
			With("progress", progressPercent).
			With("partNumber", partNumber).
			With("partSize", partSize).
			Info(ctx, "Uploading part ...\n")
		if now := time.Now(); nextProgressReport.Before(now) {
			// Setting a minimum interval between progress reports
			// to avoid slowing down the process too much
			onProgress(alreadyUploadedBytes)
			nextProgressReport = now.Add(ProgressReportInterval)
		}

		// The requested part size will be different when we reach the end of the file
		contentLength := partSize
		if alreadyUploadedBytes+partSize > fileInfo.FileSize {
			contentLength = fileInfo.FileSize - alreadyUploadedBytes
		}

		uploadedPart, err := fur.clientWithOnlyErrorLogger.UploadPart(ctx, &awsS3.UploadPartInput{
			Body:          partReader,
			Bucket:        aws.String(fileInfo.Bucket),
			Key:           aws.String(path),
			UploadId:      aws.String(uploadId),
			PartNumber:    aws.Int32(int32(partNumber)),
			ContentLength: aws.Int64(contentLength),
		})
		if err != nil {
			return err
		}

		completedParts = append(completedParts, types.CompletedPart{
			ETag:       uploadedPart.ETag,
			PartNumber: aws.Int32(int32(partNumber)),
		})

		lastUploadedPartNumber = partNumber
		alreadyUploadedBytes += partSize
	}

	_, err = fur.client.CompleteMultipartUpload(ctx, &awsS3.CompleteMultipartUploadInput{
		Bucket:   aws.String(fileInfo.Bucket),
		Key:      aws.String(path),
		UploadId: aws.String(uploadId),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})

	return err
}

// Allows reading a file in chunks via a standard io interface.
// Subsequent API calls are performed by pkg/io/s3_reader_at.go
// (so are the error cases)
func (fur *S3Repository) GetFileReaderAt(
	ctx context.Context,
	bucket string,
	path string,
) io.ReaderAt {
	return customIo.NewS3FileReaderAt(
		ctx,
		fur.clientWithOnlyErrorLogger,
		bucket,
		fur.NormalizePath(path),
	)
}

func (fur *S3Repository) GetFileReader(
	ctx context.Context,
	bucket string,
	path string,
) (
	reader io.ReadCloser,
	fileSize int64,
	mimeType *string,
	err error,
) {
	object, err := fur.client.GetObject(ctx, &awsS3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, 0, nil, err
	}
	if object.ContentLength == nil {
		return nil, 0, nil, fmt.Errorf("unexpected nil Content-Length: %+v", object)
	}

	return object.Body, *object.ContentLength, object.ContentType, nil
}

func (fur *S3Repository) RenameImagePath(ctx context.Context, bucket, srcPath, desPath string) error {
	source := bucket + "/" + srcPath
	_, err := fur.client.CopyObject(ctx, &awsS3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(url.PathEscape(source)),
		Key:        aws.String(desPath),
	})
	if err != nil {
		return err
	}

	objectExistsWaiter := awsS3.NewObjectExistsWaiter(fur.client)
	// Wait to see if the item got copied
	err = objectExistsWaiter.Wait(ctx, &awsS3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(desPath)}, MaxDurationWaiter)
	if err != nil {
		return err
	}

	// delete old object
	_, err = fur.client.DeleteObject(ctx, &awsS3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(srcPath),
	})
	if err != nil {
		return err
	}

	objectNotExistsWaiter := awsS3.NewObjectNotExistsWaiter(fur.client)
	err = objectNotExistsWaiter.Wait(ctx, &awsS3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(srcPath)}, MaxDurationWaiter)
	return err
}

// Generate presign multipart url
func (fur *S3Repository) MultipartUploadUrl(ctx context.Context, bucket, path, filename, mimeType string, fileSize int64) (*entity.MultipartUploadInformation, *string, error) {
	existed, err := fur.BucketExists(ctx, bucket)
	if err != nil {
		return nil, nil, err
	}
	if !existed {
		return nil, nil, wraperror.NewApiDisplayableError(
			http.StatusInternalServerError,
			"bucket did not existed",
			errors.New("bucket "+bucket+" dis not existed"),
		)
	}

	contentDisposition := fmt.Sprintf(`filename="%v"`, filename)
	createdUpload, err := fur.client.CreateMultipartUpload(ctx, &awsS3.CreateMultipartUploadInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(path),
		ContentType:        aws.String(mimeType),
		ContentDisposition: aws.String(contentDisposition),
	})
	if err != nil {
		return nil, nil, err
	}

	uploadId := *createdUpload.UploadId

	partNumber, partSize := fur.computeMultipartUploadPart(fileSize)
	parts, err := fur.GenerateMultipartUploadUrl(ctx, bucket, path, uploadId, partNumber, partSize, fileSize)
	if err != nil {
		return nil, nil, err
	}

	return &entity.MultipartUploadInformation{
		Method: "PUT",
		Parts:  parts,
	}, &uploadId, nil
}

// compute how many parts if each part have 10MB,
// but if parts number greater than max upload part
// then reduce size of part
func (fur *S3Repository) computeMultipartUploadPart(fileSize int64) (int64, int64) {
	var partNumber int64
	var partSize int64
	partSize = DefaultS3PartSize
	partNumber = int64(
		math.Ceil(
			float64(fileSize) / float64(partSize),
		),
	)

	if partNumber > MaxS3UploadPartsCount {
		partSize = int64(
			math.Ceil(
				float64(fileSize) / float64(MaxS3UploadPartsCount),
			),
		)

		partNumber = MaxS3UploadPartsCount
	}

	return partNumber, partSize
}

func (fur *S3Repository) CompleteMultipartUpload(ctx context.Context, bucket, path, uploadID string, fileSize int64) error {
	_, alreadyUploadedBytes, completedParts, err := fur.retrieveMultipartUploadStatus(ctx, bucket, path, uploadID)
	if err != nil {
		return err
	}

	if alreadyUploadedBytes < fileSize {
		return wraperror.NewApiDisplayableError(
			http.StatusBadRequest,
			"All parts did not finished uploading yet",
			errors.New("all parts did not finished uploading yet"),
		)
	}

	_, err = fur.client.CompleteMultipartUpload(ctx, &awsS3.CompleteMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(path),
		UploadId: aws.String(uploadID),
		MultipartUpload: &types.CompletedMultipartUpload{
			Parts: completedParts,
		},
	})

	return err
}

func (fur *S3Repository) GenerateMultipartUploadUrl(ctx context.Context, bucket, path, uploadID string, partNumber, partSize, fileSize int64) ([]entity.PartUploadInformation, error) {
	parts := make([]entity.PartUploadInformation, 0, partNumber)
	for index := int64(1); index <= partNumber; index++ {
		uploadPartReq, err := fur.presignClient.PresignUploadPart(ctx, &awsS3.UploadPartInput{
			Bucket:     aws.String(bucket),
			Key:        aws.String(path),
			UploadId:   aws.String(uploadID),
			PartNumber: aws.Int32(int32(index)),
		}, func(po *awsS3.PresignOptions) {
			po.Expires = PrivatePresignedUrlExpiryTime
		})
		if err != nil {
			return nil, err
		}

		offsetFrom := (index - 1) * partSize
		var offsetTo int64
		if index == partNumber {
			offsetTo = fileSize - 1
		} else {
			offsetTo = (index * partSize) - 1
		}

		parts = append(parts, entity.PartUploadInformation{
			OffsetFrom: offsetFrom,
			OffsetTo:   offsetTo,
			Url:        uploadPartReq.URL,
		})
	}

	return parts, nil
}

func (fur *S3Repository) AbortMultipartUpload(ctx context.Context, bucket, path, uploadID string) error {
	_, _, _, err := fur.retrieveMultipartUploadStatus(ctx, bucket, path, uploadID)
	if err != nil {
		return err
	}

	_, err = fur.client.AbortMultipartUpload(ctx, &awsS3.AbortMultipartUploadInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(path),
		UploadId: aws.String(uploadID),
	})

	return err
}

func (fur *S3Repository) DeleteAllFilesInPath(ctx context.Context, bucket, path string) error {
	existed, err := fur.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !existed {
		return wraperror.NewApiDisplayableError(
			http.StatusInternalServerError,
			"bucket did not existed",
			errors.New("bucket "+bucket+" dis not existed"),
		)
	}

	listFileInPath, err := fur.ListFiles(ctx, bucket, path)
	if err != nil {
		return err
	}

	for filename := range listFileInPath {
		srcPath := path + "/" + filename
		_, err = fur.client.DeleteObject(ctx, &awsS3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(srcPath),
		})
		if err != nil {
			return err
		}

		objectNotExistsWaiter := awsS3.NewObjectNotExistsWaiter(fur.client)
		err = objectNotExistsWaiter.Wait(ctx, &awsS3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(srcPath)}, MaxDurationWaiter)
		if err != nil {
			return err
		}
	}
	return nil
}

func (fur *S3Repository) DeleteFile(ctx context.Context, bucket, path string) error {
	path = fur.NormalizePath(path)
	existed, err := fur.BucketExists(ctx, bucket)
	if err != nil {
		return err
	}
	if !existed {
		return errors.New(utils.ErrInvalidS3Bucket)
	}

	_, err = fur.client.DeleteObject(ctx, &awsS3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return err
	}

	objectNotExistsWaiter := awsS3.NewObjectNotExistsWaiter(fur.client)
	return objectNotExistsWaiter.Wait(ctx, &awsS3.HeadObjectInput{Bucket: aws.String(bucket), Key: aws.String(path)}, MaxDurationWaiter)
}

func (fur *S3Repository) DeleteUrl(ctx context.Context, bucket, path string) (*entity.FileUploadInformation, error) {
	req, err := fur.presignClient.PresignDeleteObject(ctx, &awsS3.DeleteObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	}, func(po *awsS3.PresignOptions) {
		po.Expires = PrivatePresignedUrlExpiryTime
	})
	if err != nil {
		return nil, err
	}

	return &entity.FileUploadInformation{
		Method: req.Method,
		Url:    req.URL,
	}, nil
}

func (fur *S3Repository) UploadFileToInternal(ctx context.Context, filename, path string, body []byte, contentType string) error {
	_, err := fur.clientWithOnlyErrorLogger.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket:             aws.String(os.Getenv("S3_INTERNAL_BUCKET")),
		Key:                aws.String(path),
		Body:               bytes.NewReader(body),
		ContentType:        aws.String(contentType),
		ContentDisposition: aws.String("attachment;filename=" + filename),
	})
	return err
}

func (fur *S3Repository) NonExistPublicUrl(ctx context.Context, bucket, path string) (string, error) {
	req, err := fur.presignClient.PresignGetObject(ctx, &awsS3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	}, func(po *awsS3.PresignOptions) {
		po.Expires = PublicPresignedUrlExpiryTime
	})

	return req.URL, err
}

func (fur *S3Repository) GetHeadObject(ctx context.Context, bucket, path string) (*awsS3.HeadObjectOutput, error) {
	return fur.client.HeadObject(ctx, &awsS3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(path),
	})
}

func (fur *S3Repository) GenerateCloudFrontUrlWithTimeStamp(ctx context.Context, bucket string, domain, path *string) (*string, error) {
	if domain == nil || path == nil {
		return nil, nil
	}

	urlParsed, err := url.Parse(*domain)
	if err != nil {
		return nil, err
	}

	urlPath, err := url.Parse(urlParsed.Scheme + "://" + urlParsed.Host + utils.NormalizeRelativePath(*path))
	if err != nil {
		return nil, err
	}

	lastObject, err := fur.GetHeadObject(ctx, bucket, *path)
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return nil, nil
		}

		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			code := apiErr.ErrorCode()
			if code == codeNotFound {
				return nil, nil
			}
			return nil, err
		}

		return nil, err
	}
	if lastObject == nil {
		return nil, nil
	}

	v := urlPath.Query()
	v.Add("time", fmt.Sprintf("%d", lastObject.LastModified.Unix()))
	urlPath.RawQuery = v.Encode()

	return utils.NewStringPointer(urlPath.String()), nil
}

func (fur *S3Repository) DeleteMultipleFile(ctx context.Context, bucket string, imagePaths []string) error {
	var objectIds []types.ObjectIdentifier
	for _, key := range imagePaths {
		objectIds = append(objectIds, types.ObjectIdentifier{Key: aws.String(key)})
	}
	_, err := fur.client.DeleteObjects(ctx, &awsS3.DeleteObjectsInput{
		Bucket: aws.String(bucket),
		Delete: &types.Delete{Objects: objectIds},
	})
	return err
}

func (fur *S3Repository) GenerateCloudFrontUrlWithCDN(ctx context.Context, bucket string, domain, path *string) (*string, error) {
	lastObject, err := fur.GetHeadObject(ctx, bucket, *path)
	if err != nil {
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			return nil, nil
		}

		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			code := apiErr.ErrorCode()
			if code == codeNotFound {
				return nil, nil
			}
			return nil, err
		}

		return nil, err
	}
	if lastObject == nil {
		return nil, nil
	}

	if domain == nil || path == nil {
		return nil, nil
	}
	urlParsed, err := utils.ParseUrl(*domain)
	if err != nil {
		return nil, err
	}

	urlBuilder := urlParsed.Scheme + "://"
	if urlParsed.Subdomain != "" {
		urlSubdomain := strings.ReplaceAll(urlParsed.Subdomain, ".", "-")
		urlBuilder += urlSubdomain + "." + prefixCDN
	} else {
		urlBuilder += prefixCDN
	}
	urlBuilder += urlParsed.Domain + "." + urlParsed.TLD + utils.NormalizeRelativePath(*path)

	urlPath, err := url.Parse(urlBuilder)
	if err != nil {
		return nil, err
	}
	v := urlPath.Query()
	v.Add("time", fmt.Sprintf("%d", lastObject.LastModified.Unix()))
	urlPath.RawQuery = v.Encode()

	return utils.NewStringPointer(urlPath.String()), nil
}

// GetContentTypeAndUrlWithCDN comment
// en: get content_type and url with cdn
// en: input: bucket, domain, path
// en: output: content_type, url
// GetContentTypeAndUrlWithCDN retrieves the content type and CDN URL for an object in an S3 bucket.
func (fur *S3Repository) GetContentTypeAndUrlWithCDN(ctx context.Context, bucket string, domain, path *string) (*string, *string, error) {
	// en: get data in the S3 bucket
	lastObject, err := fur.GetHeadObject(ctx, bucket, *path)
	if err != nil {
		// en: check if error is due to the object not existing
		var nsk *types.NoSuchKey
		if errors.As(err, &nsk) {
			// en: if object does not exist, return nil
			return nil, nil, nil
		}

		// en: check if the error is an API error and handle specific cases
		var apiErr smithy.APIError
		if errors.As(err, &apiErr) {
			code := apiErr.ErrorCode()
			// en: If error is not found, return nil
			if code == codeNotFound {
				return nil, nil, nil
			}
			// en: For other API errors, return error
			return nil, nil, err
		}

		// en: For any other errors, return error
		return nil, nil, err
	}

	// en: if lastObject is nil, return nil
	if lastObject == nil {
		return nil, nil, nil
	}

	// en: check if domain or path is nil, return nil
	if domain == nil || path == nil {
		return nil, nil, nil
	}

	// en: Parse the provided domain to extract its components
	urlParsed, err := utils.ParseUrl(*domain)
	if err != nil {
		// en: parsing fails, return error
		return nil, nil, err
	}

	// en: Build the CDN URL based on the parsed domain
	urlBuilder := urlParsed.Scheme + "://"
	if urlParsed.Subdomain != "" {
		// Replace dots in the subdomain with dashes for CDN compatibility
		urlSubdomain := strings.ReplaceAll(urlParsed.Subdomain, ".", "-")
		urlBuilder += urlSubdomain + "." + prefixCDN
	} else {
		// en: ff no subdomain, just use prefixCDN
		urlBuilder += prefixCDN
	}
	// en: append domain and TLD to URL.
	urlBuilder += urlParsed.Domain + "." + urlParsed.TLD + utils.NormalizeRelativePath(*path)

	// en: Parse constructed URL to ensure it's valid
	urlPath, err := url.Parse(urlBuilder)
	if err != nil {
		// en: If parsing fails, return error
		return nil, nil, err
	}

	// em: ddd the last modified time as a query parameter to URL
	v := urlPath.Query()
	v.Add("time", fmt.Sprintf("%d", lastObject.LastModified.Unix()))
	urlPath.RawQuery = v.Encode()

	// en: Return the constructed URL and the content_type of object
	return utils.NewStringPointer(urlPath.String()), lastObject.ContentType, nil
}

// DownloadFile comment.
// en: download file from s3
func (fur *S3Repository) DownloadFile(ctx context.Context, bucket string, key string) (*os.File, error) {
	file, err := os.CreateTemp("/tmp", "tmp")
	if err != nil {
		return nil, err
	}

	_, err = fur.s3Downloader.Download(ctx, file, &awsS3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return file, err
}

// UploadFile comment.
// en: upload file to s3
func (fur *S3Repository) UploadFile(ctx context.Context, bucket string, key string, file *os.File) error {
	fileName := filepath.Base(file.Name())
	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	contentDisposition := fmt.Sprintf(`attachment;filename="%v"`, fileName)
	_, err := fur.clientWithOnlyErrorLogger.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(key),
		ContentDisposition: aws.String(contentDisposition),
		Body:               file,
		ContentType:        aws.String(contentType),
	})
	return err
}

// UploadFileMultipart comment.
// en: upload file to s3
func (fur *S3Repository) UploadFileMultipart(ctx context.Context, bucket string, key string, file *os.File) error {
	fileName := filepath.Base(file.Name())
	contentType := mime.TypeByExtension(filepath.Ext(fileName))
	contentDisposition := fmt.Sprintf(`attachment;filename="%v"`, fileName)
	_, err := fur.s3Uploader.Upload(ctx, &awsS3.PutObjectInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(key),
		ContentDisposition: aws.String(contentDisposition),
		Body:               file,
		ContentType:        aws.String(contentType),
	})
	return err
}

// DeleteFilesInPath comment.
// en: delete all objects in given path.
func (fur *S3Repository) DeleteFilesInPath(ctx context.Context, bucket string, path string) error {
	path = fur.NormalizeDirPath(path)

	paginator := awsS3.NewListObjectsV2Paginator(fur.client, &awsS3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(path),
	})
	var deleteObjectErr error
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			return err
		}

		if output.Contents == nil {
			continue
		}

		var objectIds []types.ObjectIdentifier
		for _, object := range output.Contents {
			if object.Key == nil {
				deleteObjectErr = fmt.Errorf("unexpected nil file key: %+v", object)
				break
			}
			objectIds = append(objectIds, types.ObjectIdentifier{Key: object.Key})
		}

		_, err = fur.client.DeleteObjects(ctx, &awsS3.DeleteObjectsInput{
			Bucket: aws.String(bucket),
			Delete: &types.Delete{Objects: objectIds},
		})
		if err != nil {
			return err
		}
	}

	return deleteObjectErr
}

// GenerateCloudFrontUrlForAudio comment.
// en: create url for audio content.
func (fur *S3Repository) GenerateCloudFrontUrlForAudio(domain, path *string) (*string, error) {
	if domain == nil || path == nil {
		return nil, nil
	}

	urlParsed, err := url.Parse(*domain)
	if err != nil {
		return nil, err
	}

	urlPath, err := url.Parse(urlParsed.Scheme + "://" + urlParsed.Host + utils.NormalizeRelativePath(*path))
	if err != nil {
		return nil, err
	}

	return utils.NewStringPointer(urlPath.String()), nil
}

// UploadFileWithContext comment.
// en: upload file with context to s3
func (fur *S3Repository) UploadFileWithContext(ctx context.Context, bucket string, key string, customFile CustomFile) error {
	fileName := filepath.Base(customFile.File.Name())
	if customFile.ContentType == nil {
		customFile.ContentType = utils.NewStringPointer(mime.TypeByExtension(filepath.Ext(fileName)))
	}
	contentDisposition := fmt.Sprintf(`attachment;filename="%v"`, fileName)
	_, err := fur.clientWithOnlyErrorLogger.PutObject(ctx, &awsS3.PutObjectInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(key),
		ContentDisposition: aws.String(contentDisposition),
		Body:               customFile.File,
		ContentType:        aws.String(*customFile.ContentType),
	})
	return err
}

// GenerateCloudFrontUrlCDNWithoutTime comment
// en: Generate logo_image_url cdn without time
func (fur *S3Repository) GenerateCloudFrontUrlCDNWithoutTime(domain, path *string) (*string, error) {
	if domain == nil || path == nil {
		return nil, nil
	}
	urlParsed, err := utils.ParseUrl(*domain)
	if err != nil {
		return nil, err
	}

	urlBuilder := urlParsed.Scheme + "://"
	if urlParsed.Subdomain != "" {
		urlSubdomain := strings.ReplaceAll(urlParsed.Subdomain, ".", "-")
		urlBuilder += urlSubdomain + "." + prefixCDN
	} else {
		urlBuilder += prefixCDN
	}
	urlBuilder += urlParsed.Domain + "." + urlParsed.TLD + utils.NormalizeRelativePath(*path)

	urlPath, err := url.Parse(urlBuilder)
	if err != nil {
		return nil, err
	}

	return utils.NewStringPointer(urlPath.String()), nil
}

// UploadFileStream comment.
// en: upload streaming data to s3
func (fur *S3Repository) UploadFileStream(ctx context.Context, bucket string, key string, streamData *io.PipeReader, contentType string, filename string) error {
	contentDisposition := fmt.Sprintf(`attachment;filename="%v"`, filename)
	_, err := fur.s3Uploader.Upload(ctx, &awsS3.PutObjectInput{
		Bucket:             aws.String(bucket),
		Key:                aws.String(key),
		ContentDisposition: aws.String(contentDisposition),
		Body:               streamData,
		ContentType:        aws.String(contentType),
	})
	return err
}

// GetContentTypeAndUrl comment
// en: get content_type and url with cdn
// en: input: bucket, domain, path
// en: output: content_type, url
func (fur *S3Repository) GetContentTypeAndUrl(ctx context.Context, bucket string, domain, path *string) (*string, *string, error) {
	// en: check if domain or path is nil, return nil
	if domain == nil || path == nil {
		return nil, nil, nil
	}

	// en: parse the domain URL to ensure it's valid
	urlParsed, err := url.Parse(*domain)
	if err != nil {
		return nil, nil, err
	}

	// en: construct the full URL path using the normalized path.
	urlPath, err := url.Parse(urlParsed.Scheme + "://" + urlParsed.Host + utils.NormalizeRelativePath(*path))
	if err != nil {
		return nil, nil, err
	}

	// en: Attempt to retrieve the object's metadata from S3
	lastObject, err := fur.GetHeadObject(ctx, bucket, *path)
	if err != nil {
		var nsk *types.NoSuchKey
		// en: check if the error is a NoSuchKey error (object not found)
		if errors.As(err, &nsk) {
			return nil, nil, nil
		}

		var apiErr smithy.APIError
		// en: check if the error is an API error
		if errors.As(err, &apiErr) {
			code := apiErr.ErrorCode()
			// en: If the error code is "NotFound", return nil
			if code == codeNotFound {
				return nil, nil, nil
			}
			return nil, nil, err
		}

		return nil, nil, err
	}

	// en: if the lastObject is nil, return nil
	if lastObject == nil {
		return nil, nil, nil
	}

	// en: Add the last modified time as a query parameter to the URL
	v := urlPath.Query()
	v.Add("time", fmt.Sprintf("%d", lastObject.LastModified.Unix()))
	urlPath.RawQuery = v.Encode() // en: encode the query parameters back into the URL

	// en: Return the constructed URL and the content_type of the object
	return utils.NewStringPointer(urlPath.String()), lastObject.ContentType, nil
}

// DownloadFileToLocal comment.
// en: download s3 file to local and can specify save directory
func (fur *S3Repository) DownloadFileToLocal(ctx context.Context, bucket string, key string, outputFile string) (string, error) {
	nameFile := os.TempDir() + "/" + outputFile
	folderPath := filepath.Dir(outputFile)
	if _, err := os.Stat(os.TempDir() + "/" + folderPath); os.IsNotExist(err) {
		err := os.MkdirAll(os.TempDir()+"/"+folderPath, utils.ExecutableFilePermission)
		if err != nil {
			return "", err
		}
	}
	file, err := os.Create(filepath.Clean(nameFile))
	if err != nil {
		return "", err
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Error(ctx, utils.ErrorCloseFile, "error", err)
		}
	}()

	_, err = fur.s3Downloader.Download(ctx, file, &awsS3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	return file.Name(), err
}
