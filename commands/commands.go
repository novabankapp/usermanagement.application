package commands

import (
	"github.com/novabankapp/usermanagement.application/dtos"
	registrationDtos "github.com/novabankapp/usermanagement.application/dtos/registration"
)

type UserCommands struct {
	CreateUser CreateUserCmdHandler
	UpdateUser UpdateUserCmdHandler
	DeleteUser DeleteUserCmdHandler
}

func NewUserCommands(
	createUser CreateUserCmdHandler,
	updateUser UpdateUserCmdHandler,
	deleteUser DeleteUserCmdHandler) *UserCommands {
	return &UserCommands{
		createUser,
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

func NewCreateUserCommand(dto dtos.CreateUserDto) *CreateUserCommand {
	return &CreateUserCommand{
		dto,
	}
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

type RegisterUserCommand struct {
	Dto registrationDtos.RegisterUserDto
}

func NewRegisterUserCommand(dto registrationDtos.RegisterUserDto) *RegisterUserCommand {
	return &RegisterUserCommand{
		dto,
	}
}
