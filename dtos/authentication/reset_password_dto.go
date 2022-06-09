package authentication

type ResetPasswordDto struct {
	Email  bool   `json:"email"`
	SMS    bool   `json:"sms"`
	UserId string `json:"user_id"`
}
