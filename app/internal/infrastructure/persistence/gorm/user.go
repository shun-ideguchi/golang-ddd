package gorm

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
)

type userPersistence struct {
}

func NewUserPersistence() repository.IUserRepository {
	return &userPersistence{}
}

func (p *userPersistence) FindByID(id user.UserID) (*user.User, error) {

	target := User{
		ID:    id.String(),
		Name:  "test",
		Email: "test@test.com",
	}

	user := user.ReNewUser(target.ID, target.Name, target.Email)
	return user, nil
}

func (p *userPersistence) FindByName(userName string) (*user.User, error) {
	// DBから再構築したと仮定
	target := User{
		ID:    uuid.NewString(),
		Name:  userName,
		Email: "test@test.com",
	}

	// データモデルからドメインモデルを生成
	// ルールチェックを行わない理由はDBにはルールが適用された値が永続化されているため
	// 開発者が手動で更新するケースはドメインルールに沿った値を永続化すると決める
	user := user.ReNewUser(target.ID, target.Name, target.Email)
	return user, nil
}

func (p *userPersistence) Save(user *user.User) error {
	builder := &UserDataModelBuilder{}
	user.Notify(builder)    // ← domainに「中身教えて」と依頼
	data := builder.Build() // ← DBモデルに変換

	// 永続化処理
	fmt.Println(data)

	return nil
}
