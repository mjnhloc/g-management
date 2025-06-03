package container

import (
	"g-management/internal/services/pkg/graphql/mutation"
	"g-management/internal/services/pkg/graphql/output"
	"g-management/internal/services/pkg/graphql/query"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewGraphqlSchema(
	repositories *RepositoryContainers,
	db *gorm.DB,
) (graphql.Schema, error) {
	outputTypes := make(map[string]*graphql.Object)
	for _, graphqlType := range []*graphql.Object{
		output.NewClassType(
			outputTypes,
			repositories.TrainersContainer.TrainersRepository,
		),
		output.NewTrainerType(),
		output.NewMemberType(
			outputTypes,
			repositories.MembershipsContainer.MembershipsRepository,
		),
		output.NewMembershipType(
			outputTypes,
			repositories.PaymentsContainer.PaymentsRepository,
		),
		output.NewPaymentType(),
	} {
		outputTypes[graphqlType.Name()] = graphqlType
	}

	return graphql.NewSchema(graphql.SchemaConfig{
		Query: graphql.NewObject(graphql.ObjectConfig{
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
		}),
		Mutation: graphql.NewObject(graphql.ObjectConfig{
			Name: "Mutation",
			Fields: graphql.Fields{
				"post_new_class": mutation.NewPostNewClassMutation(
					outputTypes,
					db,
					repositories.TrainersContainer.TrainersRepository,
					repositories.ClassesContainer.ClassesRepository,
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
			},
		}),
	})
}
