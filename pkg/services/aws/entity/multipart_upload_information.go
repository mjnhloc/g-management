package entity

type MultipartUploadInformation struct {
	Method string
	Parts  []PartUploadInformation
}

type PartUploadInformation struct {
	OffsetFrom int64
	OffsetTo   int64
	Url        string
}
