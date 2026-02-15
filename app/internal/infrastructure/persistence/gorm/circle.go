package gorm

import (
	"fmt"

	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/circle"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
)

type circlePersistence struct {
}

func NewCirclePersistence() repository.ICircleRepository {
	return &circlePersistence{}
}

func (p *circlePersistence) FindByCircleID(circleID circle.CircleID) (*circle.Circle, error) {
	return nil, nil
}

func (p *circlePersistence) FindByCircleName(circleName circle.CircleName) (*circle.Circle, error) {
	return nil, nil
}

func (p *circlePersistence) Save(circle *circle.Circle) error {
	builder := &CircleDataModelBuilder{}
	circle.Notify(builder) // ← domainに「中身教えて」と依頼
	m := builder.Build()   // ← DBモデルに変換

	fmt.Println(m)
	return nil
}
