package registration_handlers

import (
	"bytes"
	"context"
	"encoding/json"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	registrationcommands "github.com/novabankapp/usermanagement.application/commands/registration"
	"github.com/novabankapp/usermanagement.data/domain/account"
	"github.com/novabankapp/usermanagement.data/domain/login"
	"github.com/novabankapp/usermanagement.data/domain/registration"
	auth_repository "github.com/novabankapp/usermanagement.data/repositories/auth"
	reg_repo "github.com/novabankapp/usermanagement.data/repositories/registration"
	kafka_go "github.com/segmentio/kafka-go"
	"time"
)

type RegisterUserCmdHandler interface {
	Handle(ctx context.Context, command *registrationcommands.RegisterUserCommand) (*string, error)
}
type registerUserCmdHandler struct {
	log           logger.Logger
	cfg           *kafkaClient.Config
	repo          reg_repo.RegisterRepository
	authRepo      auth_repository.AuthRepository
	kafkaProducer kafkaClient.Producer
}

func NewRegisterUserHandler(log logger.Logger, cfg *kafkaClient.Config,
	repo reg_repo.RegisterRepository, authRepo auth_repository.AuthRepository,
	kafkaProducer kafkaClient.Producer) RegisterUserCmdHandler {
	return &registerUserCmdHandler{log: log, cfg: cfg, repo: repo, authRepo: authRepo,
		kafkaProducer: kafkaProducer}
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
	message := kafka_go.Message{
		Topic: r.cfg.KafkaTopics.UserCreated.TopicName,
		Value: msgBytes,
		Time:  time.Now().UTC(),
		Key:   []byte(*userId),
	}

	error = r.kafkaProducer.PublishMessage(ctx, message)

	return userId, error

}
