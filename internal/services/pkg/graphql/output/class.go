package output

import (
	"g-management/internal/models/classes/pkg/entity"
	"g-management/internal/models/trainers/pkg/repository"
	"g-management/pkg/shared/utils"

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
						return utils.DerefString(params.Source.(entity.Classes).Description), nil
					},
				},
				"trainer_id": &graphql.Field{
					Type: BigInt,
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						return params.Source.(entity.Classes).TrainerID, nil
					},
				},
				"trainer": &graphql.Field{
					Type: types["trainer"],
					Resolve: func(params graphql.ResolveParams) (interface{}, error) {
						trainer, err := trainersRepository.TakeByConditions(params.Context, map[string]interface{}{
							"id": params.Source.(entity.Classes).TrainerID,
						})
						if err != nil {
							return nil, err
						}

						return trainer, nil
					},
				},
			}
		}),
	})
}
