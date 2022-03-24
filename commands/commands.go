package commands

type UserCommands struct {
	CreateUser CreateUserCmdHandler
	UpdateUser UpdateUserCmdHandler
	DeleteUser DeleteUserCmdHandler
}

func NewUserCommands(createUser CreateUserCmdHandler,
	updateUser UpdateUserCmdHandler,
	deleteUser DeleteUserCmdHandler) *UserCommands {
	return &UserCommands{
		createUser,
		updateUser,
		deleteUser,
	}
}

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
