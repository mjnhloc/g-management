package trainer

import (
	"g-management/pkg/shared/handler"

	"github.com/graphql-go/graphql"
)

const (
	PostNewTrainer = "trainer/post_new_trainer.json"
	PutTrainerInfo = "trainer/put_trainer_info.json"
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
