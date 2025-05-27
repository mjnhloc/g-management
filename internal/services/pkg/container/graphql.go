package container

import (
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewGraphqlSchema(
	repositories *RepositoryContainers,
	db *gorm.DB,
) (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{}),
	})
}
