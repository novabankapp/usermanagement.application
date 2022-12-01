package dtos

type CreateUserDto struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
