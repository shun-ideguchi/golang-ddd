package service

import (
	"github.com/shun-ideguchi/golang-ddd/internal/domain/model/circle"
	"github.com/shun-ideguchi/golang-ddd/internal/domain/repository"
)

type CircleService struct {
	circleRepository repository.ICircleRepository
}

func NewCircleService(circleRepository repository.ICircleRepository) *CircleService {
	return &CircleService{circleRepository: circleRepository}
}

func (s *CircleService) IsExists(circle *circle.Circle) bool {
	found, err := s.circleRepository.FindByCircleName(circle.CircleName())
	if err != nil {
		// errorを返すべきだが省略
		return false
	}

	if found == nil {
		return false
	}

	return true
}
