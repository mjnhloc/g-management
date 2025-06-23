package mutation

import (
	"g-management/internal/models/trainers/pkg/repository"
	"g-management/internal/services/pkg/graphql/output"
	"g-management/pkg/shared/utils"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewDeleteTrainerMutation(
	typeVoid *graphql.Scalar,
	db *gorm.DB,
	trainersRepository repository.TrainersRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        typeVoid,
		Description: "Delete a trainer by ID",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: output.BigInt,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			trainerID := params.Args["id"].(int)

			_, err := trainersRepository.TakeByConditions(params.Context, map[string]interface{}{
				"id": trainerID,
			})
			if err != nil {
				return nil, err
			}

			return nil, utils.Transaction(params.Context, db, func(db *gorm.DB) error {
				err = trainersRepository.DeleteByConditions(params.Context, map[string]interface{}{
					"id": trainerID,
				})
				if err != nil {
					return err
				}

				return nil
			})
		},
	}
}
