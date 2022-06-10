package commands

import "github.com/novabankapp/usermanagement.application/dtos"

type UserCommands struct {
	UpdateUser UpdateUserCmdHandler
	DeleteUser DeleteUserCmdHandler
}

func NewUserCommands(
	updateUser UpdateUserCmdHandler,
	deleteUser DeleteUserCmdHandler) *UserCommands {
	return &UserCommands{
		updateUser,
		deleteUser,
	}
}

type CreateUserCommand struct {
	dto dtos.CreateUserDto
}

type UpdateUserCommand struct {
	dto dtos.UpdateUserDto
}

func NewUpdateUserCommand(dto dtos.UpdateUserDto) *UpdateUserCommand {
	return &UpdateUserCommand{
		dto,
	}
}

type DeleteUserCommand struct {
	dto dtos.DeleteUserDto
}

func NewDeleteUserCommand(dto dtos.DeleteUserDto) *DeleteUserCommand {
	return &DeleteUserCommand{
		dto,
	}
}
