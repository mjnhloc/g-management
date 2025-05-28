package output

import (
	"g-management/internal/models/trainers/pkg/entity"

	"github.com/graphql-go/graphql"
)

func NewTrainerType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "trainer",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type: BigInt,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Trainers).ID, nil
					},
				},
				"name": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Trainers).Name, nil
					},
				},
				"email": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Trainers).Email, nil
					},
				},
				"phone_number": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Trainers).Phone, nil
					},
				},
				"specialization": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						if params.Source.(entity.Trainers).Specialization == nil {
							return nil, nil
						}
						return *params.Source.(entity.Trainers).Specialization, nil
					},
				},
			}
		}),
	})
}
