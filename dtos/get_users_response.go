package dtos

import "github.com/novabankapp/usermanagement.data/domain"

type GetUsersResponse struct {
}

func GetResponseFromUsers(users []domain.User) *GetUsersResponse {
	return &GetUsersResponse{}
}
