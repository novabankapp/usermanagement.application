package registration

import (
	"context"
	"github.com/novabankapp/common.infrastructure/kafka"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	registrationcommands "github.com/novabankapp/usermanagement.application/commands/registration"
	registration_dtos "github.com/novabankapp/usermanagement.application/dtos/registration"
	registration_handlers "github.com/novabankapp/usermanagement.application/handlers/registration"
	regDomain "github.com/novabankapp/usermanagement.data/domain/registration"
	auth_repository "github.com/novabankapp/usermanagement.data/repositories/auth"
	baseRepository "github.com/novabankapp/usermanagement.data/repositories/base/postgres"
	"github.com/novabankapp/usermanagement.data/repositories/registration"
	"time"
)

type USSDRegistrationService interface {
	Register(ctx context.Context, phoneNumber string, pin string) (*string, error)
	VerifyPhone(ctx context.Context, phoneNumber string, otp string) (bool, error)
}

type ussdRegistrationService struct {
	repo     registration.RegisterRepository
	authRepo auth_repository.AuthRepository
	baseRepo baseRepository.PostgresRepository[regDomain.PhoneVerificationCode]
	Commands registrationcommands.RegistrationCommands
}

func NewUSSDDRegistrationService(log logger.Logger, cfg *kafka.Config,
	kafkaProducer kafkaClient.Producer,
	baseRepo baseRepository.PostgresRepository[regDomain.PhoneVerificationCode],
	repo registration.RegisterRepository,
	authRepo auth_repository.AuthRepository) USSDRegistrationService {
	regUserHandler := registration_handlers.NewRegisterUserHandler(log, cfg, repo, authRepo, kafkaProducer)
	registerCommands := registrationcommands.NewRegistrationCommands(regUserHandler)
	return &ussdRegistrationService{
		repo:     repo,
		authRepo: authRepo,
		baseRepo: baseRepo,
		Commands: *registerCommands,
	}
}

func (u ussdRegistrationService) Register(ctx context.Context, phoneNumber string, pin string) (*string, error) {
	result, err := u.Commands.RegisterUser.Handle(ctx, registrationcommands.NewRegisterUserCommand(
		registration_dtos.RegisterUserDto{
			Phone: phoneNumber,
			Pin:   pin,
		},
	))
	//insert phone verification
	if result != nil {
		u.baseRepo.Create(ctx, regDomain.PhoneVerificationCode{
			Phone:      phoneNumber,
			Used:       false,
			ExpiryDate: time.Now().Add(time.Minute * 30),
		})
	}
	return result, err

}

func (u ussdRegistrationService) VerifyPhone(ctx context.Context, phoneNumber string, otp string) (bool, error) {
	_, err := u.baseRepo.Get(ctx, 1, 1, &regDomain.PhoneVerificationCode{
		Phone: phoneNumber,
		Code:  otp,
		Used:  false,
	}, "")
	if err != nil {
		return false, err
	}
	return true, nil

}
