package user

import (
	"fmt"

	"github.com/shun-ideguchi/golang-ddd/internal/application/usecase/user/command"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/service"
)

type createUsecase struct {
	userRepository repository.IUserRepository
	userService    service.UserService
}

func NewCreateUsecase(userRepository repository.IUserRepository, userService service.UserService) *createUsecase {
	return &createUsecase{userRepository: userRepository, userService: userService}
}

func (u *createUsecase) Execute(cmd *command.CreateCommand) error {
	user, err := user.NewUser(cmd.UserID, cmd.Name, cmd.Email)
	if err != nil {
		return err
	}

	if u.userService.IsExists(user) {
		return fmt.Errorf("duplicate user: %s", cmd.Name)
	}

	return u.userRepository.Save(user)
}
