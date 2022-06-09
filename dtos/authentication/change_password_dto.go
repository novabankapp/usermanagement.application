package authentication

type ChangePasswordDto struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
	UserId      string `json:"user_id"`
}
