package user

import (
	"reflect"
)

type User struct {
	userID UserID
	name   Name
	email  Email
}

func NewUser(userID, name, email string) (*User, error) {
	newUserID, err := NewUserID(userID)
	if err != nil {
		return nil, err
	}
	newName, err := NewName(name)
	if err != nil {
		return nil, err
	}
	newEmail, err := NewEmail(email)
	if err != nil {
		return nil, err
	}

	return &User{
		userID: newUserID,
		name:   newName,
		email:  newEmail,
	}, nil
}

func ReNewUser(ID, name, email string) *User {
	return &User{
		userID: UserID(ID),
		name:   Name(name),
		email:  Email(email),
	}
}

func (u *User) ChangeName(name Name) {
	u.name = name
}

func (u *User) ChangeEmail(email Email) {
	u.email = email
}

func (u *User) Equals(other *User) bool {
	// エンティティは同一性だけの比較で良い
	return reflect.DeepEqual(u.userID, other.userID)
}

func (u *User) ID() *UserID {
	return &u.userID
}

func (u *User) Name() *Name {
	return &u.name
}

func (u *User) Email() *Email {
	return &u.email
}
