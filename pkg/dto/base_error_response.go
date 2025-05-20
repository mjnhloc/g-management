package dto

// BaseResponse struct
type BaseErrorResponse struct {
	Error *ErrorResponse `json:"error"`
}

// ErrorResponse struct
type ErrorResponse struct {
	Message          interface{} `json:"message"`
	Details          interface{} `json:"details"`
	DebugInformation interface{} `json:"debug_information"`
}
