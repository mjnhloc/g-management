package graphql

import (
	"g-management/internal/services/pkg/container"
	"g-management/internal/services/pkg/graphql/query"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func InitializeQueries(
	repositories *container.RepositoryContainers,
	db *gorm.DB,
	outputTypes map[string]*graphql.Object,
) *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"get_all_classes": query.NewGetAllClassesQuery(
				outputTypes,
				repositories.ClassesContainer.ClassesRepository,
			),
			"get_all_members": query.NewGetAllMembersQuery(
				outputTypes,
				repositories.MembersContainer.MembersRepository,
			),
			"get_all_trainers": query.NewGetAllTrainersQuery(
				outputTypes,
				repositories.TrainersContainer.TrainersRepository,
			),
			"get_class_details": query.NewGetClassDetailsQuery(
				outputTypes,
				repositories.ClassesContainer.ClassesRepository,
			),
			"get_member_details": query.NewGetMemberDetailsQuery(
				outputTypes,
				repositories.MembersContainer.MembersRepository,
			),
			"get_trainer_details": query.NewGetTrainerDetailsQuery(
				outputTypes,
				repositories.TrainersContainer.TrainersRepository,
			),
		},
	})
}
