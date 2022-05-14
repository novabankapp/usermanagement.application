package dtos

import "github.com/novabankapp/usermanagement.data/domain/registration"

type GetUserByIdResponse struct {
}

func GetResponseFromUser(user registration.User) *GetUserByIdResponse {
	return &GetUserByIdResponse{}
}
