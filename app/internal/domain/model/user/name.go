package user

import validation "github.com/go-ozzo/ozzo-validation/v4"

type Name string

func NewName(name string) (Name, error) {
	if err := validateName(name); err != nil {
		return "", err
	}

	v := Name(name)

	return v, nil
}

func validateName(name string) error {
	return validation.Validate(name,
		validation.RuneLength(3, 30).Error("名前は3~30文字以内で指定してください"),
	)
}

func (n Name) String() string {
	return string(n)
}
