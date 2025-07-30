package mutation

import (
	"errors"
	"time"

	"g-management/internal/models/activity_log"
	"g-management/internal/models/classes/pkg/entity"
	classRepo "g-management/internal/models/classes/pkg/repository"
	memberAttendanceRepo "g-management/internal/models/member-attendances/pkg/repository"
	"g-management/pkg/shared/utils"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewBookClassMutation(
	types map[string]*graphql.Object,
	db *gorm.DB,
	classesRepository classRepo.ClassesRepositoryInterface,
	waitlistsRepository classRepo.WaitlistsRepositoryInterface,
	memberAttendancesRepository memberAttendanceRepo.MemberAttendancesRepositoryInterface,
) *graphql.Field {
	return &graphql.Field{
		Type:        types["class"],
		Description: "Book a class or join waitlist if full",
		Args: graphql.FieldConfigArgument{
			"class_id":  &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
			"member_id": &graphql.ArgumentConfig{Type: graphql.NewNonNull(graphql.Int)},
		},
		Resolve: func(params graphql.ResolveParams) (interface{}, error) {
			ctx := params.Context
			classID := int64(params.Args["class_id"].(int))
			memberID := int64(params.Args["member_id"].(int))

			var bookedClass entity.Classes
			if err := utils.Transaction(ctx, db, func(tx *gorm.DB) error {
				// Check if already booked
				_, err := memberAttendancesRepository.TakeByConditions(ctx, map[string]interface{}{
					"class_id":  classID,
					"member_id": memberID,
				})
				if err == nil {
					return errors.New("You have already booked this class.")
				}

				// Check if already on waitlist
				_, err = waitlistsRepository.TakeByConditions(ctx, map[string]interface{}{
					"class_id":  classID,
					"member_id": memberID,
				})
				if err == nil {
					return errors.New("You are already on the waitlist for this class.")
				}

				// Get class info
				bookedClass, err = classesRepository.TakeByConditions(ctx, map[string]interface{}{"id": classID})
				if err != nil {
					return errors.New("Class not found")
				}

				// Count current attendance
				var attendanceCount int64
				err = tx.Table("member_attendances").Where("class_id = ?", classID).Count(&attendanceCount).Error
				if err != nil {
					return err
				}

				if int(attendanceCount) >= bookedClass.MaxCapacity {
					// Add to waitlist
					waitlist := entity.Waitlist{
						ClassID:  classID,
						MemberID: memberID,
						JoinedAt: time.Now().Format(time.RFC3339),
						Notified: false,
					}
					err = waitlistsRepository.Create(ctx, &waitlist)
					if err != nil {
						return errors.New("Failed to join waitlist")
					}
					// Log activity (optional: implement activity log repository)
					_ = tx.Table("activity_logs").Create(&activity_log.ActivityLog{
						UserID:     memberID,
						UserRole:   "member",
						Action:     "waitlist_joined",
						TargetType: "class",
						TargetID:   classID,
						Details:    "Joined waitlist for full class",
						CreatedAt:  time.Now().Format(time.RFC3339),
					}).Error
					return errors.New("Class is full. You have been added to the waitlist.")
				}

				// Book class
				attendance := map[string]interface{}{
					"class_id":    classID,
					"member_id":   memberID,
					"attended_at": time.Now().Format(time.RFC3339),
				}
				err = tx.Table("member_attendances").Create(&attendance).Error
				if err != nil {
					return errors.New("Failed to book class")
				}
				// Log activity
				_ = tx.Table("activity_logs").Create(&activity_log.ActivityLog{
					UserID:     memberID,
					UserRole:   "member",
					Action:     "class_booked",
					TargetType: "class",
					TargetID:   classID,
					Details:    "Booked class successfully",
					CreatedAt:  time.Now().Format(time.RFC3339),
				}).Error
				return nil
			}); err != nil {
				return nil, err
			}
			return bookedClass, nil
		},
	}
}
