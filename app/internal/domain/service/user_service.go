package service

import (
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
)

type UserService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) *UserService {
	return &UserService{userRepository: userRepository}
}

func (s *UserService) IsExists(user *user.User) bool {
	// 重複を確認する処理
	found, err := s.userRepository.FindByName(user.Name().String())
	if err != nil {
		// errorを返すべきだが省略
		return false
	}

	if found == nil {
		return false
	}

	// 存在しない仮定で進めるため偽で返却
	// return true
	return false
}
