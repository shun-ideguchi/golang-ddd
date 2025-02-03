package fullname

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type FullName struct {
	firstName string
	lastName  string
}

func NewFullName(firstName, lastName string) (*FullName, error) {
	if err := validateFirstName(firstName); err != nil {
		return nil, err
	}
	if err := validateLastName(lastName); err != nil {
		return nil, err
	}

	return &FullName{
		firstName: firstName,
		lastName:  lastName,
	}, nil
}

func validateFirstName(firstName string) error {
	return validation.Validate(firstName,
		validation.RuneLength(3, 10).Error("姓は3~10文字以内で指定してください"),
	)
}

func validateLastName(lastName string) error {
	return validation.Validate(lastName,
		validation.RuneLength(3, 10).Error("名は3~10文字以内で指定してください"),
	)
}
