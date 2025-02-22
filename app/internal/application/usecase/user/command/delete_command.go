package command

type DeleteCommand struct {
	UserID string
}

func NewDeleteCommand(userID string) *DeleteCommand {
	return &DeleteCommand{UserID: userID}
}
