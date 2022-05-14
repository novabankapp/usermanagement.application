package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"time"

	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.data/domain/registration"
	"github.com/novabankapp/usermanagement.data/repositories/users"
	kafka_go "github.com/segmentio/kafka-go"
)

type CreateUserCmdHandler interface {
	Handle(ctx context.Context, command *CreateUserCommand) (*string, error)
}
type createUserCmdHandler struct {
	log           logger.Logger
	cfg           *kafkaClient.Config
	repo          users.UserRepository
	kafkaProducer kafkaClient.Producer
}

func NewCreateUserHandler(log logger.Logger, cfg *kafkaClient.Config,
	repo users.UserRepository, kafkaProducer kafkaClient.Producer) CreateUserCmdHandler {
	return &createUserCmdHandler{log: log, cfg: cfg, repo: repo, kafkaProducer: kafkaProducer}
}
func (c *createUserCmdHandler) Handle(ctx context.Context, command *CreateUserCommand) (*string, error) {
	userDto := registration.User{}

	user, err := c.repo.Create(ctx, userDto)
	if err != nil {
		return nil, err
	}
	res := new(bytes.Buffer)
	json.NewEncoder(res).Encode(userDto)
	msgBytes := res.Bytes()
	message := kafka_go.Message{
		Topic: c.cfg.KafkaTopics.UserCreated.TopicName,
		Value: msgBytes,
		Time:  time.Now().UTC(),
		Key:   []byte(*user),
	}

	error := c.kafkaProducer.PublishMessage(ctx, message)
	return user, error
}
