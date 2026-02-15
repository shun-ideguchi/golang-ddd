package repository

import "github.com/shun-ideguchi/golang-ddd/internal/domain/model/circle"

type ICircleRepository interface {
	FindByCircleID(circleID circle.CircleID) (*circle.Circle, error)
	FindByCircleName(circleName circle.CircleName) (*circle.Circle, error)
	Save(circle *circle.Circle) error
}
