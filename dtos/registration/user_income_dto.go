package registration_dtos

import "github.com/shopspring/decimal"

type UserIncomeDto struct {
	UserID        string          `json:"user_id"`
	Source        string          `json:"source"`
	MonthlyIncome decimal.Decimal `json:"monthly_income"`
	ProofOfSource string          `json:"proof_of_source"`
}
