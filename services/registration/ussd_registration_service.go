package registration

import (
	"context"
	"errors"
	"fmt"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/common.notifier/sms"
	registrationCommands "github.com/novabankapp/usermanagement.application/commands/registration"
	registrationDtos "github.com/novabankapp/usermanagement.application/dtos/registration"
	registrationHandlers "github.com/novabankapp/usermanagement.application/handlers/registration"
	"github.com/novabankapp/usermanagement.application/services"
	baseService "github.com/novabankapp/usermanagement.application/services/base"
	"github.com/novabankapp/usermanagement.application/services/message_queue"
	regDomain "github.com/novabankapp/usermanagement.data/domain/registration"
	authRepository "github.com/novabankapp/usermanagement.data/repositories/auth"
	"github.com/novabankapp/usermanagement.data/repositories/registration"
	"time"
)

type ussdRegistrationService struct {
	notifier    sms.SMSService
	repo        registration.RegisterRepository
	authRepo    authRepository.AuthRepository
	baseService baseService.Service[regDomain.PhoneVerificationCode]
	Commands    registrationCommands.RegistrationCommands
}

func NewUSSDDRegistrationService(log logger.Logger,
	topics *kafkaClient.KafkaTopics,
	notifier sms.SMSService,
	messageQueue message_queue.MessageQueue,
	baseService baseService.Service[regDomain.PhoneVerificationCode],
	repo registration.RegisterRepository,
	authRepo authRepository.AuthRepository) UserRegistrationService {
	regUserHandler := registrationHandlers.NewRegisterUserHandler(log, topics, messageQueue, repo, authRepo)
	registerCommands := registrationCommands.NewRegistrationCommands(regUserHandler)
	return &ussdRegistrationService{
		notifier:    notifier,
		repo:        repo,
		authRepo:    authRepo,
		baseService: baseService,
		Commands:    *registerCommands,
	}
}

func (u ussdRegistrationService) Register(ctx context.Context, user registrationDtos.RegisterUserDto) (*string, error) {
	result, err := u.Commands.RegisterUser.Handle(ctx, registrationCommands.NewRegisterUserCommand(
		user,
	))
	//insert phone verification
	if result != nil {
		//To-Do - generate pin and send to phone
		pin := services.GenerateOTP(5)
		u.baseService.Create(ctx, regDomain.PhoneVerificationCode{
			Phone:      user.Phone,
			Used:       false,
			Code:       pin,
			ExpiryDate: time.Now().Add(time.Minute * 30),
		})
		u.notifier.SendSMS(user.Phone, fmt.Sprintf("Your One time pin is %s", pin))

	}
	return result, err

}

func (u ussdRegistrationService) VerifyOTP(ctx context.Context, user registrationDtos.RegisterUserDto, otp string) (bool, error) {
	res, err := u.baseService.Get(ctx, 1, 1, &regDomain.PhoneVerificationCode{
		Phone: user.Phone,
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
func (u ussdRegistrationService) ResendOTP(ctx context.Context, user registrationDtos.RegisterUserDto) (bool, error) {

	//To-Do - generate pin and send to phone
	pin := services.GenerateOTP(5)
	_, err := u.baseService.Create(ctx, regDomain.PhoneVerificationCode{
		Phone:      user.Phone,
		Used:       false,
		Code:       pin,
		ExpiryDate: time.Now().Add(time.Minute * 30),
	})
	if err != nil {
		return false, err
	}
	u.notifier.SendSMS(user.Phone, fmt.Sprintf("Your One time pin is %s", pin))
	return true, nil

}
