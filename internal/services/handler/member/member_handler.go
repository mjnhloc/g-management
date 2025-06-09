package member

import (
	"g-management/pkg/shared/handler"

	"github.com/graphql-go/graphql"
)

const (
	PostNewMember = "member/post_new_member.json"
	PutMemberInfo = "member/put_member_info.json"
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
