package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/novabankapp/common.application/services/message_queue"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.data/domain/registration"
	authRepository "github.com/novabankapp/usermanagement.data/repositories/auth"
	"github.com/novabankapp/usermanagement.data/repositories/users"
)

type DeleteUserCmdHandler interface {
	Handle(ctx context.Context, command *DeleteUserCommand) (bool, error)
}
type deleteUserCmdHandler struct {
	log          logger.Logger
	repo         users.UserRepository
	authRepo     authRepository.AuthRepository
	messageQueue message_queue.MessageQueue
	topics       *kafkaClient.KafkaTopics
}

func NewDeleteUserHandler(
	log logger.Logger,
	messageQueue message_queue.MessageQueue,
	topics *kafkaClient.KafkaTopics,
	repo users.UserRepository,
	authRepo authRepository.AuthRepository,
) DeleteUserCmdHandler {
	return &deleteUserCmdHandler{log: log, repo: repo, authRepo: authRepo, messageQueue: messageQueue, topics: topics}
}
func (c *deleteUserCmdHandler) Handle(ctx context.Context, command *DeleteUserCommand) (bool, error) {
	userDto := registration.User{}
	result, err := c.repo.Delete(ctx, userDto)
	if err != nil {
		return false, err
	}
	_, err = c.authRepo.DeleteUser(ctx, command.dto.UserId)
	if err != nil {
		//return false, err
	}
	res := new(bytes.Buffer)
	json.NewEncoder(res).Encode(userDto)
	msgBytes := res.Bytes()

	_, error := c.messageQueue.PublishMessage(ctx, msgBytes, userDto.ID, c.topics.UserDeleted.TopicName)
	return result, error
}
