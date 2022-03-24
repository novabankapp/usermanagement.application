package commands

type CreateUserCommand struct {
}

func NewCreateUserCommand() *CreateUserCommand {
	return &CreateUserCommand{}
}

type UpdateUserCommand struct {
}

func NewUpdateUserCommand() *UpdateUserCommand {
	return &UpdateUserCommand{}
}

type DeleteUserCommand struct {
}

func NewDeleteUserCommand() *DeleteUserCommand {
	return &DeleteUserCommand{}
}
