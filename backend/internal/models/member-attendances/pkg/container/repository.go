package container

import (
	"g-management/internal/models/member-attendances/pkg/repository"

	"gorm.io/gorm"
)

type RepositoryContainer struct {
	MemberAttendancesRepository repository.MemberAttendancesRepositoryInterface
}

func NewRepositoryContainer(db *gorm.DB) RepositoryContainer {
	return RepositoryContainer{
		MemberAttendancesRepository: repository.NewMemberAttendancesRepository(db),
	}
}
