package money

import (
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type Money struct {
	amount   int
	currency string
}

func NewMoney(amount int, currency string) (*Money, error) {
	if err := validateCurrency(currency); err != nil {
		return nil, err
	}

	return &Money{
		amount:   amount,
		currency: currency,
	}, nil
}

func validateCurrency(currency string) error {
	return validation.Validate(currency,
		validation.In("JPY", "USD").Error("通貨単位が正しくありません"),
	)
}

// Add は金額追加します
func (m *Money) Add(arg Money) (*Money, error) {
	if m.currency != arg.currency {
		return nil, fmt.Errorf("通貨単位が異なります")
	}

	return NewMoney((m.amount + arg.amount), m.currency)
}
