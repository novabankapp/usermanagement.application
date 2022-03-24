package handlers

import (
	"context"
	"github.com/novabankapp/golang.common.infrastructure/kafka"
	kafkaClient "github.com/novabankapp/golang.common.infrastructure/kafka"
	"github.com/novabankapp/golang.common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.application/commands"
	"github.com/novabankapp/usermanagement.data/repositories"
)

type CreateUserCmdHandler interface {
	Handle(ctx context.Context, command *commands.CreateUserCommand) error
}
type createProductHandler struct {
	log           logger.Logger
	cfg           *kafka.Config
	repo          repositories.UserRepository
	kafkaProducer kafkaClient.Producer
}

func NewCreateUserHandler(log logger.Logger, cfg *kafka.Config, repo repositories.UserRepository, kafkaProducer kafkaClient.Producer) *createProductHandler {
	return &createProductHandler{log: log, cfg: cfg, repo: repo, kafkaProducer: kafkaProducer}
}
