package queries

import (
	"context"

	"github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.application/dtos"
	"github.com/novabankapp/usermanagement.data/repositories/users"
)

type GetUsersHandler interface {
	Handle(ctx context.Context, query *GetUsersQuery) (*dtos.GetUsersResponse, error)
}
type getUsersHandler struct {
	log  logger.Logger
	cfg  *kafka.Config
	repo users.UserRepository
}

func NewGetUsersHandler(log logger.Logger, cfg *kafka.Config,
	repo users.UserRepository) GetUsersHandler {
	return &getUsersHandler{log: log, cfg: cfg, repo: repo}
}

func (q *getUsersHandler) Handle(ctx context.Context, query *GetUsersQuery) (*dtos.GetUsersResponse, error) {

	users, err := q.repo.GetUsers(ctx, query.Page, query.PageSize, query.Query, query.OrderBy)
	if err != nil {
		return nil, err
	}
	return dtos.GetResponseFromUsers(*users), nil
}
