package query

import (
	"g-management/internal/models/members/pkg/repository"

	"github.com/graphql-go/graphql"
)

func NewGetAllMembersQuery(
	types map[string]*graphql.Object,
	membersRepository repository.MembersRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        graphql.NewList(types["member"]),
		Args:        graphql.FieldConfigArgument{},
		Description: "Get all members",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			return membersRepository.FindByConditions(params.Context, nil)
		},
	}
}
