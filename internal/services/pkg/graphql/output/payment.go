package output

import (
	"g-management/internal/models/payments/pkg/entity"

	"github.com/graphql-go/graphql"
)

func NewPaymentType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "payment",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type: BigInt,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Payments).ID, nil
					},
				},
				"membership_id": &graphql.Field{
					Type: BigInt,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Payments).MembershipID, nil
					},
				},
				"price": &graphql.Field{
					Type: BigInt,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Payments).Price, nil
					},
				},
				"payment_date": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Payments).PaymentDate, nil
					},
				},
				"payment_method": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Payments).PaymentMethod, nil
					},
				},
				"status": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Payments).Status, nil
					},
				},
			}
		}),
	})
}
