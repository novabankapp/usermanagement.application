package registration_handlers

import (
	"bytes"
	"context"
	"encoding/json"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	registrationcommands "github.com/novabankapp/usermanagement.application/commands/registration"
	"github.com/novabankapp/usermanagement.application/services/message_queue"
	"github.com/novabankapp/usermanagement.data/domain/account"
	"github.com/novabankapp/usermanagement.data/domain/login"
	"github.com/novabankapp/usermanagement.data/domain/registration"
	auth_repository "github.com/novabankapp/usermanagement.data/repositories/auth"
	reg_repo "github.com/novabankapp/usermanagement.data/repositories/registration"
	"time"
)

type RegisterUserCmdHandler interface {
	Handle(ctx context.Context, command *registrationcommands.RegisterUserCommand) (*string, error)
}
type registerUserCmdHandler struct {
	log          logger.Logger
	topics       *kafkaClient.KafkaTopics
	messageQueue message_queue.MessageQueue
	repo         reg_repo.RegisterRepository
	authRepo     auth_repository.AuthRepository
}

func NewRegisterUserHandler(log logger.Logger,
	topics *kafkaClient.KafkaTopics,
	messageQueue message_queue.MessageQueue,
	repo reg_repo.RegisterRepository, authRepo auth_repository.AuthRepository) RegisterUserCmdHandler {
	return &registerUserCmdHandler{log: log, topics: topics, messageQueue: messageQueue, repo: repo, authRepo: authRepo}
}

func (r registerUserCmdHandler) Handle(ctx context.Context,
	command *registrationcommands.RegisterUserCommand) (*string, error) {
	userDto := registration.User{
		FirstName: command.Dto.FirstName,
		LastName:  command.Dto.LastName,
		UserName:  command.Dto.UserName,
		Email:     command.Dto.Email,
		Phone:     command.Dto.Phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	userId, error := r.repo.Create(ctx, userDto)
	if error != nil {
		return nil, error
	}
	_, err := r.authRepo.Create(ctx, account.UserAccount{
		UserID: *userId,
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
	json.NewEncoder(res).Encode(userDto)
	msgBytes := res.Bytes()

	_, error = r.messageQueue.PublishMessage(ctx, msgBytes, *userId, r.topics.UserCreated.TopicName)

	return userId, error

}
