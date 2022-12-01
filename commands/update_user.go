package commands

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	baseService "github.com/novabankapp/common.application/services/base"
	"github.com/novabankapp/common.application/services/message_queue"
	"github.com/novabankapp/common.application/utilities"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.application/dtos"
	loginDomain "github.com/novabankapp/usermanagement.data/domain/login"
	regDomain "github.com/novabankapp/usermanagement.data/domain/registration"
	authRepository "github.com/novabankapp/usermanagement.data/repositories/auth"
	"github.com/novabankapp/usermanagement.data/repositories/users"
)

type UpdateUserCmdHandler interface {
	Handle(ctx context.Context, command *UpdateUserCommand) (bool, error)
}
type Repos struct {
	LoginRepo baseService.NoSqlService[loginDomain.UserLogin]
	UsersRepo users.UserRepository
	AuthRepo  authRepository.AuthRepository
}
type updateUserCmdHandler struct {
	log          logger.Logger
	messageQueue message_queue.MessageQueue
	topics       *kafkaClient.KafkaTopics
	repos        Repos
}

func NewUpdateUserHandler(
	log logger.Logger,
	messageQueue message_queue.MessageQueue,
	topics *kafkaClient.KafkaTopics,
	repos Repos) UpdateUserCmdHandler {
	return &updateUserCmdHandler{log: log, messageQueue: messageQueue, topics: topics, repos: repos}
}

func (c *updateUserCmdHandler) Handle(ctx context.Context, command *UpdateUserCommand) (bool, error) {
	res, err := c.repos.UsersRepo.GetUser(ctx, command.dto.UserId)
	if err != nil {
		return false, err
	}
	user := *res
	user = fillUser(user, command.dto)
	queries := make([]map[string]string, 1)
	queries = utilities.MakeQueries(queries, "UserId", "=", command.dto.UserId)
	result, err := c.repos.LoginRepo.GetByCondition(ctx, queries)
	if err != nil {
		return false, err
	}
	login := *result
	login = fillLogin(login, command.dto)

	updated, err := c.repos.UsersRepo.Update(ctx, user)
	if err != nil {
		return false, err
	}
	if !updated {
		return false, errors.New("failed to update user")

	}
	saved, err := c.repos.LoginRepo.Update(ctx, login, login.ID.String())
	if err != nil {
		return false, err
	}
	if saved {
		data := new(bytes.Buffer)
		json.NewEncoder(data).Encode(command.dto)
		msgBytes := data.Bytes()
		_, _ = c.messageQueue.PublishMessage(ctx, msgBytes, command.dto.UserId, c.topics.UserUpdated.TopicName)
	}
	return saved, err
}

func fillLogin(user loginDomain.UserLogin, dto dtos.UpdateUserDto) loginDomain.UserLogin {
	if dto.Phone != "" {
		user.Phone = dto.Phone
	}
	if dto.UserName != "" {
		user.UserName = dto.UserName
	}
	if dto.LastName != "" {
		user.LastName = dto.LastName
	}
	if dto.FirstName != "" {
		user.FirstName = dto.FirstName
	}
	if dto.Email != "" {
		user.Email = dto.Email
	}
	return user
}
func fillUser(user regDomain.User, dto dtos.UpdateUserDto) regDomain.User {
	if dto.Phone != "" {
		user.Phone = dto.Phone
	}
	if dto.UserName != "" {
		user.UserName = dto.UserName
	}
	if dto.LastName != "" {
		user.LastName = dto.LastName
	}
	if dto.FirstName != "" {
		user.FirstName = dto.FirstName
	}
	if dto.Email != "" {
		user.Email = dto.Email
	}
	return user
}
