package user

import (
	"errors"

	"github.com/shun-ideguchi/golang-ddd/internal/application/usecase/user/command"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
)

type deleteUsecase struct {
	userRepository repository.IUserRepository
}

func NewDeleteUsecase(userRepository repository.IUserRepository) *deleteUsecase {
	return &deleteUsecase{userRepository: userRepository}
}

func (u *deleteUsecase) Execute(cmd *command.DeleteCommand) error {
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
	// ユーザーが見つからない場合でも成功とする場合
	// if targetUser == nil {
	// 	return nil
	// }

	return u.userRepository.Delete(targetUser)
}
