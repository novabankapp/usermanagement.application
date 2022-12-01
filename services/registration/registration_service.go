package registration

import (
	"context"
	registrationDtos "github.com/novabankapp/usermanagement.application/dtos/registration"
)

type UserRegistrationService interface {
	Register(ctx context.Context, user registrationDtos.RegisterUserDto) (*string, error)
	VerifyOTP(ctx context.Context, channels registrationDtos.VerificationChannels, otp string) (bool, error)
	ResendOTP(ctx context.Context, channels registrationDtos.VerificationChannels) (bool, error)
}
