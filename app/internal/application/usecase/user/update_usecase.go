package user

import (
	"errors"

	"github.com/shun-ideguchi/golang-ddd/internal/application/usecase/user/command"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/service"
)

type updateUsecase struct {
	userRepository repository.IUserRepository
	userService    service.UserService
}

func NewUpdateUsecase(userRepository repository.IUserRepository, userService service.UserService) *updateUsecase {
	return &updateUsecase{userRepository: userRepository, userService: userService}
}

func (u *updateUsecase) Execute(cmd *command.UpdateCommand) error {
	targetID, err := user.NewUserID(cmd.UserID)
	if err != nil {
		return err
	}

	targetUser, err := u.userRepository.FindByID(targetID)
	if err != nil {
		return err
	}
	if targetUser == nil {
		return errors.New("user not found")
	}

	if cmd.Name != nil {
		newUserName, err := user.NewName(*cmd.Name)
		if err != nil {
			return err
		}
		targetUser.ChangeName(newUserName)
	}

	if cmd.Email != nil {
		newUserEmail, err := user.NewEmail(*cmd.Email)
		if err != nil {
			return err
		}
		targetUser.ChangeEmail(newUserEmail)
	}

	if u.userService.IsExists(targetUser) {
		return errors.New("duplicate user")
	}

	return u.userRepository.Save(targetUser)
}
