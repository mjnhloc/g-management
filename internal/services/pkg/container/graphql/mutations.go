package graphql

import (
	"g-management/internal/services/pkg/container"
	"g-management/internal/services/pkg/graphql/mutation"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func InitializeMutations(
	repositories *container.RepositoryContainers,
	services *container.ServiceContainers,
	db *gorm.DB,
	outputTypes map[string]*graphql.Object,
	typeVoid *graphql.Scalar,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"post_new_class": mutation.NewPostNewClassMutation(
				outputTypes,
				db,
				repositories.TrainersContainer.TrainersRepository,
				repositories.ClassesContainer.ClassesRepository,
				services.ClassesContainer.ClassesService,
			),
			"post_new_member": mutation.NewPostNewMemberMutation(
				outputTypes,
				db,
				repositories.MembersContainer.MembersRepository,
				repositories.MembershipsContainer.MembershipsRepository,
				repositories.PaymentsContainer.PaymentsRepository,
			),
			"post_new_trainer": mutation.NewPostNewTrainerMutation(
				outputTypes,
				db,
				repositories.TrainersContainer.TrainersRepository,
			),
			"put_class_info": mutation.NewPutClassInfoMutation(
				outputTypes,
				db,
				repositories.ClassesContainer.ClassesRepository,
			),
			"delete_class": mutation.NewDeleteClassMutation(
				typeVoid,
				db,
				repositories.ClassesContainer.ClassesRepository,
			),
			"put_member_info": mutation.NewPutMemberInfoMutation(
				outputTypes,
				db,
				repositories.MembersContainer.MembersRepository,
			),
			"delete_member": mutation.NewDeleteMemberMutation(
				typeVoid,
				db,
				repositories.MembersContainer.MembersRepository,
			),
			"put_trainer_info": mutation.NewPutTrainerInfoMutation(
				outputTypes,
				db,
				repositories.TrainersContainer.TrainersRepository,
			),
			"delete_trainer": mutation.NewDeleteTrainerMutation(
				typeVoid,
				db,
				repositories.TrainersContainer.TrainersRepository,
			),
		},
	})
}
