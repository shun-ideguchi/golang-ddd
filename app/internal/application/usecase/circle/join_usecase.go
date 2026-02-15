package circle

import (
	"fmt"

	"github.com/shun-ideguchi/golang-ddd/internal/application/usecase/circle/command"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/circle"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/service"
)

type joinUsecase struct {
	circleRepository repository.ICircleRepository
	userRepository   repository.IUserRepository
	circleService    service.CircleService
}

func NewJoinUsecase(
	circleRepository repository.ICircleRepository,
	userRepository repository.IUserRepository,
	circleService service.CircleService,
) *joinUsecase {
	return &joinUsecase{
		circleRepository: circleRepository,
		userRepository:   userRepository,
		circleService:    circleService,
	}
}

func (u *joinUsecase) Execute(cmd command.JoinCommand) error {
	memberId, err := user.NewUserID(cmd.UserID)
	if err != nil {
		return fmt.Errorf("failed to create user ID: %w", err)
	}
	member, err := u.userRepository.FindByID(memberId)
	if err != nil {
		return fmt.Errorf("failed to find user: %w", err)
	}
	if member == nil {
		return fmt.Errorf("user not found: %s", memberId.String())
	}

	circleId, err := circle.NewCircleID(cmd.CircleID)
	if err != nil {
		return fmt.Errorf("failed to create circle ID: %w", err)
	}
	circle, err := u.circleRepository.FindByCircleID(circleId)
	if err != nil {
		return fmt.Errorf("failed to find circle: %w", err)
	}
	if circle == nil {
		return fmt.Errorf("circle not found: %s", circleId.String())
	}

	if err := circle.Join(member); err != nil {
		return fmt.Errorf("failed to join circle: %w", err)
	}
	return u.circleRepository.Save(circle)
}
