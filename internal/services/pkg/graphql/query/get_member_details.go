package query

import (
	"g-management/internal/models/members/pkg/repository"
	"g-management/internal/services/pkg/graphql/output"

	"github.com/graphql-go/graphql"
)

func NewGetMemberDetailsQuery(
	types map[string]*graphql.Object,
	membersRepository repository.MembersRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        types["member"],
		Description: "Get member details by ID",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: output.BigInt,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return membersRepository.TakeByConditions(params.Context, map[string]interface{}{
				"id": params.Args["id"].(int),
			})
		},
	}
}
