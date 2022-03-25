package services

import (
	"context"
	"github.com/novabankapp/golang.common.infrastructure/kafka"
	kafkaClient "github.com/novabankapp/golang.common.infrastructure/kafka"
	"github.com/novabankapp/golang.common.infrastructure/logger"
	"github.com/novabankapp/usermanagement.application/commands"
	"github.com/novabankapp/usermanagement.application/dtos"
	"github.com/novabankapp/usermanagement.application/queries"
	"github.com/novabankapp/usermanagement.data/repositories"
)

type UserService interface {
	CreateUser(ctx context.Context, dto dtos.CreateUserDto) (*string, error)
	DeleteUser(ctx context.Context, dto dtos.DeleteUserDto) (bool, error)
	UpdateUser(ctx context.Context, dto dtos.UpdateUserDto) (bool, error)
	GetUserById(ctx context.Context, id string) (*dtos.GetUserByIdResponse, error)
	GetUsers(ctx context.Context, query string, orderBy string, page int, pageSize int) (*dtos.GetUsersResponse, error)
}
type userService struct {
	Repo     repositories.UserRepository
	Commands *commands.UserCommands
	Queries  *queries.UserQueries
}

func NewUserService(log logger.Logger, cfg *kafka.Config,
	kafkaProducer kafkaClient.Producer, repo repositories.UserRepository) UserService {

	createUserHandler := commands.NewCreateUserHandler(log, cfg, repo, kafkaProducer)
	updateUserHandler := commands.NewUpdateUserHandler(log, cfg, repo, kafkaProducer)
	deleteUserHandler := commands.NewDeleteUserHandler(log, cfg, repo, kafkaProducer)

	getUserByIdHandler := queries.NewGetUserByIdHandler(log, cfg, repo)
	getUsersHandler := queries.NewGetUsersHandler(log, cfg, repo)

	usersCommands := commands.NewUserCommands(createUserHandler, updateUserHandler, deleteUserHandler)
	usersQueries := queries.NewUsersQueries(getUserByIdHandler, getUsersHandler)

	return &userService{Commands: usersCommands, Queries: usersQueries, Repo: repo}
}
func (s *userService) CreateUser(ctx context.Context, dto dtos.CreateUserDto) (*string, error) {
	return s.Commands.CreateUser.Handle(ctx, commands.NewCreateUserCommand(
		dto,
	))

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
