package mutation

import (
	"g-management/internal/models/members/pkg/entity"
	membersRepository "g-management/internal/models/members/pkg/repository"
	"g-management/internal/services/pkg/graphql/output"
	"g-management/pkg/shared/utils"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewPutMemberInfoMutation(
	types map[string]*graphql.Object,
	db *gorm.DB,
	membersRepository membersRepository.MembersRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        types["member"],
		Description: "Update member information",
		Args: graphql.FieldConfigArgument{
			"id": &graphql.ArgumentConfig{
				Type: output.BigInt,
			},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// en: handle member attributes
			memberAttributes := map[string]interface{}{}
			memberAttributes["id"] = params.Args["id"].(int)

			memberInput := utils.GetSubMap(params.Source, "member")
			memberInputAttributes := utils.GetOnlyScalar(memberInput)
			if memberInputAttributes["name"] != nil {
				memberAttributes["name"] = memberInputAttributes["name"].(string)
			}
			if memberInputAttributes["email"] != nil {
				memberAttributes["email"] = memberInputAttributes["email"].(string)
			}
			if memberInputAttributes["phone"] != nil {
				memberAttributes["phone"] = memberInputAttributes["phone"].(string)
			}
			if memberInputAttributes["date_of_birth"] != nil {
				memberAttributes["date_of_birth"] = memberInputAttributes["date_of_birth"].(string)
			}
			if memberInputAttributes["is_active"] != nil {
				memberAttributes["is_active"] = memberInputAttributes["is_active"].(bool)
			}

			member := entity.Members{}
			var err error
			if err := utils.Transaction(params.Context, db, func(tx *gorm.DB) error {
				member, err = membersRepository.UpsertWithTransaction(tx, memberAttributes)
				if err != nil {
					return err
				}

				return nil
			}); err != nil {
				return nil, err
			}

			return member, err
		},
	}
}
