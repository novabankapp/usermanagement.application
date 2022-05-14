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

type DeleteUserCmdHandler interface {
	Handle(ctx context.Context, command *DeleteUserCommand) (bool, error)
}
type deleteUserCmdHandler struct {
	log           logger.Logger
	cfg           *kafkaClient.Config
	repo          users.UserRepository
	kafkaProducer kafkaClient.Producer
}

func NewDeleteUserHandler(log logger.Logger, cfg *kafkaClient.Config,
	repo users.UserRepository, kafkaProducer kafkaClient.Producer) DeleteUserCmdHandler {
	return &deleteUserCmdHandler{log: log, cfg: cfg, repo: repo, kafkaProducer: kafkaProducer}
}
func (c *deleteUserCmdHandler) Handle(ctx context.Context, command *DeleteUserCommand) (bool, error) {
	userDto := registration.User{}
	result, err := c.repo.Delete(ctx, userDto)
	if err != nil {
		return false, err
	}
	res := new(bytes.Buffer)
	json.NewEncoder(res).Encode(userDto)
	msgBytes := res.Bytes()
	message := kafka_go.Message{
		Topic: c.cfg.KafkaTopics.UserDeleted.TopicName,
		Value: msgBytes,
		Time:  time.Now().UTC(),
		Key:   []byte(userDto.ID),
	}

	error := c.kafkaProducer.PublishMessage(ctx, message)
	return result, error
}
