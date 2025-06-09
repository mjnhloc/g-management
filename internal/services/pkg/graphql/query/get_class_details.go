package query

import (
	"g-management/internal/models/classes/pkg/repository"
	"g-management/internal/services/pkg/graphql/output"

	"github.com/graphql-go/graphql"
)

func NewGetClassDetailsQuery(
	types map[string]*graphql.Object,
	classesRepository repository.ClassesRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        types["class"],
		Description: "Get class details by ID",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: output.BigInt,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return classesRepository.TakeByConditions(params.Context, map[string]interface{}{
				"id": params.Args["id"].(int),
			})
		},
	}
}
