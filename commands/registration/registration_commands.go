package registration_commands

import (
	registrationHandlers "github.com/novabankapp/usermanagement.application/handlers/registration"
)

type RegistrationCommands struct {
	RegisterUser registrationHandlers.RegisterUserCmdHandler
}

func NewRegistrationCommands(registerUser registrationHandlers.RegisterUserCmdHandler) *RegistrationCommands {
	return &RegistrationCommands{
		RegisterUser: registerUser,
	}
}
