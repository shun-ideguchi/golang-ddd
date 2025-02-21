package user

import validation "github.com/go-ozzo/ozzo-validation/v4"

type UserID string

func newUserID(userID string) (UserID, error) {
	if err := validateUserID(userID); err != nil {
		return "", err
	}

	v := UserID(userID)

	return v, nil
}

func validateUserID(userID string) error {
	return validation.Validate(userID,
		validation.RuneLength(1, 31).Error("ユーザーIDは1~31文字以内で指定してください"),
	)
}

func (u UserID) String() string {
	return string(u)
}
