package user

type IFactory interface {
	Create(name, email string) (*User, error)
}
