package user

import (
	"errors"

	dto "github.com/shun-ideguchi/golang-ddd/internal/application/dto/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
)

type getUsecase struct {
	userRepository repository.IUserRepository
}

func NewGetUsecase(userRepository repository.IUserRepository) *getUsecase {
	return &getUsecase{userRepository: userRepository}
}

func (u *getUsecase) Execute(userID string) (*dto.UserDto, error) {
	targetID, err := user.NewUserID(userID)
	if err != nil {
		return nil, err
	}

	user, err :=u.userRepository.FindByID(targetID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return dto.NewUserDto(user), nil
}
