package services

import (
	"github.com/novabankapp/golang.common.infrastructure/kafka"
	kafkaClient "github.com/novabankapp/golang.common.infrastructure/kafka"
	"github.com/novabankapp/golang.common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.application/commands"
	"github.com/novabankapp/usermanagement.application/queries"
	"github.com/novabankapp/usermanagement.data/repositories"
)

type UserService struct {
	Commands *commands.UserCommands
	Queries  *queries.UserQueries
}

func NewUserService(log logger.Logger, cfg *kafka.Config,
	kafkaProducer kafkaClient.Producer, repo repositories.UserRepository) *UserService {

	createUserHandler := commands.NewCreateUserHandler(log, cfg, repo, kafkaProducer)
	updateUserHandler := commands.NewUpdateUserHandler(log, cfg, repo, kafkaProducer)
	deleteUserHandler := commands.NewDeleteUserHandler(log, cfg, repo, kafkaProducer)

	getUserByIdHandler := queries.NewGetUserByIdHandler(log, cfg, repo)
	getUsersHandler := queries.NewGetUsersHandler(log, cfg, repo)

	usersCommands := commands.NewUserCommands(createUserHandler, updateUserHandler, deleteUserHandler)
	usersQueries := queries.NewUsersQueries(getUserByIdHandler, getUsersHandler)

	return &UserService{Commands: usersCommands, Queries: usersQueries}
}
