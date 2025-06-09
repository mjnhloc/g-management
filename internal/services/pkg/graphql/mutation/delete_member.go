package mutation

import (
	"g-management/internal/models/members/pkg/repository"
	"g-management/internal/services/pkg/graphql/output"
	"g-management/pkg/shared/utils"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewDeleteMemberMutation(
	typeVoid *graphql.Scalar,
	db *gorm.DB,
	membersRepository repository.MembersRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        typeVoid,
		Description: "Delete a member by ID",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: output.BigInt,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			memberID := params.Args["id"].(int)

			_, err := membersRepository.TakeByConditions(params.Context, map[string]interface{}{
				"id": memberID,
			})
			if err != nil {
				return nil, err
			}

			return nil, utils.Transaction(params.Context, db, func(db *gorm.DB) error {
				err = membersRepository.DeleteByConditions(params.Context, map[string]interface{}{
					"id": memberID,
				})
				if err != nil {
					return err
				}

				return nil
			})
		},
	}
}
