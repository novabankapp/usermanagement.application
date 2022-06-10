package dtos

type UpdateUserDto struct {
	UserId    string `json:"userId"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	UserName  string `json:"username"`
	Email     string `json:"email"`
	Phone     string `json:"phone"`
}
