package user

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Email string

func NewEmail(email string) (Email, error) {
	if err := validateEmail(email); err != nil {
		return "", err
	}

	v := Email(email)

	return v, nil
}

func validateEmail(email string) error {
	return validation.Validate(email,
		is.EmailFormat.Error("メールアドレスの形式が正しくありません"),
	)
}

func (e Email) String() string {
	return string(e)
}
