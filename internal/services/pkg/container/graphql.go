package container

import (
	"g-management/internal/services/pkg/graphql/mutation"
	"g-management/internal/services/pkg/graphql/query"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewGraphqlSchema(
	repositories *RepositoryContainers,
	db *gorm.DB,
) (graphql.Schema, error) {
	outputTypes := make(map[string]*graphql.Object)
	for _, graphqlType := range []*graphql.Object{} {
		outputTypes[graphqlType.Name()] = graphqlType
	}

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
			Name: "Query",
			Fields: graphql.Fields{
				"get_all_classes": query.NewGetAllClassesQuery(
					outputTypes,
					repositories.ClassesContainer.ClassesRepository,
				),
			},
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"post_class": mutation.NewPostClassMutation(
					outputTypes,
					db,
				),
			},
		}),
	})
}
