# エンティティ

### エンティティとは
同一性によって識別するドメインモデルのドメインオブジェクト

### 性質
- 可変である
- 同じ属性でも区別される
- 同一性により区別される

## 可変である とは
エンティティは、属性が変化しても同じオブジェクトとして扱われるという性質を持ちます。例えば、人の名前が変わったとしても、その人自体が別の存在になるわけではない   
これは、値オブジェクトとは異なり、「同一性」を持つことを意味する   

↓以下はエンティティの属性を可変させる例
```go
// リスト1-1 
package user

import validation "github.com/go-ozzo/ozzo-validation/v4"

type User struct {
	name string
}

func NewUser(name string) (*User, error) {
	user := new(User)
	if err := user.ChangeName(name); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangeName は値が正常か確認しUserエンティティのname属性を変更します
func (u *User) ChangeName(name string) error {
	if err := validation.Validate(name,
		validation.RuneLength(3, 10).Error("名前は3~10文字以内で指定してください"),
	); err != nil {
		return err
	}

	u.name = name

	return nil
}
```
処理が内包されているセッターなどと違い、メソッド名を見ればその振る舞いが語られているので視認性も高い点がポイント   
値オブジェクトとの違いは交換(代入)によって変更を表現せず、メソッド名に振る舞いが定義されているものを通じて変化させる   
値オブジェクトと同様なのは異常値のチェックを行うので、不完全なオブジェクトが生成されることはない   

## 同じ属性でも区別される とは
山田 太郎さんという名前の人が二人いた場合、同じ「山田 太郎」という名前だけで識別すると、システム上は同一人物と見なされてしまう
しかし、エンティティは「同じ属性を持っていても異なる存在である」ことを保証するため、識別子（ID）によって区別される   

## 同一性により区別される とは
あるユーザーが名前を変更したとしても、そのユーザーはシステム上では同じ人物   
これは、エンティティが「名前」ではなく「識別子（ID）」を持つことで同一性を維持しているから   
同じ名前の人が複数いても識別子によって区別できるため、エンティティはシステム上で正しく扱われる   

↓以下は同一性を判断するために識別子(identify)を追加した例
```go
// リスト1-1 
package user

import validation "github.com/go-ozzo/ozzo-validation/v4"

type User struct {
    userID UserID  // 追加
	name   string
}

func NewUser(userID, name string) (*User, error) {
	user := new(User)

    user.userID = userID // 追加

	if err := user.ChangeName(name); err != nil {
		return nil, err
	}

	return user, nil
}

// ChangeName は値が正常か確認しUserエンティティのname属性を変更します
func (u *User) ChangeName(name string) error {
	if err := validation.Validate(name,
		validation.RuneLength(3, 10).Error("名前は3~10文字以内で指定してください"),
	); err != nil {
		return err
	}

	u.name = name

	return nil
}
```

# 値オブジェクトとエンティティの違い
- エンティティと値オブジェクトの違いは、「ライフサイクルを持つかどうか」によって判断できる   
例えば、ユーザーは登録・更新・退会といったライフサイクルがあるためエンティティと考える   
一方、氏名や住所のように、値そのものが重要でライフサイクルを持たないものは値オブジェクトとして扱われる   
ユーザーを例にすれば以下のライフサイクルがある
```
新規登録
↓
更新
↓
退会
```
上記のようなライフサイクルがある場合はエンティティとして定義する   
- エンティティの比較処理は識別子のみでOK   
値オブジェクトの場合は全てのプロパティを比較し等価性を確認していたが、エンティティの場合等価保証するのは識別子になる
```go
func (u *User) Equals(other *User) bool {
	// エンティティは同一性だけの比較で良い
	return reflect.DeepEqual(u.userID, other.userID)
}
```

# サンプルコード
値オブジェクト → /app/internal/domain/model/*
実行サンプル → /app/cmd/main.go
