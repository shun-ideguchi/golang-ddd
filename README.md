# 集約 - ドメインのルールを守る

## 漏れだしたルール

### 漏れだしたルールとは
ドメインが本来守るべきルールがユースケース（アプリケーション層）に記述されてしまうことを「漏れだしたルール」と呼ぶ   
漏れだしたルールは仕様変更時に修正箇所が散在し、保守コストを増大させる

### 具体例: サークルへの加入人数制限
サークルには「30人まで加入できる」というドメインルールがある   
このルールをユースケースに直接記述してしまうと以下のような問題が起きる

↓ルールがユースケースに漏れだした悪い例
```go
// サークル加入ユースケース（悪い例）
func (u *joinUsecase) Execute(cmd command.JoinCommand) error {
    // ...（省略: メンバーとサークルの取得処理）

    // ドメインルールがユースケースに漏れだしている
    if len(circle.Members()) >= 29 {
        return fmt.Errorf("circle is full")
    }
    circle.Members = append(circle.Members, member)

    return u.circleRepository.Save(circle)
}
```

この書き方には以下の問題がある
1. サークル加入に関するユースケースが複数存在する場合、全てのユースケースに同じルールを記述する必要がある
2. 仕様変更（例: 上限を30人から50人に変更）が入った場合、全ユースケースを漏れなく修正する必要がある
3. 「漏れなく」という慎重さを求める作業は開発者を疲弊させ、バグの温床になる

### 解決策: ドメインモデルにルールを集約する
ドメインルールはドメインモデル自身に持たせるべきである   
以下は実際のCircleドメインモデルの実装である

```go
// /app/internal/domain/model/circle/circle.go

func (c *Circle) IsFull() bool {
    return c.CountMembers() >= 30
}

func (c *Circle) Join(member *user.User) error {
    if member == nil {
        return fmt.Errorf("member is nil: %s", c.id.String())
    }
    if c.IsFull() {
        return fmt.Errorf("circle is full: %s", c.id.String())
    }
    c.members = append(c.members, *member)
    return nil
}
```

`IsFull()` と `Join()` にルールを集約することで、ユースケースは `circle.Join(member)` を呼ぶだけでよくなる   
ルールの存在を意識する必要がない

```go
// /app/internal/application/usecase/circle/join_usecase.go

func (u *joinUsecase) Execute(cmd command.JoinCommand) error {
    // ...（省略: メンバーとサークルの取得処理）

    // ユースケースはドメインルールを意識せず Join を呼ぶだけ
    if err := circle.Join(member); err != nil {
        return fmt.Errorf("failed to join circle: %w", err)
    }
    return u.circleRepository.Save(circle)
}
```

仕様変更が入っても `IsFull()` の1箇所を修正するだけで全てのユースケースに反映される

## 言葉とコードの齟齬をなくす

### 問題: 30人なのに29?
「サークルは30人まで入れる」というルールがある   
サークルにはオーナー1人 + メンバーが存在する   
もし人数チェックを以下のように書いた場合、コード上に「29」という数字が現れる

```go
// 悪い例: 言葉とコードに齟齬がある
func (c *Circle) IsFull() bool {
    // membersにはオーナーが含まれないため29でチェックしている
    // しかし仕様書の「30人」とコードの「29」に齟齬がある
    return len(c.members) >= 29
}
```

「30人まで入れる」というルールに対してコードが `29` で判定しているため、仕様書やドメインエキスパートの言葉とコードの間に齟齬が生まれる   
この齟齬は後からコードを読む開発者に混乱を与え、バグの原因になりうる

### 解決策: CountMembersメソッドで齟齬を解消する

```go
// /app/internal/domain/model/circle/circle.go

func (c *Circle) CountMembers() int {
    // オーナーを含めた人数を返すことで「30人」と一致させる
    return len(c.members) + 1
}

func (c *Circle) IsFull() bool {
    // 「サークルは30人まで」という言葉がそのままコードに現れる
    return c.CountMembers() >= 30
}
```

