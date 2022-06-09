package registration_dtos

type ContactDto struct {
	ID     *uint  `json:"id"`
	UserID string `json:"user_id"`
	TypeID int    `json:"type_id"`
	Value  string `json:"value"`
}
