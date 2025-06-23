package query

import (
	"g-management/internal/models/trainers/pkg/repository"

	"github.com/graphql-go/graphql"
)

func NewGetAllTrainersQuery(
	types map[string]*graphql.Object,
	trainersRepository repository.TrainersRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(types["trainer"]),
		Description: "Get all trainers",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return trainersRepository.FindByConditions(params.Context, nil)
		},
	}
}
