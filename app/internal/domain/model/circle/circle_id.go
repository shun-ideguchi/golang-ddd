package circle

import validation "github.com/go-ozzo/ozzo-validation/v4"

type CircleID string

func NewCircleID(circleID string) (CircleID, error) {
	if err := validateCircleID(circleID); err != nil {
		return "", err
	}

	v := CircleID(circleID)

	return v, nil
}

func validateCircleID(circleID string) error {
	return validation.Validate(circleID,
		validation.RuneLength(1, 36).Error("サークルIDは1~36文字以内で指定してください"),
	)
}

func (c CircleID) String() string {
	return string(c)
}
