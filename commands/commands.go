package commands

import "github.com/novabankapp/usermanagement.application/dtos"

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
	dto *dtos.CreateUserDto
}

func NewCreateUserCommand(dto *dtos.CreateUserDto) *CreateUserCommand {
	return &CreateUserCommand{
		dto,
	}
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
