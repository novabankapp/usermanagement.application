package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/novabankapp/common.application/services"
	baseService "github.com/novabankapp/common.application/services/base"
	"github.com/novabankapp/common.application/services/message_queue"
	"github.com/novabankapp/common.application/utilities"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.notifier/sms"
	"github.com/novabankapp/usermanagement.application/dtos/authentication/login"
	loginDomain "github.com/novabankapp/usermanagement.data/domain/login"
	"github.com/novabankapp/usermanagement.data/repositories/auth"
	"time"
)

type LoginService interface {
	LoginByUsername(ctx context.Context, dto login.UserLoginDto) (bool, error)
	LoginByUserId(ctx context.Context, userId string) (bool, error)
	VerifyLoginOTP(ctx context.Context, dto login.VerifyLoginPinDto) (bool, error)
}
type loginService struct {
	authRepo     auth.AuthRepository
	otpRepo      baseService.NoSqlService[loginDomain.OtpLogin]
	messageQueue message_queue.MessageQueue
	topics       *kafkaClient.KafkaTopics
	notifier     sms.SMSService
}

func NewLoginService(
	authRepo auth.AuthRepository,
	otpRepo baseService.NoSqlService[loginDomain.OtpLogin],
	messageQueue message_queue.MessageQueue,
	topics *kafkaClient.KafkaTopics,
	notifier sms.SMSService,
) LoginService {
	return &loginService{
		authRepo:     authRepo,
		messageQueue: messageQueue,
		topics:       topics,
		notifier:     notifier,
		otpRepo:      otpRepo,
	}
}
func (l loginService) LoginByUsername(ctx context.Context, dto login.UserLoginDto) (bool, error) {
	accounts, err := l.authRepo.Login(ctx, dto.Username, dto.Password)
	if err != nil {
		return false, err
	}

	if accounts != nil && len(*accounts) > 0 {
		accs := *accounts
		val := new(bytes.Buffer)
		e := json.NewEncoder(val).Encode(struct {
			Username string
		}{
			Username: dto.Username,
		})
		if e == nil {
			msgBytes := val.Bytes()
			_, _ = l.messageQueue.PublishMessage(ctx, msgBytes, accs[0].UserID, l.topics.UserLoggedIn.TopicName)
		}
	}
	return true, err
}

func (l loginService) LoginByUserId(ctx context.Context, userId string) (bool, error) {
	usr, err := l.authRepo.GetUserById(ctx, userId)
	if err != nil {
		return false, err
	}
	otp := services.GenerateOTP(6)
	expiry := time.Now().Add(time.Hour * 5)
	_, err = l.otpRepo.Create(ctx, loginDomain.OtpLogin{
		UserID:     userId,
		Pin:        otp,
		ExpiryDate: expiry,
	})
	if err != nil {
		return false, err
	}
	l.notifier.SendSMS(usr.Phone, utilities.FormatPhoneLoginMessage(otp, "5 hours"))

	return true, nil
}

func (l loginService) VerifyLoginOTP(ctx context.Context, dto login.VerifyLoginPinDto) (bool, error) {
	queries := make([]map[string]string, 1)
	queries = utilities.MakeQueries(queries, "UserId", "=", dto.UserId)
	queries = utilities.MakeQueries(queries, "Pin", "=", dto.Pin)
	queries = utilities.MakeQueries(queries, "ExpiryDate", ">", time.Now().String())
	_, err := l.otpRepo.GetByCondition(ctx, queries)
	if err != nil {
		return false, err
	}
	return true, nil

}
