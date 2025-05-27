package class

import (
	"g-management/pkg/shared/handler"

	"github.com/graphql-go/graphql"
)

const (
	PostClass = "post_class.json"
)

type HTTPHandler struct {
	handler.ApplicationHandler
	graphql graphql.Schema
}

func NewHTTPHandler(
	ah handler.ApplicationHandler,
	graphql graphql.Schema,
) *HTTPHandler {
	return &HTTPHandler{
		ApplicationHandler: ah,
		graphql:            graphql,
	}
}
