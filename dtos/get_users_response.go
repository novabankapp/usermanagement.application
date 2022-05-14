package dtos

import "github.com/novabankapp/usermanagement.data/domain/registration"

type GetUsersResponse struct {
}

func GetResponseFromUsers(users []registration.User) *GetUsersResponse {
	return &GetUsersResponse{}
}
