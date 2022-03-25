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

type UpdateUserCmdHandler interface {
	Handle(ctx context.Context, command *UpdateUserCommand) (bool, error)
}
type updateUserCmdHandler struct {
	log           logger.Logger
	cfg           *kafka.Config
	repo          repositories.UserRepository
	kafkaProducer kafkaClient.Producer
}

func NewUpdateUserHandler(log logger.Logger, cfg *kafka.Config,
	repo repositories.UserRepository, kafkaProducer kafkaClient.Producer) UpdateUserCmdHandler {
	return &updateUserCmdHandler{log: log, cfg: cfg, repo: repo, kafkaProducer: kafkaProducer}
}
func (c *updateUserCmdHandler) Handle(ctx context.Context, command *UpdateUserCommand) (bool, error) {
	userDto := domain.User{}
	result, err := c.repo.Update(ctx, userDto)
	if err != nil {
		return false, err
	}
	res := new(bytes.Buffer)
	json.NewEncoder(res).Encode(userDto)
	msgBytes := res.Bytes()
	message := kafka_go.Message{
		Topic: c.cfg.Topics.UserUpdated.TopicName,
		Value: msgBytes,
		Time:  time.Now().UTC(),
		Key:   []byte(userDto.ID),
	}

	error := c.kafkaProducer.PublishMessage(ctx, message)
	return result, error
}
