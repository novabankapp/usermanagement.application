package queries

import (
	"context"

	"github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.application/dtos"
	"github.com/novabankapp/usermanagement.data/repositories/users"
)

type GetUserByIdHandler interface {
	Handle(ctx context.Context, query *GetUserByIdQuery) (*dtos.GetUserByIdResponse, error)
}
type getUserByIdHandler struct {
	log  logger.Logger
	cfg  *kafka.Config
	repo users.UserRepository
}

func NewGetUserByIdHandler(log logger.Logger, cfg *kafka.Config,
	repo users.UserRepository) GetUserByIdHandler {
	return &getUserByIdHandler{log: log, cfg: cfg, repo: repo}
}

func (q *getUserByIdHandler) Handle(ctx context.Context, query *GetUserByIdQuery) (*dtos.GetUserByIdResponse, error) {

	user, err := q.repo.GetUser(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	return dtos.GetResponseFromUser(*user), nil
}
