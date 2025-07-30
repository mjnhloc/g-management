package container

import (
	classesContainer "g-management/internal/models/classes/pkg/container"
	memberAttendancesContainer "g-management/internal/models/member-attendances/pkg/container"
	membersContainer "g-management/internal/models/members/pkg/container"
	membershipsContainer "g-management/internal/models/memberships/pkg/container"
	paymentsContainer "g-management/internal/models/payments/pkg/container"
	trainersContainer "g-management/internal/models/trainers/pkg/container"
	"g-management/pkg/services/elasticsearch/client"

	"gorm.io/gorm"
)

type RepositoryContainers struct {
	ClassesContainer           classesContainer.RepositoryContainer
	MemberAttendancesContainer memberAttendancesContainer.RepositoryContainer
	MembersContainer           membersContainer.RepositoryContainer
	MembershipsContainer       membershipsContainer.RepositoryContainer
	PaymentsContainer          paymentsContainer.RepositoryContainer
	TrainersContainer          trainersContainer.RepositoryContainer
	WaitlistsContainer         classesContainer.RepositoryContainer
	ElasticSearchClientRepo    client.ClientInterface
}

func NewRepositoryContainers(db *gorm.DB, client client.ClientInterface) *RepositoryContainers {
	return &RepositoryContainers{
		ClassesContainer:           classesContainer.NewRepositoryContainer(db),
		MemberAttendancesContainer: memberAttendancesContainer.NewRepositoryContainer(db),
		MembersContainer:           membersContainer.NewRepositoryContainer(db),
		MembershipsContainer:       membershipsContainer.NewRepositoryContainer(db),
		PaymentsContainer:          paymentsContainer.NewRepositoryContainer(db),
		TrainersContainer:          trainersContainer.NewRepositoryContainer(db),
		WaitlistsContainer:         classesContainer.NewRepositoryContainer(db),
		ElasticSearchClientRepo:    client,
	}
}
