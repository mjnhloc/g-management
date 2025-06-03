package output

import (
	"g-management/internal/models/memberships/pkg/entity"
	"g-management/internal/models/payments/pkg/repository"

	"github.com/graphql-go/graphql"
)

func NewMembershipType(
	types map[string]*graphql.Object,
	paymentsRepository repository.PaymentsRepositoryInterface,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "membership",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type: BigInt,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Memberships).ID, nil
					},
				},
				"member_id": &graphql.Field{
					Type: BigInt,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Memberships).MemberID, nil
					},
				},
				"membership_type": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Memberships).MembershipType, nil
					},
				},
				"start_date": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Memberships).StartDate, nil
					},
				},
				"end_date": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Memberships).EndDate, nil
					},
				},
				"payment": &graphql.Field{
					Type: types["payment"],
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return paymentsRepository.TakeByConditions(params.Context, map[string]interface{}{
							"membership_id": params.Source.(entity.Memberships).ID,
						})
					},
				},
			}
		}),
	})
}
