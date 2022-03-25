package queries

type UserQueries struct {
	GetUserById GetUserByIdHandler
	GetUsers    GetUsersHandler
}

func NewUsersQueries(getUserById GetUserByIdHandler, getUsers GetUsersHandler) *UserQueries {
	return &UserQueries{GetUserById: getUserById, GetUsers: getUsers}
}

type GetUserByIdQuery struct {
	UserID string `json:"userId" validate:"required,gte=0,lte=255"`
}

func NewGetUserByIdQuery(userID string) *GetUserByIdQuery {
	return &GetUserByIdQuery{UserID: userID}
}

type GetUsersQuery struct {
	Query    string `json:"query"`
	Page     int    `json:"page"`
	PageSize int    `json:"pageSize"`
	OrderBy  string `json:"orderBy"`
}

func NewGetUsersQuery(query string, page int, pageSize int, orderBy string) *GetUsersQuery {
	return &GetUsersQuery{Query: query, Page: page, PageSize: pageSize, OrderBy: orderBy}
}
