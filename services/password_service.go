package services

type PasswordService interface {
}
type passwordService struct {
}

func NewPasswordService() PasswordService {
	return passwordService{}
}
