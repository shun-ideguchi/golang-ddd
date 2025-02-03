# 値オブジェクト

### 性質
値のオブジェクトの性質は以下のように定義されている
- 不変である
- 交換が可能である
- 透過性によって比較される

## 不変である とは
オブジェクトが生成された時点で値が確定し、セッターなどで途中更新されない   
値オブジェクトはコンストラクタなどで生成されることが一般的

↓以下は変更されているとはならない（変更ではなく代入）
```go
// リスト1
text := "こんにちは"
fmt.Println(text)   // こんにちは
text = "さようなら"
fmt.Println(text)   // さようなら
```

↓以下は変更されていると言える（text変数にChangeメソッドは存在するはずがないが疑似的に値を変更するを再現した処理）
```go
// リスト2
text := "こんにちは"
text.Change("さようなら")
fmt.Println(text)   // さようなら？
```

↓リスト2をより分かりやすく値を変更するを再現した処理
```go
// リスト3
"こんにちは".Change("さようなら")
fmt.Println("こんにちは")   // さようならにはならない
```

↓実際の現場で見られる実装（リスト3と同じである）
```go
fullName := fullname.NewFullName("山田", "太郎")
fullName.Change("花子")
```
