package registration_dtos

type ResidenceDetailsDto struct {
	UserID               string `json:"user_id" `
	ResidentialStatus    string `json:"residential_status" `
	ProofOfResidency     string `json:"proof_of_residency" `
	NationalityCountryID int
	CountryOfBirthID     int
}
