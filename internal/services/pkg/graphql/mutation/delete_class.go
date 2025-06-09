package mutation

import (
	"g-management/internal/models/classes/pkg/repository"
	"g-management/internal/services/pkg/graphql/output"
	"g-management/pkg/shared/utils"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewDeleteClassMutation(
	typeVoid *graphql.Scalar,
	db *gorm.DB,
	classesRepository repository.ClassesRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        typeVoid,
		Description: "Delete class by ID",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: output.BigInt,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			classID := params.Args["id"].(int)

			// Check if the class exists
			_, err := classesRepository.TakeByConditions(params.Context, map[string]interface{}{
				"id": classID,
			})
			if err != nil {
				return nil, err
			}

			return nil, utils.Transaction(params.Context, db, func(tx *gorm.DB) error {
				err = classesRepository.DeleteByConditions(params.Context, map[string]interface{}{
					"id": classID,
				})
				if err != nil {
					return err
				}

				return nil
			})
		},
	}
}
