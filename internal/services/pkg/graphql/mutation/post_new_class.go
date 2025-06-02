package mutation

import (
	"g-management/internal/models/classes/pkg/entity"
	classesRepo "g-management/internal/models/classes/pkg/repository"
	"g-management/internal/models/trainers/pkg/repository"
	"g-management/pkg/shared/utils"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewPostNewClassMutation(
	types map[string]*graphql.Object,
	db *gorm.DB,
	trainersRepository repository.TrainersRepositoryInterface,
	classesRepository classesRepo.ClassesRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        types["class"],
		Description: "Create a new class",
		Args:        graphql.FieldConfigArgument{},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			classAttributes := map[string]interface{}{}
			trainerIDPtr := utils.GetSubInteger(params.Source, "class", "trainer", "id")

			if trainerID := utils.DerefInt(trainerIDPtr); trainerIDPtr != nil {
				_, err := trainersRepository.TakeByConditions(params.Context, map[string]interface{}{
					"id": trainerID,
				})
				if err != nil {
					return nil, err
				}

				classAttributes["trainer_id"] = trainerID
			}

			classInput := utils.GetSubMap(params.Source, "class")
			classInputAttributes := utils.GetOnlyScalar(classInput)
			if classInputAttributes["name"] != nil {
				classAttributes["name"] = classInputAttributes["name"].(string)
			}

			if classInputAttributes["schedule"] != nil {
				classAttributes["schedule"] = classInputAttributes["schedule"].(string)
			}

			if classInputAttributes["duration"] != nil {
				classAttributes["duration"] = classInputAttributes["duration"].(float64)
			}
			if classInputAttributes["max_capacity"] != nil {
				classAttributes["max_capacity"] = classInputAttributes["max_capacity"].(float64)
			}
			if classInputAttributes["description"] != nil {
				classAttributes["description"] = classInputAttributes["description"].(string)
			}

			var class entity.Classes
			var err error
			if err := utils.Transaction(params.Context, db, func(tx *gorm.DB) error {
				class, err = classesRepository.CreateWithTransaction(tx, classAttributes)
				if err != nil {
					return err
				}

				return nil
			}); err != nil {
				return nil, err
			}

			return class, err
		},
	}
}
