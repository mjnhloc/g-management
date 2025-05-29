package container

import (
	classHandler "g-management/internal/services/handler/class"
	memberHandler "g-management/internal/services/handler/member"
	trainerHandler "g-management/internal/services/handler/trainer"
	baseHandler "g-management/pkg/shared/handler"
	"g-management/pkg/shared/validator"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

type HandlerContainer struct {
	Classes  *classHandler.HTTPHandler
	Members  *memberHandler.HTTPHandler
	Trainers *trainerHandler.HTTPHandler
}

func NewHandlerContainer(
	inputValidator *validator.JsonSchemaValidator,
	graphql graphql.Schema,
	db *gorm.DB,
) HandlerContainer {
	base := baseHandler.NewApplicationHandler()
	base.Validator = inputValidator

	classContainer := classHandler.NewHTTPHandler(*base, graphql)

	memberContainer := memberHandler.NewHTTPHandler(*base)

	trainerContainer := trainerHandler.NewHTTPHandler(*base)

	return HandlerContainer{
		Classes:  classContainer,
		Members:  memberContainer,
		Trainers: trainerContainer,
	}
}
