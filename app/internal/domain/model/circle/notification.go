package circle

import "github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"

type CircleNotification interface {
	ID(CircleID) CircleNotification
	Name(CircleName) CircleNotification
	Owner(user.User) CircleNotification
	Members([]user.User) CircleNotification
}
