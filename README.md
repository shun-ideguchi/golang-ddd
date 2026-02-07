# ドメインサービス

### ドメインサービスとは
値オブジェクトやエンティティに定義すると不自然になる振る舞いを記述するためのドメインオブジェクト

### なぜドメインサービスが必要か
ドメインには、特定のエンティティや値オブジェクトに持たせると不自然になる振る舞いが存在する   
例えば「ユーザーの重複確認」をUserエンティティに定義すると以下のようになる

```go
// Userエンティティに重複確認を持たせた場合
func (u *User) IsExists(other *User) bool {
	// 自分自身に「他のユーザーが存在するか」を聞いている — 不自然
	return u.name == other.name
}
```

このコードは「ユーザーに対して、他のユーザーが重複しているか聞いている」という不自然な構造になる   
重複確認は特定のユーザー自身の振る舞いではなく、ユーザー全体を見渡して判断する処理であるため、エンティティに持たせるべきではない   
このような場合にドメインサービスを利用する

## ドメインサービスの実装

```go
// app/internal/domain/service/user_service.go

package service

import "github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"

type userService struct{}

func NewUserService() *userService {
	return &userService{}
}

func (s *userService) IsExists(user *user.User) bool {
	// 重複を確認する処理
	return true
}
```

ポイント:
- ドメインサービスは状態を持たない。`userService` 構造体にフィールドがないのは、ドメインサービスがデータではなく振る舞いのみを提供するため
- メソッド名にドメインの意図が表現されている（`IsExists` = ユーザーが既に存在するか）

## メソッドの引数にはドメインオブジェクトを指定する

ドメインサービスのメソッド引数には、具体的な値（名前やメールアドレス）ではなくドメインオブジェクトを指定する

↓以下は具体的な値を引数にした場合

```go
func (s *userService) IsExists(name string) bool {
	// 名前で重複確認
}
```

この設計では、重複確認の条件が「名前 + メールアドレス」に変更された場合、メソッドのシグネチャを変更しなければならない

```go
// 要件変更: メールアドレスでも重複を確認したい
func (s *userService) IsExists(name, email string) bool {
	// シグネチャが変わる → 呼び出し側も全て修正が必要
}
```

↓以下はドメインオブジェクトを引数にした場合

```go
func (s *userService) IsExists(user *user.User) bool {
	// ユーザーエンティティから必要な値を取得して重複確認
}
```

ドメインオブジェクトを引数にすることで、重複確認の条件が変わってもシグネチャは変わらない   
メソッド内部でエンティティから必要な値を取得するだけで対応できるため、呼び出し側への影響がない

## ドメインサービスの乱用に注意する

ドメインサービスは便利だが、何でもドメインサービスに定義してしまうと**ドメインモデル貧血症**に陥る

### ドメインモデル貧血症とは

ドメインオブジェクト（エンティティ・値オブジェクト）にロジックがほとんどなく、データの入れ物だけになってしまう状態

↓以下はドメインサービスを乱用した場合

```go
// ユーザー名変更サービス — 本来エンティティに定義すべき振る舞い
type userService struct{}

func (s *userService) ChangeName(user *user.User, name string) error {
	// ユーザー名の変更をサービスが行う
}

// ユーザーメール変更サービス
func (s *userService) ChangeEmail(user *user.User, email string) error {
	// メールアドレスの変更もサービスが行う
}
```

```go
// Userエンティティ — データだけで振る舞いが何もない（貧血症）
type User struct {
	userID UserID
	name   Name
}
// メソッドが一つもない…
```

このように、エンティティが持つべき振る舞い（名前の変更など）までドメインサービスに移してしまうと、エンティティはただのデータの入れ物になる   
ドメインモデルを見てもドメインの知識が読み取れず、DDDの利点が失われてしまう

### 正しい設計

ユーザー名の変更はユーザー自身の振る舞いであるため、エンティティに定義する

```go
// app/internal/domain/model/user/user.go

// ChangeName はUserエンティティ自身の振る舞い — エンティティに定義するのが自然
func (u *User) ChangeName(name string) error {
	v, err := newName(name)
	if err != nil {
		return err
	}
	u.name = v
	return nil
}
```

- `ChangeName` はユーザー自身の属性を変更する振る舞いであるため、エンティティに定義するのが自然
- バリデーション（`newName`）もエンティティ内で完結しており、不正な値への変更を防げる

## ドメインサービスに定義すべきかの判断基準

ドメインサービスに定義すべきか迷った場合は、**まずエンティティや値オブジェクトに定義する**ことを検討する

| 判断基準 | 定義先 | 例 |
|----------|--------|------|
| 特定のオブジェクト自身の振る舞いか | エンティティ / 値オブジェクト | ユーザー名の変更、メールアドレスのフォーマット検証 |
| 複数のオブジェクトにまたがる判断か | ドメインサービス | ユーザーの重複確認 |
| 特定のオブジェクトに持たせると不自然か | ドメインサービス | 「ユーザーに自分が重複しているか聞く」のは不自然 |

**原則: 可能な限りドメインサービスを避ける**   
ドメインサービスは最終手段であり、エンティティや値オブジェクトに自然に定義できないものだけを置く

## 利用例

```go
// app/cmd/api/main.go

func main() {
	userService := service.NewUserService()

	user, err := user.NewUser("uuid", "test")
	if err != nil {
		fmt.Println("failed to initialize user model")
	}

	// ドメインサービスで重複確認
	isExist := userService.IsExists(user)
	if isExist {
		fmt.Println("duplicate user")
	}
}
```

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
| ドメインサービス | `/app/internal/domain/service/user_service.go` |
| エンティティ | `/app/internal/domain/model/user/user.go` |
| 実行サンプル | `/app/cmd/api/main.go` |
