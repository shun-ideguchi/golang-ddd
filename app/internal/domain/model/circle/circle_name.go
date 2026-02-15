package circle

import (
	"reflect"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CircleName string

func NewCircleName(circleName string) (CircleName, error) {
	if err := validateCircleName(circleName); err != nil {
		return "", err
	}

	v := CircleName(circleName)

	return v, nil
}

func validateCircleName(circleName string) error {
	return validation.Validate(circleName,
		validation.RuneLength(3, 20).Error("サークル名は3~20文字以内で指定してください"),
	)
}

func (c CircleName) String() string {
	return string(c)
}

func (c CircleName) Equals(other CircleName) bool {
	return reflect.DeepEqual(c, other)
}
