package queries

import "github.com/google/uuid"

type UserQueries struct {
	GetUserById GetUserByIdHandler
	GetUsers    GetUsersHandler
}

func NewUsersQueries(getUserById GetUserByIdHandler, getUsers GetUsersHandler) *UserQueries {
	return &UserQueries{GetUserById: getUserById, GetUsers: getUsers}
}

type GetUserByIdQuery struct {
	UserID uuid.UUID `json:"userId" validate:"required,gte=0,lte=255"`
}

func NewGetUserByIdQuery(userID uuid.UUID) *GetUserByIdQuery {
	return &GetUserByIdQuery{UserID: userID}
}

type GetUsersQuery struct {
	Query    string `json:"query"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	OrderBy  string `json:"orderBy"`
}

func NewSearchProductQuery(query string, page int, pageSize int, orderBy string) *GetUsersQuery {
	return &GetUsersQuery{Query: query, Page: page, PageSize: pageSize, OrderBy: orderBy}
}
