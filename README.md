# ファクトリー

### ファクトリーとは
オブジェクトの生成処理を専門に担うドメインオブジェクト   
生成ロジックが複雑な場合に、その知識をファクトリーに集約させることで呼び出し側の負担を軽減する

### なぜファクトリーが必要か
例えば、ユーザーの新規作成時にUUIDの採番が必要だとする   
この採番ロジックをユースケース層に書いてしまうと、以下の問題が起きる
- ID生成という技術的関心がユースケースに漏れ出す
- 複数箇所でユーザーを生成する場合、採番ロジックが散らばる
- 採番方式の変更時に複数箇所を修正する必要がある

ファクトリーに生成ロジックを集約することで、呼び出し側は「名前」や「メールアドレス」などドメインに関連する値だけを渡せばよくなる

## なぜDBの自動採番（AUTO_INCREMENT）ではなくアプリケーションで採番するのか

MySQLなどのRDBには自動採番（AUTO_INCREMENT）機能があるが、エンティティの識別子にこれを採用するとオブジェクトが不安定になる

### 永続化するまでIDが存在しない

自動採番はINSERT時にDBが割り当てるため、インスタンス生成時点ではIDが存在しない   
エンティティは識別子によって識別されるオブジェクトであるにもかかわらず、生成直後は「自分が何者か分からない」不完全な状態になる

```go
// 自動採番を採用した場合のイメージ
user := &User{
	userID: "",           // IDがまだない — エンティティとして不完全
	name:   "山田 太郎",
}

// DBにINSERTして初めてIDが確定する
db.Create(user)  // ここでようやく user.userID = 1 になる
```

この設計では、永続化前のエンティティは識別子を持たない「不完全なオブジェクト」として存在してしまう

### エンティティの同一性比較ができない

`Equals` メソッドは識別子で同一性を判断するが、永続化前のエンティティはIDが空のため正しく比較できない

```go
func (u *User) Equals(other *User) bool {
	return reflect.DeepEqual(u.userID, other.userID)
}

// 永続化前のエンティティ同士を比較すると…
user1 := &User{userID: ""}  // IDなし
user2 := &User{userID: ""}  // IDなし
user1.Equals(user2) // true — 別のオブジェクトなのに同一と判定されてしまう
```

### ドメインロジックがインフラストラクチャ層に依存する

IDを得るためにまず永続化が必要になると、ドメインロジックの実行順序が永続化層に縛られる   
例えば、生成後すぐにドメインイベントの発行やログ記録でIDを参照したくても、先にDBへ保存しなければならない   
これはドメイン層がインフラストラクチャ層に依存することを意味し、DDDの依存関係の原則に反する

### テストが困難になる

自動採番はDBに依存するため、ユニットテストでエンティティを生成するたびにDBアクセスが必要になる   
テストの実行速度が低下し、DB環境がなければテストできないという制約が生まれる

### 本リポジトリでの解決策

本リポジトリではファクトリー内でUUIDを採番することで、上記の問題を解決している

```go
func (f *userFactory) Create(name, email string) (*user.User, error) {
	userID := uuid.NewString()  // 永続化に依存せずIDを生成
	return user.NewUser(userID, name, email)
}
```

- インスタンス生成時点で識別子が確定するため、エンティティは常に完全な状態で存在する
- 永続化前でも `Equals` による同一性比較が正しく動作する
- ドメインロジックが永続化のタイミングに依存しない
- DBなしでテスト可能

## ファクトリーのインターフェース

ファクトリーのインターフェースはエンティティと同じパッケージに定義する   
これにより、エンティティの生成方法がドメイン層の中で表現される

```go
// app/internal/domain/model/user/factory.go

package user

type IFactory interface {
	Create(name, email string) (*User, error)
}
```

ポイント:
- 引数にはユーザーが入力する値（`name`, `email`）のみを受け取る。識別子（`userID`）は含めない
- 識別子の生成はファクトリーの実装に委ねることで、呼び出し側はID生成の詳細を知る必要がない
- 戻り値はドメインモデル（`*User`）を返す

## ファクトリーの実装

インターフェースの実装はエンティティとは別のパッケージに配置する   
本リポジトリではUUID生成ライブラリを使用しているため、エンティティパッケージから分離している

```go
// app/internal/domain/model/factory/user.go

package factory

import (
	"github.com/google/uuid"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"
)

type userFactory struct{}

func NewUserFactory() user.IFactory {
	return &userFactory{}
}

func (f *userFactory) Create(name, email string) (*user.User, error) {
	// ID生成が複雑なものと仮定
	userID := uuid.NewString()
	return user.NewUser(userID, name, email)
}
```

ポイント:
- `NewUserFactory` の戻り値の型は `user.IFactory`（インターフェース）にする
- `Create` 内でUUIDを採番し、`user.NewUser` を呼び出す。`NewUser` はバリデーションを行うため、不正な値でのエンティティ生成を防げる
- ID生成方式を変更したい場合、このファクトリー実装を差し替えるだけで済む

## ファクトリーを使わない場合との比較

↓以下はファクトリーを使わない場合

```go
func (u *createUsecase) Execute(cmd *command.CreateCommand) error {
	// ユースケースにID生成の知識が漏れ出している
	userID := uuid.NewString()
	user, err := user.NewUser(userID, cmd.Name, cmd.Email)
	if err != nil {
		return err
	}
	// ...
}
```

