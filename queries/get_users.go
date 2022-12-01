package queries

import (
	"context"
	"fmt"
	"github.com/novabankapp/common.application/services/message_queue"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.application/dtos"
	"github.com/novabankapp/usermanagement.data/repositories/users"
)

type GetUsersHandler interface {
	Handle(ctx context.Context, query *GetUsersQuery) (*dtos.GetUsersResponse, error)
}
type getUsersHandler struct {
	log          logger.Logger
	messageQueue message_queue.MessageQueue
	topics       *kafkaClient.KafkaTopics
	repo         users.UserRepository
}

func NewGetUsersHandler(log logger.Logger,
	messageQueue message_queue.MessageQueue,
	topics *kafkaClient.KafkaTopics,
	repo users.UserRepository) GetUsersHandler {
	return &getUsersHandler{log: log, messageQueue: messageQueue, topics: topics, repo: repo}
}

func (q *getUsersHandler) Handle(ctx context.Context, query *GetUsersQuery) (*dtos.GetUsersResponse, error) {
	fmt.Println("In GetUsers Handle")
	fmt.Println(ctx)
	users, err := q.repo.GetUsers(ctx, query.Page, query.PageSize, query.Query, query.OrderBy)
	fmt.Println(err.Error())
	if err != nil {
		return nil, err
	}
	return dtos.GetResponseFromUsers(*users), nil
}
