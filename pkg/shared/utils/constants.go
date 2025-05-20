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

	ErrorInputFail = "MSGCM002"

	ErrorInputRequired = "MSGCM001"

	ErrorPasswordFail = "MSGCM005"
	ErrorEmailFail    = "MSGCM003"

	ErrorRegisteredPassword = "MSGNFJO03001"

	ErrorPasswordIsIncorrect = "MSGCM004"

	ErrorPasswordDuplicate  = "MSGCM009"
	ErrorTimeEndBeforeStart = "MSGCM007"

	ErrorInputCharacterLimit = "MSGCM008"
	ErrorInputByteLimit      = "MSGCM013"

	ErrorEmailAddressIsDuplicated = "MSGPFPM02001"

	ErrorRoleNameIsDuplicated = "MSGPFMC11002"

	ErrorAccountPermission = "MSGPFAU05002"

	ErrorStartAfterEnd = "MSGPFME01001"

	ErrorMissingFanclubBillingPlan         = "MCGPFMC09009"
	ErrorMissingFreeTrialSetting           = "MCGPFMC09010"
	ErrorMissingPersonalInfoPolicy         = "MCGPFMC09011"
	ErrorMissingCommissionRate             = "MCGPFMC09012"
	ErrorSiteHostingTypeIncorrect          = "MCGPFMC09013"
	ErrorMissingSiteDomain                 = "MCGPFMC09014"
	ErrorSiteNameIncorrect                 = "MCGPFMC09015"
	ErrorCommissionRateForOtherIncorrect   = "MCGPFMC09016"
	ErrorCommissionRateForAppIncorrect     = "MCGPFMC09017"
	ErrorMissingContentProviderCopyright   = "MCGPFMC09018"
	ErrorMissingContentProviderMailAddress = "MCGPFMC09019"
	ErrorDisplayNickNameSettingIncorrect   = "MCGPFMC09020"
	ErrorPersonalInfoSettingIncorrect      = "MCGPFMC09021"
	ErrorMissingSiteBaseInfo               = "MCGPFMC09022"
	ErrorNicknameAndMailAddressAreRequired = "MSGPFMC12010"
	ErrorFileNotFound                      = "MSGPFMC12011"
	ErrorNotEnoughColumnInCSV              = "MSGPFMC12014"
	ErrorSiteCodeMismatch                  = "MSGPFMC12015"
	ErrorMemberNoAlreadyUsed               = "MSGPFMC12006"
	ErrorPersonalInfoIsMissing             = "MSGPFMC12007"
	ErrorWrongCountryOrPrefectureOrGender  = "MSGPFMC12008"
	ErrorAnswerIsRequired                  = "MSGPFMC12009"
	ErrorNotFoundBillingPlanAvailable      = "MSGPFMC12016"
	ErrorScheduledExpireAtIsInThePast      = "MSGPFMC12017"
	ErrorAccountRegistered                 = "MSGPFMC12018"
	ErrorDuplicateTermAtMaintenance        = "MSGPFNM02010"
	ErrorEmptyContentType                  = "MSGPFNM08001"
	ErrorBillingEndAtIsInThePast           = "MSGPFMC17001"

	ErrorDuplicateHashtagName = "MSGCM021"

	ErrorDuplicatePortalMenuHashtagName = "MSGPFMC20003"

	ErrorDuplicatePortalSiteHashtagName = "MSGPFMC20504"

	ErrorLimitedSiteHashtag = "MSGPFMC20501"

	ErrorLimitedSiteHashtagNew = "MSGPFMC20502"

	ErrorLimitedPortalBanners       = "MSGPFMC20202"
	ErrorLimitedPublicPortalBanners = "MSGPFMC20201"
	ErrorDuplicateSiteHashtagTitle  = "MSGPFMC20503"

	ErrorLimitedVideoTags = "MSGPFMC20301"

	ErrorLimitedPortalPickups       = "MSGPFMC20402"
	ErrorLimitedPublicPortalPickups = "MSGPFMC20401"

	ErrorFailToPostComment = "MSGNFAP01001"

	ErrorStartDateAfterEndDate = "MSGPFME01001"

	ErrorNoPerformancePlanToRevoke = "MSGPFME02010"

	ErrorNotFoundPayment = "MSGPFRP02001"

	ErrorFcSiteNotValid = "MSGPFNM02009"

	ErrorReComfirmMailFail = "MSGNFMP06001"

	ErrorLimitedHeadlines = "MSGPFMC20101"

	ErrorNotFollowed = "MSGNFCMP05003"

	ErrorPaidMember               = "MSGNFCM04001" // en: error code with paid member
	ErrorMemberDuringReliefPeriod = "MSGNFCM04002" // en: error code with member during relief period

	// LoginPath frontend
	LoginPath = "/login"

	// SignUpPath frontend
	SignUpPath = "/signup"

	// SetPasswordPath frontend
	SetPasswordPath = "/signup/password"

	// SetPasswordPath frontend
	SetNFCPasswordPath = "/join/password"

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

	FCSite = "FC"
	CPSite = "CP"
	PFSite = "PF"

	ErrorTokenExpired      = "Token is expired"
	ErrorTokenExpiredMsgJp = "ログイン期限が切れました"

	S3FolderSiteImagesPath   = "site_design/"
	S3FolderPublicPath       = "public_html/"
	S3FolderContentsPath     = "contents/"
	S3FolderPortalImagesPath = "portal_design/"

	JstreamSftpFolderPublicPath     = "/public_html"
	JstreamSftpFolderArchivedPath   = "/archived"
	JstreamSftpFolderTranscodedPath = "/transcoded"

	MaxVodFileSize   = 30 * 1024 * 1024 * 1024
	MaxAudioFileSize = 12 * 1024 * 1024 * 1024

	VodExtension = ".mp4"

	ErrInvalidS3Bucket = "invalid s3 bucket"

	CpZendeskLink             = "https://partner-support.sheeta.com/"
	FcZendeskChannelPlusLink  = "https://help.nicochannel.jp/hc/ja"
	FcZendeskOtherChannelLink = "https://help.sheeta.com/"
	PfZendeskLink             = "https://help.sheeta.com/"

	FcCurrentPlanPath = "/my/membership-plan"

	ChannelPlusHostName = "ニコニコチャンネルプラス"
	NickNameDefault     = "ゲスト"

	ErrorGAInvalid  = "MSGCPSM06007"
	ErrorGTMInvalid = "MSGCPSM06008"

	ErrorUserLoginDisable = "MSGNFAU01005"

	LogEnqConfigReset = "EnqConfigReset"

	ErrorDuplicateEmail = "MSGCPSM05001"

	ErrCheckMaxLengthUnder50Characters = "MSGCP012"
	DefaultBatchSize                   = 100
	BatchSize1k                        = 1000

	ErrInvalidDomain = "MSGPFMC09027"

	ErrorInvalidHiragana = "MSGNFPI02004"

	DefaultProtocol    = "rtmp"
	DefaultHttpsScheme = "https"
	DefaultHttpScheme  = "http"

	SheetaPlatformID                       = "SHTA" // en: Platform ID for PF sheeta
	QloverPlatformID                       = "JOQR" // en: Platform ID for Qlover
	ChannelPlusPlatformID                  = "CHPL" // en: Platform ID for ChannelPlus
	TokyoFMPlatformID                      = "TKFM" // en: Platform ID for TOKYO FM Broadcasting
	KADOKAWABusinessContactsID             = 13     // en: KADOKAWA Business Contact
	YearMonthShowTheSalesAndCommissionData = "2024-01"
	YearShowDataAnnualSalesReport          = 2024
	ErrorMailAddressNotMatchDomain         = "MSGPFMC09029"
	YearMonthReleasePonp                   = "2024-04"

	FunctionCodeShowDisplayAdvertisementFlg   = "PFAM010-01"
	FunctionCodeUpdateDisplayAdvertisementFlg = "PFAM010-02"

	AppUserAgent = "SheetaLive"

	MaxByteContent = 65535

	SyncSlaveDBWaitTime = 1 // en: Maximum waiting time to sync from master db to slave db. Average is around 10〜15ms

	TrueStatusString  = "true"
	FalseStatusString = "false"
	NullStatusString  = "null"

	ExecutableFilePermission = 0o750

	SetReadHeaderTimeOut = 5

	ErrorCloseResponseBody = "failed to close response body"
	ErrorCloseRows         = "failed to close Rows"
	ErrorCloseReader       = "failed to close reader"
	ErrorCloseWriter       = "failed to close writer"
	ErrorCloseFile         = "failed to close file"
	ErrorCloseSftp         = "failed to close sftp"
)
