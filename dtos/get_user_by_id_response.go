package dtos

import "github.com/novabankapp/usermanagement.data/domain"

type GetUserByIdResponse struct {
}

func GetResponseFromUser(user domain.User) *GetUserByIdResponse {
	return &GetUserByIdResponse{}
}
