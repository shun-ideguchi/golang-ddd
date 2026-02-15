package user

type UserNotification interface {
	ID(UserID) UserNotification
	Name(Name) UserNotification
	Email(Email) UserNotification
}
