package output

import (
	"g-management/internal/models/classes/pkg/entity"
	"g-management/internal/models/trainers/pkg/repository"

	"github.com/graphql-go/graphql"
)

func NewClassType(
	types map[string]*graphql.Object,
	trainersRepository repository.TrainersRepositoryInterface,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "class",
		Fields: graphql.FieldsThunk(func() graphql.Fields {
			return graphql.Fields{
				"id": &graphql.Field{
					Type: BigInt,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Classes).ID, nil
					},
				},
				"name": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Classes).Name, nil
					},
				},
				"schedule": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Classes).Schedule, nil
					},
				},
				"duration": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Classes).Duration, nil
					},
				},
				"max_capacity": &graphql.Field{
					Type: graphql.Int,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Classes).MaxCapacity, nil
					},
				},
				"description": &graphql.Field{
					Type: graphql.String,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						if params.Source.(entity.Classes).Description == nil {
							return nil, nil
						}
						return *params.Source.(entity.Classes).Description, nil
					},
				},
				"trainer": &graphql.Field{
					Type: types["trainer"],
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return trainersRepository.TakeByConditions(params.Context, map[string]interface{}{
							"id": params.Source.(entity.Classes).TrainerID,
						})
					},
				},
			}
		}),
	})
}
