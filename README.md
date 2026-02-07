# アプリケーションサービス
### アプリケーションサービスとは
ドメインオブジェクトが行うタスクの進行を管理、まとめあげ問題を解決するもの

### 具体例

#### 例1: ユーザー登録

```go
func (s *UserApplicationService) Register(name string, email string) error {
    // 1. ドメインオブジェクト（Value Object）を生成
    userName, err := domain.NewUserName(name)
    // 2. ドメインオブジェクト（Entity）を生成
    user, err := domain.NewUser(userName, email)
    // 3. ドメインサービスで重複チェック
    exists, err := s.userDomainService.IsExists(user)
    // 4. リポジトリに保存
    err = s.userRepository.Save(user)
    return nil
}
```

#### 例2: 注文処理

```go
func (s *OrderApplicationService) PlaceOrder(customerID string, items []OrderItemDTO) error {
    // 1. 顧客を取得（リポジトリ）
    customer, err := s.customerRepository.FindByID(customerID)
    // 2. 注文エンティティを生成（ドメインオブジェクト）
    order, err := domain.NewOrder(customer, items)
    // 3. 在庫チェック（ドメインサービス）
    err = s.inventoryService.CheckAvailability(order)
    // 4. 注文を保存（リポジトリ）
    err = s.orderRepository.Save(order)
    // 5. イベント発行（通知など）
    s.eventPublisher.Publish(domain.OrderPlacedEvent{OrderID: order.ID()})
    return nil
}
```

### ドメインオブジェクトの公開について(DTOの採用)
ユーザーの取得処理などにおいて、結果を返却する場合の注意点を示す。

#### 例1: 戻り値としてドメインオブジェクトを公開
```go
func (s *UserApplicationService) Get(id string) (domain.User, error) {
    userId, err := domain.NewUserId(id)
    user, err := s.userRepository.Find(userId)

    return user, nil
}
```
```go
func (c *Client) ChangeName(id, name string) error {
    user, err := c.UserApplicationService.Get(id)
    newName, err := domain.NewUserName(name)
    // 自由にドメインオブジェクトを操作できてしまう
    user.ChangeName(newName)
}
```
アプリケーションサービス以外のオブジェクトがドメインオブジェクトの直接のクライアントとなり自由に操作できてしまう問題が発生する。
その結果、ドメインオブジェクトのふるまいを呼び出せるのはアプリケーションサービスの役割であるにもかかわらず、他のクライアント各所に散りばめられることでドメインオブジェクトに対する多くの依存を生み出す。

#### 例2: DTOの採用
```go
type UserData struct {
    Id   string
    Name string
}

func NewUserData(user domain.User) UserData {
    return UserData{
        Id:   user.Id(),
        Name: user.Name(),
    }
}
```
```go
func (s *UserApplicationService) Get(id string) (userData UserData, err error) {
    userId, err := domain.NewUserId(id)
    user, err := s.userRepository.Find(userId)

    userData = NewUserData(user)
    return
}
```
DTOを採用することでクライアントからドメインオブジェクトのふるまいを呼び出すことはできなくなる。
なお、DTOのコンストラクタの引数にドメインオブジェクトを指定しているのは、外部に公開するパラメータが増えてもコンストラクタのシグネチャが変更されないことで修正範囲の肥大化を防いでいる。
DTOはドメインオブジェクトと密な関係であることから依存することに問題はないと考えられている。

### コマンドオブジェクトの採用（ユーザー更新処理）
アプリケーションサービスのメソッド引数で `name`, `email` などを一つずつ受け取ると、ユーザーのプロパティが増えるごとにメソッドのシグネチャが変更されてしまう。
コマンドオブジェクトをシグネチャに設定することで、この問題を解決できる。

#### 問題: 引数を直接受け取る場合
```go
// プロパティが増えるたびにシグネチャが変わる
func (s *UserApplicationService) Update(id string, name string, mailAddress string) error { ... }
// ↓ 電話番号が増えた場合、シグネチャが変更される
func (s *UserApplicationService) Update(id string, name string, mailAddress string, phone string) error { ... }
```
呼び出し元すべてに修正が波及してしまう。

#### 解決: コマンドオブジェクトの導入
```go
type UserUpdateCommand struct {
    Id          string
    Name        *string // nilなら「更新しない」を表す
    MailAddress *string
}
```
```go
// プロパティが増えてもシグネチャは変わらない
func (s *UserApplicationService) Update(command UserUpdateCommand) error {
    userId, err := domain.NewUserId(command.Id)
    user, err := s.userRepository.Find(userId)

    // nilでなければ名前を更新
    if command.Name != nil {
        name, err := domain.NewUserName(*command.Name)
        user.ChangeName(name)
    }

    // nilでなければメールアドレスを更新
    if command.MailAddress != nil {
        mail, err := domain.NewMailAddress(*command.MailAddress)
        user.ChangeMailAddress(mail)
    }

    err = s.userRepository.Save(user)
    return nil
}
```

#### 呼び出し側
```go
// 名前だけ更新
name := "naruse"
updateNameCommand := UserUpdateCommand{
    Id:   id,
    Name: &name,
}
userApplicationService.Update(updateNameCommand)

// メールアドレスだけ更新
mail := "aaa@e.com"
updateMailAddressCommand := UserUpdateCommand{
    Id:          id,
    MailAddress: &mail,
}
userApplicationService.Update(updateMailAddressCommand)
```

#### コマンドオブジェクトのポイント

| 観点 | 説明 |
|---|---|
| **`nil` の意味** | 「そのフィールドは更新しない」ことを表す |
| **シグネチャの安定** | プロパティ追加時にメソッドシグネチャが変わらない |
| **部分更新が可能** | 変更したいフィールドだけセットすればよい |
| **呼び出し元への影響を最小化** | 新フィールド追加時も既存コードの修正が不要 |

なお、Goではオプショナルな値を表現するために `*string` のようなポインタ型を使う。`nil` であれば「未指定（更新しない）」、値が入っていれば「更新する」という判定ができる。

### アプリケーションサービスのポイント

| 特徴 | 説明 |
|---|---|
| **ビジネスロジックを持たない** | ロジックはドメインオブジェクトに任せる |
| **タスクの調整役** | 複数のドメインオブジェクトやリポジトリを組み合わせる |
| **ユースケースに対応** | 「ユーザーを登録する」「注文を確定する」などの操作単位 |
| **トランザクション管理** | DB のトランザクション境界を定めることが多い |
| **入出力の変換** | DTO（外部向けデータ）⇔ ドメインオブジェクトの変換を行う |

アプリケーションサービスは **「オーケストレーター（指揮者）」** のような存在で、ドメインオブジェクトたちに「何をするか」を指示し、その進行を管理してユースケースを完遂させる役割を担っている。
