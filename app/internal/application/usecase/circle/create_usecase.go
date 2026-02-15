package circle

import (
	"fmt"

	"github.com/shun-ideguchi/golang-ddd/internal/application/usecase/circle/command"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/circle"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/service"
)

type createUsecase struct {
	circleFactory    circle.ICircleFactory
	circleRepository repository.ICircleRepository
	circleService    service.CircleService
	userRepository   repository.IUserRepository
}

func NewCreateUsecase(
	circleFactory circle.ICircleFactory,
	circleRepository repository.ICircleRepository,
	circleService service.CircleService,
	userRepository repository.IUserRepository,
) *createUsecase {
	return &createUsecase{
		circleFactory:    circleFactory,
		circleRepository: circleRepository,
		circleService:    circleService,
		userRepository:   userRepository,
	}
}

func (u *createUsecase) Execute(cmd command.CreateCommand) error {
	ownerId, err := user.NewUserID(cmd.UserID)
	if err != nil {
		return fmt.Errorf("failed to create user ID: %w", err)
	}

	owner, err := u.userRepository.FindByName(ownerId.String())
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if owner == nil {
		return fmt.Errorf("user not found: %s", ownerId.String())
	}

	circleName, err := circle.NewCircleName(cmd.CircleName)
	if err != nil {
		return fmt.Errorf("failed to create circle name: %w", err)
	}

	circle, err := u.circleFactory.Create(circleName, owner)
	if err != nil {
		return fmt.Errorf("failed to create circle: %w", err)
	}

	if u.circleService.IsExists(circle) {
		return fmt.Errorf("duplicate circle: %s", circleName.String())
	}

	return u.circleRepository.Save(circle)
}
