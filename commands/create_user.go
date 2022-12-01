package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/novabankapp/common.application/services/message_queue"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.data/domain/registration"
	"github.com/novabankapp/usermanagement.data/repositories/users"
)

type CreateUserCmdHandler interface {
	Handle(ctx context.Context, command *CreateUserCommand) (*string, error)
}
type createUserCmdHandler struct {
	log          logger.Logger
	cfg          *kafkaClient.Config
	repo         users.UserRepository
	messageQueue message_queue.MessageQueue
	topics       *kafkaClient.KafkaTopics
}

func NewCreateUserHandler(
	log logger.Logger,
	messageQueue message_queue.MessageQueue,
	topics *kafkaClient.KafkaTopics,
	repo users.UserRepository) CreateUserCmdHandler {
	return &createUserCmdHandler{log: log, topics: topics, repo: repo, messageQueue: messageQueue}
}
func (c *createUserCmdHandler) Handle(ctx context.Context, command *CreateUserCommand) (*string, error) {
	userDto := registration.User{
		Phone:     command.dto.Phone,
		FirstName: command.dto.FirstName,
		LastName:  command.dto.LastName,
		UserName:  command.dto.UserName,
		Email:     command.dto.Email,
	}

	user, err := c.repo.Create(ctx, userDto)
	if err != nil {
		return nil, err
	}
	res := new(bytes.Buffer)
	json.NewEncoder(res).Encode(userDto)
	msgBytes := res.Bytes()
	_, _ = c.messageQueue.PublishMessage(ctx, msgBytes, command.dto.UserName, c.topics.UserCreated.TopicName)
	return user, err
}
