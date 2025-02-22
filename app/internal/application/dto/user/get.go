package user

import "github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"

type UserDto struct {
	ID    string
	Name  string
	Email string
}

func NewUserDto(source *user.User) *UserDto {
	return &UserDto{
		ID:    source.ID().String(),
		Name:  source.Name().String(),
		Email: source.Email().String(),
	}
}
