package data_model

import "github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"

type User struct {
	ID   string
	Name string
}

func ToUserDataModel(from *user.User) *User {
	return &User{
		ID:   from.ID().String(),
		Name: from.Name().String(),
	}
}
