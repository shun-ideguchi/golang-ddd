package user

import (
	"fmt"

	"github.com/shun-ideguchi/golang-ddd/internal/application/usecase/user/command"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/service"
)

type createUsecase struct {
	userFactory    user.IFactory
	userRepository repository.IUserRepository
	userService    service.UserService
}

func NewCreateUsecase(userFactory user.IFactory, userRepository repository.IUserRepository, userService service.UserService) *createUsecase {
	return &createUsecase{userFactory: userFactory, userRepository: userRepository, userService: userService}
}

func (u *createUsecase) Execute(cmd *command.CreateCommand) error {
	user, err := u.userFactory.Create(cmd.Name, cmd.Email)
	if err != nil {
		return fmt.Errorf("failed to create user entity: %w", err)
	}

	if u.userService.IsExists(user) {
		return fmt.Errorf("duplicate user: %s", cmd.Name)
	}

	return u.userRepository.Save(user)
}
