# リポジトリ

### リポジトリとは
エンティティの永続化と再構築を抽象化するためのドメインオブジェクト   
ドメイン層にインターフェースを定義し、インフラストラクチャ層に実装を配置することで、依存関係逆転の原則（DIP）を実現する

### 役割
- ドメインモデルの永続化（保存）
- 永続化されたデータからドメインモデルを再構築（取得）
- ドメイン層が永続化の技術的詳細（DB、ファイル等）を知らなくて済むようにする

## リポジトリのインターフェース

リポジトリのインターフェースはドメイン層に定義する   
これにより、ドメイン層はインフラストラクチャ層に依存せず、永続化の手段を抽象的に扱える

```go
// app/internal/domain/repository/user_repository.go

package repository

import "github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"

type IUserRepository interface {
	Find(userName string) (*user.User, error)
	Save(user *user.User) error
}
```

ポイント:
- インターフェースの引数・戻り値にはドメインモデル（`*user.User`）を使用する
- SQL文やDBドライバなどの技術的な詳細はインターフェースに一切含めない
- ドメイン層に定義することで、ドメインロジックからリポジトリを利用する際に外部パッケージへの依存が発生しない

## リポジトリの実装

インターフェースの実装はインフラストラクチャ層に配置する   
実装ではデータベースとのやり取りを行い、ドメインモデルとデータモデルの変換を担う

```go
// app/internal/infrastructure/persistence/user.go

package persistence

import (
	"fmt"

	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
	"github.com/shun-ideguchi/golang-ddd/internal/infrastructure/data_model"
)

type userPersistence struct{}

func NewUserPersistence() repository.IUserRepository {
	return &userPersistence{}
}

func (p *userPersistence) Find(userName string) (*user.User, error) {
	// DBから再構築したと仮定
	target := data_model.User{
		ID:   "uuid",
		Name: userName,
	}

	// データモデルからドメインモデルを生成
	// ルールチェックを行わない理由はDBにはルールが適用された値が永続化されているため
	// 開発者が手動で更新するケースはドメインルールに沿った値を永続化すると決める
	user := user.ReNewUser(target.ID, target.Name)
	return user, nil
}

func (p *userPersistence) Save(user *user.User) error {
	data := data_model.ToUserDataModel(user)

	// 永続化処理
	fmt.Println(data)

	return nil
}
```

ポイント:
- `NewUserPersistence` の戻り値の型は `repository.IUserRepository`（インターフェース）にする。これにより呼び出し側は実装の詳細を知る必要がない
- `Find` ではDBから取得したデータモデルを `ReNewUser` でドメインモデルに再構築する。`NewUser` ではなく `ReNewUser` を使う理由は、DBに永続化された値は既にバリデーション済みであるため、再度ルールチェックを行う必要がないから
- `Save` ではドメインモデルをデータモデルに変換してから永続化する

## データモデル（DTO）

ドメインモデルと永続化層の間のデータ変換を担うオブジェクト   
ドメインモデルのフィールドは非公開（小文字始まり）であるため、直接アクセスできない   
データモデルを介すことで、ドメインモデルのカプセル化を保ちつつ永続化層がデータを扱える

```go
// app/internal/infrastructure/data_model/user.go

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
```

ポイント:
- データモデルのフィールドは公開（大文字始まり）にする。DBやJSONとのマッピングを容易にするため
- ドメインモデルが公開しているメソッド（`ID()`, `Name()`）を通じて値を取得し、`String()` で基本型に変換する
- ドメインモデルのカプセル化を破壊しないことが重要

## NewUser と ReNewUser の使い分け

エンティティには2つの生成方法がある

```go
// app/internal/domain/model/user/user.go

// NewUser は新規生成時に使用する。バリデーションを行い、不正な値の生成を防ぐ
func NewUser(userID, name string) (*User, error) {
	newUserID, err := newUserID(userID)
	if err != nil {
		return nil, err
	}
	newName, err := newName(name)
	if err != nil {
		return nil, err
	}

	return &User{
		userID: newUserID,
		name:   newName,
	}, nil
}

// ReNewUser はDBからの再構築時に使用する。バリデーションを行わない
func ReNewUser(ID, name string) *User {
	return &User{
		userID: UserID(ID),
		name:   Name(name),
	}
}
```

| メソッド | 用途 | バリデーション | 使用箇所 |
|----------|------|---------------|----------|
| `NewUser` | 新規作成 | あり | ユースケース層・アプリケーション層 |
| `ReNewUser` | DB からの再構築 | なし | リポジトリ実装（インフラストラクチャ層） |

## 依存関係の流れ

```
ドメイン層（domain）                  インフラストラクチャ層（infrastructure）
┌──────────────────────┐             ┌──────────────────────┐
│  IUserRepository     │◄─実装──────│  userPersistence     │
│  （インターフェース）   │             │  （リポジトリ実装）    │
└──────────────────────┘             └──────────┬───────────┘
         ▲                                      │
         │ 依存                                  │ 使用
         │                                      ▼
┌──────────────────────┐             ┌──────────────────────┐
│  userService         │             │  data_model.User     │
│  （ドメインサービス）   │             │  （データモデル）      │
└──────────────────────┘             └──────────────────────┘
```

- ドメイン層はインフラストラクチャ層を知らない（インターフェースのみ定義）
- インフラストラクチャ層がドメイン層のインターフェースに依存する（依存関係の逆転）
- ドメインサービスはインターフェースに対してプログラムするため、永続化先がDBでもメモリでも差し替え可能

## 利用例

```go
// app/cmd/api/main.go

func main() {
	// インフラストラクチャ層の実装をインターフェース経由で注入
	userRepository := persistence.NewUserPersistence()
	userService := service.NewUserService(userRepository)

	user, err := user.NewUser("uuid", "test")
	if err != nil {
		// エラーハンドリング
	}

	// ドメインサービスがリポジトリを通じて重複確認
	isExist := userService.IsExists(user)
}
```

`main.go` で具体的な実装（`userPersistence`）を生成し、インターフェース（`IUserRepository`）として注入する   
これにより、テスト時にはモックに差し替えるといった柔軟な構成が可能になる

## サンプルコード

| 種類 | パス |
|------|------|
| リポジトリ（インターフェース） | `/app/internal/domain/repository/user_repository.go` |
| リポジトリ（実装） | `/app/internal/infrastructure/persistence/user.go` |
| データモデル | `/app/internal/infrastructure/data_model/user.go` |
| エンティティ | `/app/internal/domain/model/user/user.go` |
| ドメインサービス | `/app/internal/domain/service/user_service.go` |
| 実行サンプル | `/app/cmd/api/main.go` |
