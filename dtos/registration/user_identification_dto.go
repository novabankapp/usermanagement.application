package registration_dtos

import "time"

type UserIdentificationDto struct {
	UserID     string    `json:"user_id" binding:"required"`
	TypeOfID   string    `json:"type_of_id" binding:"required"`
	IDNumber   string    `json:"id_number" binding:"required"`
	IssueDate  time.Time `json:"issue_date" binding:"required"`
	ExpiryDate time.Time `json:"expiry_date" binding:"required"`
}
