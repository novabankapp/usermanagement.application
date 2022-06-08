package registration_dtos

type ContactDto struct {
	UserID string `json:"user_id"`
	TypeID int    `json:"type_id"`
	Value  string `json:"value"`
}
