package registration

import (
	"context"
	"errors"
	"github.com/novabankapp/common.infrastructure/kafka"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/common.notifier/sms"
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
	ResendPhoneOTP(ctx context.Context, phoneNumber string) (bool, error)
}

type ussdRegistrationService struct {
	notifier sms.SMSService
	repo     registration.RegisterRepository
	authRepo auth_repository.AuthRepository
	baseRepo baseRepository.PostgresRepository[regDomain.PhoneVerificationCode]
	Commands registrationcommands.RegistrationCommands
}

func NewUSSDDRegistrationService(log logger.Logger, cfg *kafka.Config,
	kafkaProducer kafkaClient.Producer,
	notifier sms.SMSService,
	baseRepo baseRepository.PostgresRepository[regDomain.PhoneVerificationCode],
	repo registration.RegisterRepository,
	authRepo auth_repository.AuthRepository) USSDRegistrationService {
	regUserHandler := registration_handlers.NewRegisterUserHandler(log, cfg, repo, authRepo, kafkaProducer)
	registerCommands := registrationcommands.NewRegistrationCommands(regUserHandler)
	return &ussdRegistrationService{
		notifier: notifier,
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
		//To-Do - generate pin and send to phone
		u.notifier.SendSMS("", phoneNumber, "")

		u.baseRepo.Create(ctx, regDomain.PhoneVerificationCode{
			Phone:      phoneNumber,
			Used:       false,
			ExpiryDate: time.Now().Add(time.Minute * 30),
		})
	}
	return result, err

}

func (u ussdRegistrationService) VerifyPhone(ctx context.Context, phoneNumber string, otp string) (bool, error) {
	res, err := u.baseRepo.Get(ctx, 1, 1, &regDomain.PhoneVerificationCode{
		Phone: phoneNumber,
		Code:  otp,
		Used:  false,
	}, "")
	if err != nil {
		return false, err
	}
	if res != nil {
		result := *res
		ver := result[0]
		now := time.Now()
		if ver.ExpiryDate.Before(now) {
			return false, errors.New("code expired")
		}
		return true, nil
	}
	return false, errors.New("code not found")

}
func (u ussdRegistrationService) ResendPhoneOTP(ctx context.Context, phoneNumber string) (bool, error) {
	_, err := u.baseRepo.Create(ctx, regDomain.PhoneVerificationCode{
		Phone:      phoneNumber,
		Used:       false,
		ExpiryDate: time.Now().Add(time.Minute * 30),
	})
	if err != nil {
		return false, err
	}
	//To-Do - generate pin and send to phone
	u.notifier.SendSMS("", phoneNumber, "")
	return true, nil

}
