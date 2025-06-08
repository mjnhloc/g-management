package mutation

import (
	"errors"
	"g-management/internal/models/classes/pkg/entity"
	"g-management/internal/models/classes/pkg/repository"
	"g-management/internal/services/pkg/graphql/output"
	"g-management/pkg/shared/utils"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewPutClassMutation(
	types map[string]*graphql.Object,
	db *gorm.DB,
	classesRepository repository.ClassesRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        types["class"],
		Description: "Update class information",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: output.BigInt,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			updateClassAttributes := map[string]interface{}{}
			classID := params.Args["id"].(int)
			updateClassAttributes["id"] = classID

			classInput := utils.GetSubMap(params.Source, "class")
			if classInput == nil {
				return nil, errors.New("No class input provided")
			}

			trainerIDPtr := utils.GetSubInteger(classInput, "trainer", "id")
			if trainerIDPtr != nil {
				trainerID := utils.DerefInt(trainerIDPtr)
				_, err := classesRepository.TakeByConditions(params.Context, map[string]interface{}{
					"id": trainerID,
				})
				if err != nil {
					return nil, err
				}

				updateClassAttributes["trainer_id"] = trainerID
			}

			classInputAttributes := utils.GetOnlyScalar(classInput)

			if classInputAttributes["name"] != nil {
				updateClassAttributes["name"] = classInputAttributes["name"].(string)
			}
			if classInputAttributes["schedule"] != nil {
				updateClassAttributes["schedule"] = classInputAttributes["schedule"].(string)
			}
			if classInputAttributes["duration"] != nil {
				updateClassAttributes["duration"] = classInputAttributes["duration"].(float64)
			}
			if classInputAttributes["max_capacity"] != nil {
				updateClassAttributes["max_capacity"] = classInputAttributes["max_capacity"].(float64)
			}
			if classInputAttributes["description"] != nil {
				updateClassAttributes["description"] = classInputAttributes["description"].(string)
			}

			var class entity.Classes
			var err error
			if err := utils.Transaction(params.Context, db, func(tx *gorm.DB) error {
				class, err = classesRepository.UpsertWithTransaction(tx, updateClassAttributes)
				if err != nil {
					return err
				}

				return nil
			}); err != nil {
				return nil, err
			}

			return class, nil
		},
	}
}
