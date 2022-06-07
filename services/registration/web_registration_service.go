package registration

import (
	"context"
	"errors"
	"fmt"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/common.notifier/email"
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

type webRegistrationService struct {
	smsNotifier              sms.SMSService
	mailNotifier             email.MailService
	repo                     registration.RegisterRepository
	authRepo                 authRepository.AuthRepository
	phoneVerificationService baseService.RdbmsService[regDomain.PhoneVerificationCode]
	emailVerificationService baseService.RdbmsService[regDomain.EmailVerificationCode]
	Commands                 registrationCommands.RegistrationCommands
}

func NewWebRegistrationService(log logger.Logger,
	topics *kafkaClient.KafkaTopics,
	smsNotifier sms.SMSService,
	mailNotifier email.MailService,
	messageQueue message_queue.MessageQueue,
	phoneVerificationService baseService.RdbmsService[regDomain.PhoneVerificationCode],
	emailVerificationService baseService.RdbmsService[regDomain.EmailVerificationCode],
	repo registration.RegisterRepository,
	authRepo authRepository.AuthRepository) UserRegistrationService {
	regUserHandler := registrationHandlers.NewRegisterUserHandler(log, topics, messageQueue, repo, authRepo)
	registerCommands := registrationCommands.NewRegistrationCommands(regUserHandler)
	return &webRegistrationService{
		smsNotifier:              smsNotifier,
		mailNotifier:             mailNotifier,
		repo:                     repo,
		authRepo:                 authRepo,
		phoneVerificationService: phoneVerificationService,
		emailVerificationService: emailVerificationService,
		Commands:                 *registerCommands,
	}
}

func (w webRegistrationService) Register(ctx context.Context, user registrationDtos.RegisterUserDto) (*string, error) {
	result, err := w.Commands.RegisterUser.Handle(ctx, registrationCommands.NewRegisterUserCommand(
		user,
	))
	//insert phone verification
	if result != nil {

		pin := services.GenerateOTP(5)

		if user.VerificationChannel.Sms {
			w.phoneVerificationService.Create(ctx, regDomain.PhoneVerificationCode{
				Phone:      user.Phone,
				Used:       false,
				Code:       pin,
				ExpiryDate: time.Now().Add(time.Minute * 30),
			})
			w.smsNotifier.SendSMS(user.Phone, fmt.Sprintf("Your One time pin is %s", pin))
		}
		if user.VerificationChannel.Email {
			w.emailVerificationService.Create(ctx, regDomain.EmailVerificationCode{
				Email:      user.Email,
				Used:       false,
				Code:       pin,
				ExpiryDate: time.Now().Add(time.Minute * 30),
			})
			dest := []string{user.Email}
			w.mailNotifier.SendEmail(dest, fmt.Sprintf("Your One time pin is %s", pin))
		}
	}
	return result, err
}

func (w webRegistrationService) VerifyOTP(ctx context.Context, user registrationDtos.RegisterUserDto, otp string) (bool, error) {
	if user.VerificationChannel.Sms {
		res, err := w.phoneVerificationService.Get(ctx, 1, 1, &regDomain.PhoneVerificationCode{
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
	}
	if user.VerificationChannel.Email {
		res, err := w.emailVerificationService.Get(ctx, 1, 1, &regDomain.EmailVerificationCode{
			Email: user.Email,
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
	}
	return false, errors.New("code not found")
}

func (w webRegistrationService) ResendOTP(ctx context.Context, user registrationDtos.RegisterUserDto) (bool, error) {
	//To-Do - generate pin and send to phone
	pin := services.GenerateOTP(5)
	if user.VerificationChannel.Sms {
		w.phoneVerificationService.Create(ctx, regDomain.PhoneVerificationCode{
			Phone:      user.Phone,
			Used:       false,
			Code:       pin,
			ExpiryDate: time.Now().Add(time.Minute * 30),
		})
		w.smsNotifier.SendSMS(user.Phone, fmt.Sprintf("Your One time pin is %s", pin))
	}
	if user.VerificationChannel.Email {
		w.emailVerificationService.Create(ctx, regDomain.EmailVerificationCode{
			Email:      user.Email,
			Used:       false,
			Code:       pin,
			ExpiryDate: time.Now().Add(time.Minute * 30),
		})
		dest := []string{user.Email}
		w.mailNotifier.SendEmail(dest, fmt.Sprintf("Your One time pin is %s", pin))
	}

	return true, nil
}
