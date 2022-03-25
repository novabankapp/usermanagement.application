package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/novabankapp/golang.common.infrastructure/kafka"
	kafkaClient "github.com/novabankapp/golang.common.infrastructure/kafka"
	"github.com/novabankapp/golang.common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.data/domain"
	"github.com/novabankapp/usermanagement.data/repositories"
	kafka_go "github.com/segmentio/kafka-go"
	"time"
)

type CreateUserCmdHandler interface {
	Handle(ctx context.Context, command *CreateUserCommand) (*string, error)
}
type createUserCmdHandler struct {
	log           logger.Logger
	cfg           *kafka.Config
	repo          repositories.UserRepository
	kafkaProducer kafkaClient.Producer
}

func NewCreateUserHandler(log logger.Logger, cfg *kafka.Config,
	repo repositories.UserRepository, kafkaProducer kafkaClient.Producer) CreateUserCmdHandler {
	return &createUserCmdHandler{log: log, cfg: cfg, repo: repo, kafkaProducer: kafkaProducer}
}
func (c *createUserCmdHandler) Handle(ctx context.Context, command *CreateUserCommand) (*string, error) {
	userDto := domain.User{}

	user, err := c.repo.Create(ctx, userDto)
	if err != nil {
		return nil, err
	}
	res := new(bytes.Buffer)
	json.NewEncoder(res).Encode(userDto)
	msgBytes := res.Bytes()
	message := kafka_go.Message{
		Topic: c.cfg.Topics.UserCreated.TopicName,
		Value: msgBytes,
		Time:  time.Now().UTC(),
		Key:   []byte(*user),
	}

	error := c.kafkaProducer.PublishMessage(ctx, message)
	return user, error
}
