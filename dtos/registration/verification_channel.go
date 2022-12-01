package registration_dtos

type VerificationChannel struct {
	Sms   bool
	Email bool
}
type VerificationChannels struct {
	Phone *string
	Email *string
}
