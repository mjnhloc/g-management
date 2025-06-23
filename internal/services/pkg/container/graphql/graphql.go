package graphql

import (
	"g-management/internal/services/pkg/container"
	"g-management/internal/services/pkg/graphql/output"

	"github.com/graphql-go/graphql"
	"gorm.io/gorm"
)

func NewGraphqlSchema(
	repositories *container.RepositoryContainers,
	services *container.ServiceContainers,
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
		output.NewClassElasticsearchType(outputTypes),
	} {
		outputTypes[graphqlType.Name()] = graphqlType
	}

	voidOutputType := output.NewVoidType()

	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    InitializeQueries(repositories, db, outputTypes),
		Mutation: InitializeMutations(repositories, services, db, outputTypes, voidOutputType),
	})
}
