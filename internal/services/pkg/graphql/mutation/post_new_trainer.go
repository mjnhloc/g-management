package mutation

import (
	"g-management/internal/models/trainers/pkg/entity"
	"g-management/internal/models/trainers/pkg/repository"
	"g-management/pkg/shared/utils"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewPostNewTrainerMutation(
	outputTypes map[string]*graphql.Object,
	db *gorm.DB,
	trainersRepository repository.TrainersRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        outputTypes["trainer"],
		Description: "Create a new trainer",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			trainerAttributes := map[string]interface{}{}
			trainerInput := utils.GetSubMap(params.Source, "trainer")
			trainerInputAttributes := utils.GetOnlyScalar(trainerInput)
			if trainerInputAttributes["name"] != nil {
				trainerAttributes["name"] = trainerInputAttributes["name"].(string)
			}
			if trainerInputAttributes["email"] != nil {
				trainerAttributes["email"] = trainerInputAttributes["email"].(string)
			}
			if trainerInputAttributes["phone"] != nil {
				trainerAttributes["phone"] = trainerInputAttributes["phone"].(string)
			}
			if trainerInputAttributes["specialization"] != nil {
				trainerAttributes["specialization"] = trainerInputAttributes["specialization"].(string)
			}
			if trainerInputAttributes["hired_at"] != nil {
				trainerAttributes["hired_at"] = trainerInputAttributes["hired_at"].(string)
			}

			var trainer entity.Trainers
			var err error
			if err := utils.Transaction(params.Context, db, func(tx *gorm.DB) error {
				trainer, err = trainersRepository.CreateWithTransaction(tx, trainerAttributes)
				if err != nil {
					return err
				}

				return nil
			}); err != nil {
				return nil, err
			}

			return trainer, nil
		},
	}
}