`CountMembers()` でオーナーを含めた総人数を返すことにより、`IsFull()` の条件が `>= 30` となる   
ドメインエキスパートの「30人まで」という言葉とコードの数値が一致し、齟齬がなくなる

## 通知オブジェクトによるデータモデル構築

### 課題: 永続化のためにドメインの内部データをどう取得するか
リポジトリでドメインオブジェクトを永続化する際、内部データをデータモデルに変換する必要がある   
しかしドメインオブジェクトにgetterを安易に公開すると、集約のルールが外部から破られるリスクがある

### 解決策: 通知パターン
通知パターンでは以下の3つのコンポーネントで構成される

**1. 通知インタフェース（ドメイン層）**   
ドメイン層に通知用のインタフェースを定義する   
これにより、ドメインが「何を通知するか」だけを定義し、「誰に通知するか」は関知しない

```go
// /app/internal/domain/model/circle/notification.go

type CircleNotification interface {
    ID(CircleID) CircleNotification
    Name(CircleName) CircleNotification
    Owner(user.User) CircleNotification
    Members([]user.User) CircleNotification
}
```

**2. ドメインオブジェクトのNotifyメソッド（ドメイン層）**   
ドメインオブジェクトは通知インタフェースを受け取り、自身の内部データを通知する

```go
// /app/internal/domain/model/circle/circle.go

func (c *Circle) Notify(n CircleNotification) {
    n.ID(c.id).Name(c.name).Owner(c.owner).Members(c.members)
}
```

**3. 通知ビルダー（インフラ層）**   
インフラ層で通知インタフェースを実装し、受け取ったデータからデータモデルを構築する

```go
// /app/internal/infrastructure/persistence/gorm/circle_notification_builder.go

type CircleDataModelBuilder struct {
    id      circle.CircleID
    name    circle.CircleName
    owner   user.User
    members []user.User
}

func (b *CircleDataModelBuilder) ID(id circle.CircleID) circle.CircleNotification {
    b.id = id
    return b
}

func (b *CircleDataModelBuilder) Name(name circle.CircleName) circle.CircleNotification {
    b.name = name
    return b
}

func (b *CircleDataModelBuilder) Owner(owner user.User) circle.CircleNotification {
    b.owner = owner
    return b
}

func (b *CircleDataModelBuilder) Members(members []user.User) circle.CircleNotification {
    b.members = members
    return b
}

func (b *CircleDataModelBuilder) Build() *Circle {
    members := make([]User, len(b.members))
    for i, member := range b.members {
        members[i] = User{
            ID:    member.ID().String(),
            Name:  member.Name().String(),
            Email: member.Email().String(),
        }
    }
    return &Circle{
        ID:      b.id.String(),
        Name:    b.name.String(),
        Owner:   b.owner.ID().String(),
        Members: members,
    }
}
```

### リポジトリでの使用例
リポジトリではビルダーを生成し、ドメインオブジェクトにNotifyで通知を依頼するだけでデータモデルが構築できる

```go
// /app/internal/infrastructure/persistence/gorm/circle.go

func (p *circlePersistence) Save(circle *circle.Circle) error {
    builder := &CircleDataModelBuilder{}
    circle.Notify(builder) // ← domainに「中身教えて」と依頼
    m := builder.Build()   // ← DBモデルに変換

    fmt.Println(m)
    return nil
}
```

この通知パターンにより以下のメリットがある
- ドメインオブジェクトが永続化用のデータモデルの存在を知らなくてよい
- ドメインオブジェクトにgetterを公開する必要がない（集約のルールが守られる）
- インフラ層の変更（DB変更など）がドメイン層に影響しない

## サンプルコード
ドメインモデル → /app/internal/domain/model/*   
ユースケース → /app/internal/application/usecase/*   
リポジトリ → /app/internal/infrastructure/persistence/gorm/*   
実行サンプル → /app/cmd/api/main.go
