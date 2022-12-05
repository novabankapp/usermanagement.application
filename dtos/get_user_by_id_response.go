package dtos

import (
	"github.com/novabankapp/usermanagement.data/domain/registration"
)

type GetUserByIdResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Phone     string
}

func GetResponseFromUser(user registration.User) *GetUserByIdResponse {
	return &GetUserByIdResponse{}
}
