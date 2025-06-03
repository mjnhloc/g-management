package query

import (
	"g-management/internal/models/trainers/pkg/repository"
	"g-management/internal/services/pkg/graphql/output"

	"github.com/graphql-go/graphql"
)

func NewGetTrainerDetailsQuery(
	types map[string]*graphql.Object,
	trainersRepository repository.TrainersRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        types["trainer"],
		Description: "Get trainer details by ID",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: output.BigInt,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return trainersRepository.TakeByConditions(params.Context, map[string]interface{}{
				"id": params.Args["id"].(int),
			})
		},
	}
}
