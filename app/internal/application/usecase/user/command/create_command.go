package command

import (
	"github.com/google/uuid"
)

type CreateCommand struct {
	UserID string
	Name   string
	Email  string
}

func NewCreateCommand(name, email string) *CreateCommand {
	return &CreateCommand{
		UserID: uuid.NewString(),
		Name:   name,
		Email:  email,
	}
}
