package registration_handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/novabankapp/common.application/services/message_queue"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	registrationCommands "github.com/novabankapp/usermanagement.application/commands/registration"
	"github.com/novabankapp/usermanagement.data/domain/account"
	"github.com/novabankapp/usermanagement.data/domain/login"
	"github.com/novabankapp/usermanagement.data/domain/registration"
	authRepository "github.com/novabankapp/usermanagement.data/repositories/auth"
	regRepo "github.com/novabankapp/usermanagement.data/repositories/registration"
	"time"
)

type RegisterUserCmdHandler interface {
	Handle(ctx context.Context, command *registrationCommands.RegisterUserCommand) (*string, error)
}
type registerUserCmdHandler struct {
	log          logger.Logger
	topics       *kafkaClient.KafkaTopics
	messageQueue message_queue.MessageQueue
	repo         regRepo.RegisterRepository
	authRepo     authRepository.AuthRepository
}

func NewRegisterUserHandler(log logger.Logger,
	topics *kafkaClient.KafkaTopics,
	messageQueue message_queue.MessageQueue,
	repo regRepo.RegisterRepository,
	authRepo authRepository.AuthRepository) RegisterUserCmdHandler {
	return &registerUserCmdHandler{log: log, topics: topics, messageQueue: messageQueue, repo: repo, authRepo: authRepo}
}

func (r registerUserCmdHandler) Handle(ctx context.Context,
	command *registrationCommands.RegisterUserCommand) (*string, error) {
	userDto := registration.User{
		FirstName: command.Dto.FirstName,
		LastName:  command.Dto.LastName,
		UserName:  command.Dto.UserName,
		Email:     command.Dto.Email,
		Phone:     command.Dto.Phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userId, err2 := r.repo.Create(ctx, userDto)
	if err2 != nil {
		return nil, err2
	}
	accountId, _, err := r.authRepo.Create(ctx, account.UserAccount{
		UserID:    *userId,
		CreatedAt: time.Now(),
	}, login.UserLogin{
		UserID:    *userId,
		FirstName: command.Dto.FirstName,
		LastName:  command.Dto.LastName,
		UserName:  command.Dto.UserName,
		Email:     command.Dto.Email,
		Phone:     command.Dto.Phone,
		Password:  command.Dto.Password,
		Pin:       command.Dto.Pin,
	})
	if err != nil {
		return nil, err
	}

	res := new(bytes.Buffer)
	e := json.NewEncoder(res).Encode(userDto)
	if e == nil {
		msgBytes := res.Bytes()

		_, err2 = r.messageQueue.PublishMessage(ctx, msgBytes, *userId, r.topics.UserCreated.TopicName)
	}
	res2 := new(bytes.Buffer)
	er := json.NewEncoder(res).Encode(struct {
		AccountId string
		UserId    string
	}{
		AccountId: *accountId,
		UserId:    *userId,
	})
	if er == nil {
		msgBytes := res2.Bytes()

		_, err = r.messageQueue.PublishMessage(ctx, msgBytes, *userId, r.topics.AccountCreated.TopicName)
	}

	return userId, err2

}
