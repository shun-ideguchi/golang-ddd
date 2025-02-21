package service

import (
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
)

type userService struct {
	userRepository repository.IUserRepository
}

func NewUserService(userRepository repository.IUserRepository) *userService {
	return &userService{userRepository: userRepository}
}

func (s *userService) IsExists(user *user.User) bool {
	// 重複を確認する処理
	found, err := s.userRepository.Find(user.Name().String())
	if err != nil {
		// errorを返すべきだが省略
		return false
	}

	if found == nil {
		return false
	}

	return true
}
