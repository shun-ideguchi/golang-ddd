package factory

import (
	"github.com/google/uuid"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
)

type userFactory struct {
}

func NewUserFactory() user.IFactory {
	return &userFactory{}
}

func (f *userFactory) Create(name, email string) (*user.User, error) {
	// ID生成が複雑なものと仮定
	userID := uuid.NewString()
	return user.NewUser(userID, name, email)
}
