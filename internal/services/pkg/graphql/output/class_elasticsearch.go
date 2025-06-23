package output

import (
	"g-management/internal/models/classes/pkg/entity"
	"g-management/pkg/shared/utils"

	"github.com/graphql-go/graphql"
)

func NewClassElasticsearchType(
	types map[string]*graphql.Object,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "class_elasticsearch",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type: BigInt,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.ClassDocument).ID, nil
					},
				},
				"name": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.ClassDocument).Name, nil
					},
				},
				"schedule": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.ClassDocument).Schedule, nil
					},
				},
				"description": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return utils.DerefString(params.Source.(entity.ClassDocument).Description), nil
					},
				},
			}
		}),
	})
}
