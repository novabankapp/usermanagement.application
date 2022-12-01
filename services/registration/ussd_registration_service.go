package registration

import (
	"context"
	"errors"
	"fmt"
	commonServices "github.com/novabankapp/common.application/services"
	baseService "github.com/novabankapp/common.application/services/base"
	"github.com/novabankapp/common.application/services/message_queue"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/common.notifier/sms"
	registrationCommands "github.com/novabankapp/usermanagement.application/commands/registration"
	registrationDtos "github.com/novabankapp/usermanagement.application/dtos/registration"
	registrationHandlers "github.com/novabankapp/usermanagement.application/handlers/registration"
	regDomain "github.com/novabankapp/usermanagement.data/domain/registration"
	authRepository "github.com/novabankapp/usermanagement.data/repositories/auth"
	"github.com/novabankapp/usermanagement.data/repositories/registration"
	"time"
)

type UssdRegistrationService struct {
	notifier    sms.SMSService
	repo        registration.RegisterRepository
	authRepo    authRepository.AuthRepository
	baseService baseService.RdbmsService[regDomain.PhoneVerificationCode]
	Commands    registrationCommands.RegistrationCommands
}

func NewUSSDDRegistrationService(log logger.Logger,
	topics *kafkaClient.KafkaTopics,
	notifier sms.SMSService,
	messageQueue message_queue.MessageQueue,
	baseService baseService.RdbmsService[regDomain.PhoneVerificationCode],
	repo registration.RegisterRepository,
	authRepo authRepository.AuthRepository) UssdRegistrationService {
	regUserHandler := registrationHandlers.NewRegisterUserHandler(log, topics, messageQueue, repo.Create, authRepo.Create)
	registerCommands := registrationCommands.NewRegistrationCommands(regUserHandler)
	return UssdRegistrationService{
		notifier:    notifier,
		repo:        repo,
		authRepo:    authRepo,
		baseService: baseService,
		Commands:    *registerCommands,
	}
}

func (u UssdRegistrationService) Register(ctx context.Context, user registrationDtos.RegisterUserDto) (*string, error) {
	result, err := u.Commands.RegisterUser.Handle(ctx, registrationCommands.NewRegisterUserCommand(
		user,
	))
	//insert phone verification
	if result != nil {
		//To-Do - generate pin and send to phone
		pin := commonServices.GenerateOTP(5)
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

func (u UssdRegistrationService) VerifyOTP(ctx context.Context, channels registrationDtos.VerificationChannels, otp string) (bool, error) {
	if channels.Phone != nil {
		res, err := u.baseService.Get(ctx, 1, 1, &regDomain.PhoneVerificationCode{
			Phone: *channels.Phone,
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
	return false, errors.New("no verification channel was selected")

}
func (u UssdRegistrationService) ResendOTP(ctx context.Context, user registrationDtos.RegisterUserDto) (bool, error) {

	//To-Do - generate pin and send to phone
	pin := commonServices.GenerateOTP(5)
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
