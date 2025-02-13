package user

import "reflect"

type User struct {
	userID UserID
	name   Name
}

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

func (u *User) ChangeName(name string) error {
	v, err := newName(name)
	if err != nil {
		return err
	}

	u.name = v

	return nil
}

func (u *User) Equals(other *User) bool {
	// エンティティは同一性だけの比較で良い
	return reflect.DeepEqual(u.userID, other.userID)
}
