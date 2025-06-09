package mutation

import (
	"fmt"
	"time"

	"g-management/internal/models/members/pkg/entity"
	membersRepository "g-management/internal/models/members/pkg/repository"
	membershipsRepository "g-management/internal/models/memberships/pkg/repository"
	paymentsRepository "g-management/internal/models/payments/pkg/repository"
	"g-management/pkg/shared/utils"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewPostNewMemberMutation(
	types map[string]*graphql.Object,
	db *gorm.DB,
	membersRepository membersRepository.MembersRepositoryInterface,
	membershipsRepository membershipsRepository.MembershipsRepositoryInterface,
	paymentsRepository paymentsRepository.PaymentsRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        types["member"],
		Description: "Create a new member",
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			// en: handle member attributes
			memberAttributes := map[string]interface{}{}
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

			// en: handle membership attributes
			membershipAttributes := map[string]interface{}{}
			membershipInput := utils.GetSubMap(params.Source, "membership")
			membershipInputAttributes := utils.GetOnlyScalar(membershipInput)
			membershipType := membershipInputAttributes["membership_type"].(string)
			membershipAttributes["membership_type"] = membershipType

			var startDateStr string
			if memberInputAttributes["start_date"] != nil {
				startDateStr = membershipInputAttributes["start_date"].(string)
			} else {
				startDateStr = time.Now().Format(utils.FormatDate)
			}
			membershipAttributes["start_date"] = startDateStr

			var endDate time.Time
			startDate, err := time.Parse(utils.FormatDate, startDateStr)
			if err != nil {
				return nil, err
			}
			switch membershipType {
			case "monthly":
				endDate = startDate.AddDate(0, 1, 0)
			case "quarter":
				endDate = startDate.AddDate(0, 3, 0)
			case "annual":
				endDate = startDate.AddDate(1, 0, 0)
			default:
				return nil, fmt.Errorf("invalid membership type: %s", membershipType)
			}
			membershipAttributes["end_date"] = endDate.Format(utils.FormatDate)

			// en: handle payment attributes
			paymentAttributes := map[string]interface{}{}
			paymentInput := utils.GetSubMap(params.Source, "payment")
			paymentInputAttributes := utils.GetOnlyScalar(paymentInput)
			if paymentInputAttributes["price"] != nil {
				paymentAttributes["price"] = paymentInputAttributes["price"].(float64)
			}
			if paymentInputAttributes["payment_date"] != nil {
				paymentAttributes["payment_date"] = paymentInputAttributes["payment_date"].(string)
			}
			if paymentInputAttributes["payment_method"] != nil {
				paymentAttributes["payment_method"] = paymentInputAttributes["payment_method"].(string)
			}
			if paymentInputAttributes["status"] != nil {
				paymentAttributes["status"] = paymentInputAttributes["status"].(string)
			}

			var member entity.Members
			if err := utils.Transaction(params.Context, db, func(tx *gorm.DB) error {
				memberAttributes["is_active"] = true
				member, err = membersRepository.CreateWithTransaction(tx, memberAttributes)
				if err != nil {
					return err
				}

				membershipAttributes["member_id"] = member.ID
				membership, err := membershipsRepository.CreateWithTransaction(tx, membershipAttributes)
				if err != nil {
					return err
				}

				paymentAttributes["membership_id"] = membership.ID
				_, err = paymentsRepository.CreateWithTransaction(tx, paymentAttributes)
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
