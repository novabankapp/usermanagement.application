package registration_commands

import (
	registration_dtos "github.com/novabankapp/usermanagement.application/dtos/registration"
	registration_handlers "github.com/novabankapp/usermanagement.application/handlers/registration"
)

type RegistrationCommands struct {
	RegisterUser registration_handlers.RegisterUserCmdHandler
}

func NewRegistrationCommands(registerUser registration_handlers.RegisterUserCmdHandler) *RegistrationCommands {
	return &RegistrationCommands{
		RegisterUser: registerUser,
	}
}

type RegisterUserCommand struct {
	Dto registration_dtos.RegisterUserDto
}

func NewRegisterUserCommand(dto registration_dtos.RegisterUserDto) *RegisterUserCommand {
	return &RegisterUserCommand{
		dto,
	}
}
