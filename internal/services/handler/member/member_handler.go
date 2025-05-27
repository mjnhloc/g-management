package member

import "g-management/pkg/shared/handler"

type HTTPHandler struct {
	handler.ApplicationHandler
}

func NewHTTPHandler(
	ah handler.ApplicationHandler,
) *HTTPHandler {
	return &HTTPHandler{
		ApplicationHandler: ah,
	}
}
