package output

import (
	"g-management/internal/models/members/pkg/entity"
	"g-management/internal/models/memberships/pkg/repository"
	"g-management/pkg/shared/utils"

	"github.com/graphql-go/graphql"
)

func NewMemberType(
	types map[string]*graphql.Object,
	membershipsRepository repository.MembershipsRepositoryInterface,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "member",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type: BigInt,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Members).ID, nil
					},
				},
				"name": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Members).Name, nil
					},
				},
				"email": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return utils.DerefString(params.Source.(entity.Members).Email), nil
					},
				},
				"phone": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Members).Phone, nil
					},
				},
				"date_of_birth": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return utils.DerefString(params.Source.(entity.Members).DateOfBirth), nil
					},
				},
				"is_active": &graphql.Field{
					Type: graphql.Boolean,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Members).IsActive, nil
					},
				},
				"membership": &graphql.Field{
					Type: types["membership"],
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return membershipsRepository.TakeByConditions(params.Context, map[string]interface{}{
							"member_id": params.Source.(entity.Members).ID,
						})
					},
				},
			}
		}),
	})
}
