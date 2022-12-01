package services

import (
	"context"
	"fmt"
	"github.com/novabankapp/common.application/services/message_queue"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	"github.com/novabankapp/common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.application/commands"
	"github.com/novabankapp/usermanagement.application/dtos"
	"github.com/novabankapp/usermanagement.application/queries"
)

type UserService interface {
	CreateUser(cxt context.Context, dto dtos.CreateUserDto) (*string, error)
	DeleteUser(ctx context.Context, dto dtos.DeleteUserDto) (bool, error)
	UpdateUser(ctx context.Context, dto dtos.UpdateUserDto) (bool, error)
	GetUserById(ctx context.Context, id string) (*dtos.GetUserByIdResponse, error)
	CheckUsername(cxt context.Context, username string) (bool, error)
	CheckEmail(cxt context.Context, email string) (bool, error)
	GetUsers(ctx context.Context, query *string, orderBy *string, page int, pageSize int) (*dtos.GetUsersResponse, error)
}

type userService struct {
	Repos        commands.Repos
	Commands     *commands.UserCommands
	Queries      *queries.UserQueries
	messageQueue message_queue.MessageQueue
	topics       *kafkaClient.KafkaTopics
}

func NewUserService(log logger.Logger,
	messageQueue message_queue.MessageQueue,
	topics *kafkaClient.KafkaTopics,
	repos commands.Repos) UserService {

	updateUserHandler := commands.NewUpdateUserHandler(log, messageQueue, topics, repos)
	createUserHander := commands.NewCreateUserHandler(log, messageQueue, topics, repos.UsersRepo)
	deleteUserHandler := commands.NewDeleteUserHandler(log, messageQueue, topics, repos.UsersRepo, repos.AuthRepo)

	getUserByIdHandler := queries.NewGetUserByIdHandler(log, messageQueue, topics, repos.UsersRepo)
	getUsersHandler := queries.NewGetUsersHandler(log, messageQueue, topics, repos.UsersRepo)

	usersCommands := commands.NewUserCommands(createUserHander, updateUserHandler, deleteUserHandler)
	usersQueries := queries.NewUsersQueries(getUserByIdHandler, getUsersHandler)

	return &userService{Commands: usersCommands, Queries: usersQueries, Repos: repos}
}

func (s *userService) CreateUser(cxt context.Context, dto dtos.CreateUserDto) (*string, error) {
	return s.Commands.CreateUser.Handle(cxt, commands.NewCreateUserCommand(dto))
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
func (s *userService) GetUsers(ctx context.Context, query *string, orderBy *string, page int, pageSize int) (*dtos.GetUsersResponse, error) {
	fmt.Println("Is UserService GetUsers")
	fmt.Println(ctx)
	return s.Queries.GetUsers.Handle(ctx, queries.NewGetUsersQuery(
		query,
		page,
		pageSize,
		orderBy,
	))
}
