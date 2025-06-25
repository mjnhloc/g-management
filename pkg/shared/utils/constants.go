package utils

import (
	"errors"
	"time"
)

var ErrS3ObjectNotFound error = errors.New("object not found")

const (
	// FormatYear year
	FormatYear = "2006"

	// FormatMonth month
	FormatMonth = "01"

	// FormatDate time
	FormatDate = time.DateOnly

	// FormatFiscalMonthDate time
	FormatFiscalMonthDate = "01"

	// FormatTime time
	FormatTime = "15:04"

	// FormatTimeDB time
	FormatTimeDB = time.TimeOnly

	// FormatDateTime time
	FormatDateTime = time.RFC3339

	// FormatDateTimeDb time
	FormatDateTimeDb    = time.DateTime
	FormatDateTimeDbCSV = "2006-01-02-150405"

	FormatDateTimeDbTrimSec = "2006-01-02 15:04:00"

	FormatDateInline = "20060102150405"

	FormatSeparateDateTimeInline = "20060102_150405"

	FormatMonthYear = "2006-01"
	FormatYearMonth = "200601"

	FormatDateForEmails = "2006/01/02"

	FormatDateTimeForEmails = "2006/01/02 15:04"

	FormatDateTimeJST = "2006/01/02 15:04:05"

	FormatDateWithoutSeparate = "20060102"

	FormatMonthYearWithoutSeparate = "200601"

	FormatMonthDate = "01/02"

	FormatDateJStream = "20060102-15-04-05"

	FormatDayMonthYearDate = "02/01/2006"

	// Default response message of the API for unhandled errors
	MessageInternalServerError = "Internal server error"

	// Page default
	DefaultPage = 1

	// PerPage default
	DefaultPerPage = 30

	// Min page and per page
	MinPaging = 1

	// Max page and per page
	MaxPerPage = 100

	// LoginPath frontend
	LoginPath = "/login"

	// SignUpPath frontend
	SignUpPath = "/signup"

	// SetPasswordPath frontend
	SetPasswordPath = "/signup/password"

	// PaymentPath for register to plan
	PaymentPath = "/join/membership-plan"

	// Personal information path
	PersonalInfoPath = "/join/personal-information"

	// Set Forgot Password path
	SetForgotPasswordPath = "/renew-password"

	// Reset Password FC path
	ResetPasswordFCPath = "/password-reset/input"

	// Verified Mail path
	VerifiedMailFCPath = "/my/mail/change/complete"

	// Url active path
	UrlActivePath = "/my/mail/register/complete"

	// Url verify email path
	UrlVerifyEmailPath = "/login/login-redirect"

	// Live List path
	LiveListCPPath  = "/live"
	VideoListCPPath = "/video"
	ArticlesPath    = "/articles"

	ErrorS3BucketDoesNotExist = "S3 bucket not found"
	ErrorKCUserDoesNotExist   = "User not found"

	MySqlDuplicateErrorNumber = 1062

	JstreamSftpFolderPublicPath     = "/public_html"
	JstreamSftpFolderArchivedPath   = "/archived"
	JstreamSftpFolderTranscodedPath = "/transcoded"

	MaxVodFileSize   = 30 * 1024 * 1024 * 1024
	MaxAudioFileSize = 12 * 1024 * 1024 * 1024

	ErrInvalidS3Bucket = "invalid s3 bucket"

	MaxByteContent = 65535

	SyncSlaveDBWaitTime = 1 // en: Maximum waiting time to sync from master db to slave db. Average is around 10〜15ms

	TrueStatusString  = "true"
	FalseStatusString = "false"
	NullStatusString  = "null"

	ExecutableFilePermission = 0o750

	SetReadHeaderTimeOut = 5

	DefaultHttpPath = "https"

	ErrorCloseResponseBody = "failed to close response body"
	ErrorCloseRows         = "failed to close Rows"
	ErrorCloseReader       = "failed to close reader"
	ErrorCloseWriter       = "failed to close writer"
	ErrorCloseFile         = "failed to close file"
	ErrorCloseSftp         = "failed to close sftp"
	ErrorInputFail         = "input invalid"
	ErrorInputEmail        = "input email invalid"
	ErrorPasswordFail      = "input password invalid"
	ErrorDomain            = "input domain invalid"
	ErrorTokenExpired      = "token expired"
)
