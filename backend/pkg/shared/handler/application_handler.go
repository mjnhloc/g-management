package handler

type ApplicationHandler struct {
	BaseHTTPHandler
}

func NewApplicationHandler() *ApplicationHandler {
	return &ApplicationHandler{
		BaseHTTPHandler: BaseHTTPHandler{},
	}
}
