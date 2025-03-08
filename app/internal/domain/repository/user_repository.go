package repository

import "github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"

type IUserRepository interface {
	FindByName(name string) (*user.User, error)
	Save(user *user.User) error
}
