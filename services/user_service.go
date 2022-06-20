package services

import (
	"context"
	baseService "github.com/novabankapp/common.application/services/base"
	"github.com/novabankapp/common.application/services/message_queue"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.application/commands"
	"github.com/novabankapp/usermanagement.application/dtos"
	"github.com/novabankapp/usermanagement.application/queries"
	loginDomain "github.com/novabankapp/usermanagement.data/domain/login"
	authRepository "github.com/novabankapp/usermanagement.data/repositories/auth"
	"github.com/novabankapp/usermanagement.data/repositories/users"
)

type UserService interface {
	DeleteUser(ctx context.Context, dto dtos.DeleteUserDto) (bool, error)
	UpdateUser(ctx context.Context, dto dtos.UpdateUserDto) (bool, error)
	GetUserById(ctx context.Context, id string) (*dtos.GetUserByIdResponse, error)
	CheckUsername(cxt context.Context, username string) (bool, error)
	CheckEmail(cxt context.Context, email string) (bool, error)
	GetUsers(ctx context.Context, query string, orderBy string, page int, pageSize int) (*dtos.GetUsersResponse, error)
}
type Repos struct {
	LoginRepo baseService.NoSqlService[loginDomain.UserLogin]
	UsersRepo users.UserRepository
	AuthRepo  authRepository.AuthRepository
}
type userService struct {
	Repos        Repos
	Commands     *commands.UserCommands
	Queries      *queries.UserQueries
	messageQueue message_queue.MessageQueue
	topics       *kafkaClient.KafkaTopics
}

func NewUserService(log logger.Logger,
	messageQueue message_queue.MessageQueue,
	topics *kafkaClient.KafkaTopics,
	repos Repos) UserService {

	updateUserHandler := commands.NewUpdateUserHandler(log, messageQueue, topics, repos)
	deleteUserHandler := commands.NewDeleteUserHandler(log, messageQueue, topics, repos.UsersRepo, repos.AuthRepo)

	getUserByIdHandler := queries.NewGetUserByIdHandler(log, messageQueue, topics, repos.UsersRepo)
	getUsersHandler := queries.NewGetUsersHandler(log, messageQueue, topics, repos.UsersRepo)

	usersCommands := commands.NewUserCommands(updateUserHandler, deleteUserHandler)
	usersQueries := queries.NewUsersQueries(getUserByIdHandler, getUsersHandler)

	return &userService{Commands: usersCommands, Queries: usersQueries, Repos: repos}
}

func (s *userService) CheckUsername(cxt context.Context, username string) (bool, error) {
	return s.Repos.AuthRepo.CheckUsername(cxt, username)
}

func (s *userService) CheckEmail(cxt context.Context, email string) (bool, error) {
	return s.Repos.AuthRepo.CheckEmail(cxt, email)
}

func (s *userService) DeleteUser(ctx context.Context, dto dtos.DeleteUserDto) (bool, error) {
	return s.Commands.DeleteUser.Handle(ctx, commands.NewDeleteUserCommand(
		dto,
	))
}
func (s *userService) UpdateUser(ctx context.Context, dto dtos.UpdateUserDto) (bool, error) {
	return s.Commands.UpdateUser.Handle(ctx, commands.NewUpdateUserCommand(
		dto,
	))
}
func (s *userService) GetUserById(ctx context.Context, id string) (*dtos.GetUserByIdResponse, error) {
	return s.Queries.GetUserById.Handle(ctx, queries.NewGetUserByIdQuery(id))
}
func (s *userService) GetUsers(ctx context.Context, query string, orderBy string, page int, pageSize int) (*dtos.GetUsersResponse, error) {
	return s.Queries.GetUsers.Handle(ctx, queries.NewGetUsersQuery(
		query,
		page,
		pageSize,
		orderBy,
	))
}
