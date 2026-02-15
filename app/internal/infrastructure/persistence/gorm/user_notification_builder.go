package gorm

import "github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"

type UserDataModelBuilder struct {
	id    user.UserID
	name  user.Name
	email user.Email
}

func (b *UserDataModelBuilder) ID(id user.UserID) user.UserNotification {
	b.id = id
	return b
}

func (b *UserDataModelBuilder) Name(name user.Name) user.UserNotification {
	b.name = name
	return b
}

func (b *UserDataModelBuilder) Email(email user.Email) user.UserNotification {
	b.email = email
	return b
}

func (b *UserDataModelBuilder) Build() *User {
	return &User{
		ID:    b.id.String(),
		Name:  b.name.String(),
		Email: b.email.String(),
	}
}
