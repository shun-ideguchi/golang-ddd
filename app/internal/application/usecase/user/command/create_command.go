package command

type CreateCommand struct {
	UserID string
	Name   string
	Email  string
}

func NewCreateCommand(name, email string) *CreateCommand {
	return &CreateCommand{
		Name:  name,
		Email: email,
	}
}
