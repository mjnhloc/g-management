package entity

type FileUploadInformation struct {
	Method string
	Url    string
}

type MultipartUploadInformation struct {
	Method string
	Parts  []PartUploadInformation
}

type PartUploadInformation struct {
	OffsetFrom int64
	OffsetTo   int64
	Url        string
}
