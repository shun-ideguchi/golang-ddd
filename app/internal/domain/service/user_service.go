package service

import "github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"

type userService struct {
}

func NewUserService() *userService {
	return &userService{}
}

func (s *userService) IsExists(user *user.User) bool {
	// 重複を確認する処理
	return true
}