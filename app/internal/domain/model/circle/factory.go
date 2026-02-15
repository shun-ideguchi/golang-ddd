package circle

import "github.com/shun-ideguchi/golang-ddd/internal/domain/model/user"

type ICircleFactory interface {
	Create(circleName CircleName, owner *user.User) (*Circle, error)
}
