package registration_dtos

import "time"

type RegisterUserDto struct {
	FirstName           string              `json:"firstname"`
	LastName            string              `json:"lastname"`
	UserName            string              `json:"username"`
	CreatedAt           time.Time           `json:"created_at"`
	UpdatedAt           time.Time           `json:"updated_at"`
	Email               string              `json:"email"`
	Phone               string              `json:"phone"`
	Password            string              `json:"password"`
	Pin                 string              `json:"pin"`
	VerificationChannel VerificationChannel `json:"verification_channel"`
}
