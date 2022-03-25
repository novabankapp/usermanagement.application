package queries

import (
	"context"
	"github.com/novabankapp/golang.common.infrastructure/kafka"
	"github.com/novabankapp/golang.common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.application/dtos"
	"github.com/novabankapp/usermanagement.data/repositories"
)

type GetUserByIdHandler interface {
	Handle(ctx context.Context, query *GetUserByIdQuery) (*dtos.GetUserByIdResponse, error)
}
type getUserByIdHandler struct {
	log  logger.Logger
	cfg  *kafka.Config
	repo repositories.UserRepository
}

func NewGetUserByIdHandler(log logger.Logger, cfg *kafka.Config,
	repo repositories.UserRepository) GetUserByIdHandler {
	return &getUserByIdHandler{log: log, cfg: cfg, repo: repo}
}

func (q *getUserByIdHandler) Handle(ctx context.Context, query *GetUserByIdQuery) (*dtos.GetUserByIdResponse, error) {

	user, err := q.repo.GetUser(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	return dtos.GetResponseFromUser(*user), nil
}
