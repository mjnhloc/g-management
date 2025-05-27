package mutation

import (
	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewPostClassMutation(
	types map[string]*graphql.Object,
	db *gorm.DB,
) *graphql.Field {
	return &graphql.Field{
		Type:        types["class"],
		Description: "Create a new class",
		Args:        graphql.FieldConfigArgument{},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return nil, nil
		},
	}
}