↓以下はファクトリーを使う場合

```go
func (u *createUsecase) Execute(cmd *command.CreateCommand) error {
	// 名前とメールアドレスだけ渡せばよい
	user, err := u.userFactory.Create(cmd.Name, cmd.Email)
	if err != nil {
		return fmt.Errorf("failed to create user entity: %w", err)
	}
	// ...
}
```

ファクトリーを使うことで、ユースケースは「何を作るか」に集中でき、「どうやって識別子を生成するか」を意識しなくてよくなる

## ユースケースでの利用

ユースケース層ではファクトリーをインターフェース経由で受け取り、エンティティの生成に使用する

```go
// app/internal/application/usecase/user/create_usecase.go

package user

type createUsecase struct {
	userFactory    user.IFactory
	userRepository repository.IUserRepository
	userService    service.UserService
}

func NewCreateUsecase(
	userFactory user.IFactory,
	userRepository repository.IUserRepository,
	userService service.UserService,
) *createUsecase {
	return &createUsecase{
		userFactory:    userFactory,
		userRepository: userRepository,
		userService:    userService,
	}
}

func (u *createUsecase) Execute(cmd *command.CreateCommand) error {
	// ファクトリーでエンティティを生成
	user, err := u.userFactory.Create(cmd.Name, cmd.Email)
	if err != nil {
		return fmt.Errorf("failed to create user entity: %w", err)
	}

	// ドメインサービスで重複チェック
	if u.userService.IsExists(user) {
		return fmt.Errorf("duplicate user: %s", cmd.Name)
	}

	// リポジトリで永続化
	return u.userRepository.Save(user)
}
```

ユースケースはコマンド（`CreateCommand`）を受け取り、以下の流れで処理する
1. ファクトリーでエンティティを生成
2. ドメインサービスで重複チェック
3. リポジトリで永続化

### コマンドオブジェクト

ユースケースへの入力値はコマンドオブジェクトとして定義する   
プレゼンテーション層からの入力をユースケースの引数として整理する役割を持つ

```go
// app/internal/application/usecase/user/command/create_command.go

package command

type CreateCommand struct {
	UserID string
	Name   string
	Email  string
}

func NewCreateCommand(name, email string) *CreateCommand {
	return &CreateCommand{
		Name:  name,
		Email: email,
	}
}
```

`NewCreateCommand` では `Name` と `Email` のみ受け取る。`UserID` はファクトリーが生成するため、コマンド生成時には不要である

## 依存関係の流れ

```
ドメイン層（domain）
┌──────────────────────────────────────────────┐
│                                              │
│  user パッケージ        factory パッケージ     │
│  ┌────────────────┐    ┌──────────────────┐  │
│  │ User           │    │ userFactory      │  │
│  │ IFactory       │◄───│ （IFactory 実装） │  │
│  │ NewUser()      │    │ uuid.NewString() │  │
│  └────────────────┘    └──────────────────┘  │
│          ▲                      ▲            │
└──────────┼──────────────────────┼────────────┘
           │                      │
           │ 利用                  │ 注入
           │                      │
┌──────────┼──────────────────────┼─────────────┐
│  アプリケーション層（application）              │
│  ┌──────────────────────────────────────────┐ │
│  │ createUsecase                            │ │
│  │  - userFactory.Create(name, email)       │ │
│  └──────────────────────────────────────────┘ │
└───────────────────────────────────────────────┘
```

- ファクトリーのインターフェース（`IFactory`）はエンティティと同じ `user` パッケージに定義
- ファクトリーの実装（`userFactory`）は別パッケージ `factory` に配置し、UUID生成の依存を分離
- ユースケースはインターフェースに対してプログラムするため、ID生成方式を差し替え可能

## 利用例

```go
// app/cmd/api/main.go

func main() {
	userRepository := persistence.NewUserPersistence()
	userService := service.NewUserService(userRepository)
	userFactory := factory.NewUserFactory()
	userCreateUsecase := user.NewCreateUsecase(userFactory, userRepository, *userService)

	// コマンドを作成してユースケースを実行
	createCommand := command.NewCreateCommand("test name", "test@test.com")
	if err := userCreateUsecase.Execute(createCommand); err != nil {
		fmt.Println(err.Error())
	}
}
```

`main.go` でファクトリーの具体的な実装を生成し、インターフェース経由でユースケースに注入する   
テスト時にはUUIDではなく固定値を返すモックファクトリーに差し替えることで、再現性のあるテストが可能になる

## 起動方法

```bash
docker compose up --build
```

コンテナ内で以下を実行する。

```bash
go run ./cmd/api/
```

## サンプルコード

| 種類 | パス |
|------|------|
| ファクトリー（インターフェース） | `/app/internal/domain/model/user/factory.go` |
| ファクトリー（実装） | `/app/internal/domain/model/factory/user.go` |
| ユースケース | `/app/internal/application/usecase/user/create_usecase.go` |
| コマンド | `/app/internal/application/usecase/user/command/create_command.go` |
| エンティティ | `/app/internal/domain/model/user/user.go` |
| 実行サンプル | `/app/cmd/api/main.go` |
