package login

type VerifyLoginPinDto struct {
	UserId string `json:"user_id"`
	Pin    string `json:"pin"`
}
