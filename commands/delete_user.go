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

type DeleteUserCmdHandler interface {
	Handle(ctx context.Context, command *DeleteUserCommand) (bool, error)
}
type deleteUserCmdHandler struct {
	log           logger.Logger
	cfg           *kafka.Config
	repo          repositories.UserRepository
	kafkaProducer kafkaClient.Producer
}

func NewDeleteUserHandler(log logger.Logger, cfg *kafka.Config,
	repo repositories.UserRepository, kafkaProducer kafkaClient.Producer) DeleteUserCmdHandler {
	return &deleteUserCmdHandler{log: log, cfg: cfg, repo: repo, kafkaProducer: kafkaProducer}
}
func (c *deleteUserCmdHandler) Handle(ctx context.Context, command *DeleteUserCommand) (bool, error) {
	userDto := domain.User{}
	result, err := c.repo.Delete(ctx, userDto)
	if err != nil {
		return false, err
	}
	res := new(bytes.Buffer)
	json.NewEncoder(res).Encode(userDto)
	msgBytes := res.Bytes()
	message := kafka_go.Message{
		Topic: c.cfg.Topics.UserDeleted.TopicName,
		Value: msgBytes,
		Time:  time.Now().UTC(),
		Key:   []byte(userDto.ID),
	}

	error := c.kafkaProducer.PublishMessage(ctx, message)
	return result, error
}
