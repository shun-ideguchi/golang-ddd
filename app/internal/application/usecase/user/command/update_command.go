package command

type UpdateCommand struct {
	UserID string
	Name   *string
	Email  *string
}

func NewUpdateCommand(userID string, options ...func(*UpdateCommand)) *UpdateCommand {
	command := &UpdateCommand{
		UserID: userID,
	}

	for _, option := range options {
		option(command)
	}

	return command
}

func WithName(name string) func(*UpdateCommand) {
	return func(command *UpdateCommand) {
		command.Name = &name
	}
}

func WithEmail(email string) func(*UpdateCommand) {
	return func(command *UpdateCommand) {
		command.Email = &email
	}
}
