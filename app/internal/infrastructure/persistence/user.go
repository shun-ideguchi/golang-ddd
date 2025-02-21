package persistence

import (
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
)

type userPersistence struct {
}

func NewUserPersistence() repository.IUserRepository {
	return &userPersistence{}
}

func (p *userPersistence) Find(userName string) (*user.User, error) {
	// 再構築処理

	// 仮モデル作成
	user, _ := user.NewUser("uuid", "test name")
	return user, nil
}

func (p *userPersistence) Save(user *user.User) error {
	// 永続化処理

	return nil
}
