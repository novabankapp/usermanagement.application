package auth

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.notifier/email"
	"github.com/novabankapp/common.notifier/sms"
	"github.com/novabankapp/usermanagement.application/dtos/authentication"
	"github.com/novabankapp/usermanagement.application/services"
	baseService "github.com/novabankapp/usermanagement.application/services/base"
	"github.com/novabankapp/usermanagement.application/services/message_queue"
	"github.com/novabankapp/usermanagement.application/utilities"
	loginDomain "github.com/novabankapp/usermanagement.data/domain/login"
	passwordDomain "github.com/novabankapp/usermanagement.data/domain/password"
	"time"
)

type PasswordService interface {
	RecoverPassword(ctx context.Context, dto authentication.ResetPasswordDto) (*string, error)
	VerifyResetPin(ctx context.Context, pin string, userId string) (bool, error)
	VerifyResetPhrase(ctx context.Context, phrase string, userId string) (bool, error)
	ChangePassword(ctx context.Context, dto authentication.ChangePasswordDto) (bool, error)
}
type PasswordRepos struct {
	ChangePasswordEmailRepo baseService.NoSqlService[passwordDomain.EmailPasswordReset]
	ChangePasswordPhoneRepo baseService.NoSqlService[passwordDomain.PhonePasswordReset]
	LoginRepo               baseService.NoSqlService[loginDomain.UserLogin]
}
type passwordService struct {
	repos         PasswordRepos
	smsNotifier   sms.SMSService
	emailNotifier email.MailService
	messageQueue  message_queue.MessageQueue
	topics        *kafkaClient.KafkaTopics
}

func NewPasswordService(repos PasswordRepos,
	smsNotifier sms.SMSService,
	emailNotifier email.MailService,
	messageQueue message_queue.MessageQueue,
	topics *kafkaClient.KafkaTopics) PasswordService {
	return &passwordService{
		repos:         repos,
		smsNotifier:   smsNotifier,
		emailNotifier: emailNotifier,
		messageQueue:  messageQueue,
		topics:        topics,
	}
}

func (p passwordService) RecoverPassword(ctx context.Context, dto authentication.ResetPasswordDto) (*string, error) {
	queries := make([]map[string]string, 1)
	queries = utilities.MakeQueries(queries, "UserId", "=", dto.UserId)
	user, err := p.repos.LoginRepo.GetByCondition(ctx, queries)
	if err != nil {
		return nil, errors.New("user not found")
	}
	var usr = *user
	if dto.SMS {
		otp := services.GenerateOTP(6)
		expiry := time.Now().Add(time.Hour * 5)
		_, err := p.repos.ChangePasswordPhoneRepo.Create(ctx, passwordDomain.PhonePasswordReset{
			UserID:     dto.UserId,
			Pin:        otp,
			ExpiryDate: expiry,
		})
		if err != nil {
			return nil, err
		}
		p.smsNotifier.SendSMS(usr.Phone, utilities.FormatPhonePasswordResetMessage(otp, "5 hours"))
		return &otp, nil
	}
	if dto.Email {
		otp := services.GenerateOTP(6)
		hash := services.GenerateSha1Hash(otp)
		expiry := time.Now().Add(time.Hour * 5)
		_, err := p.repos.ChangePasswordEmailRepo.Create(ctx, passwordDomain.EmailPasswordReset{
			UserID:     dto.UserId,
			Phrase:     hash,
			ExpiryDate: expiry,
		})
		if err != nil {
			return nil, err
		}
		dest := []string{usr.Email}
		p.emailNotifier.SendEmail(dest, utilities.FormatEmailPasswordResetMessage(hash, "5 hours"))
		return &hash, nil
	}
	return nil, errors.New("invalid action")
}

func (p passwordService) VerifyResetPin(ctx context.Context, pin string, userId string) (bool, error) {
	queries := make([]map[string]string, 2)
	queries = utilities.MakeQueries(queries, "UserId", "=", userId)
	queries = utilities.MakeQueries(queries, "Pin", "=", pin)
	_, err := p.repos.ChangePasswordPhoneRepo.GetByCondition(ctx, queries)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p passwordService) VerifyResetPhrase(ctx context.Context, phrase string, userId string) (bool, error) {
	queries := make([]map[string]string, 2)
	queries = utilities.MakeQueries(queries, "UserId", "=", userId)
	queries = utilities.MakeQueries(queries, "Phrase", "=", phrase)
	_, err := p.repos.ChangePasswordEmailRepo.GetByCondition(ctx, queries)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (p passwordService) ChangePassword(ctx context.Context, dto authentication.ChangePasswordDto) (bool, error) {
	queries := make([]map[string]string, 1)
	queries = utilities.MakeQueries(queries, "UserId", "=", dto.UserId)
	user, err := p.repos.LoginRepo.GetByCondition(ctx, queries)
	if err != nil {
		return false, errors.New("user not found")
	}
	var usr = *user
	passError := usr.ComparePasswords(dto.OldPassword)
	if passError != nil {
		return false, errors.New("password does not match")
	}
	usr.Password = dto.NewPassword
	e := usr.HashPassword()
	if e != nil {
		return false, e
	}
	done, err := p.repos.LoginRepo.Update(ctx)
	if done {
		val := new(bytes.Buffer)
		e := json.NewEncoder(val).Encode(struct {
			Username string
			Name     string
		}{
			Username: usr.UserName,
			Name:     fmt.Sprintf("%s %s", usr.FirstName, usr.LastName),
		})
		if e == nil {
			msgBytes := val.Bytes()
			_, _ = p.messageQueue.PublishMessage(ctx, msgBytes, dto.UserId, p.topics.UserPasswordChanged.TopicName)
		}
	}
	return done, err
}
