package query

import (
	"g-management/internal/models/classes/pkg/repository"

	"github.com/graphql-go/graphql"
)

func NewGetAllClassesQuery(
	types map[string]*graphql.Object,
	classesRepository repository.ClassesRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(types["class"]),
		Description: "Get all classes",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return classesRepository.FindByConditions(params.Context, nil)
		},
	}
}
