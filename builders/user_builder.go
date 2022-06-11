package builders

import (
	"github.com/novabankapp/usermanagement.data/domain/registration"
	"github.com/shopspring/decimal"
	"time"
)

type UserBuilder struct {
	user *registration.User
}

func NewUserBuilder() *UserBuilder {
	return &UserBuilder{
		user: &registration.User{},
	}
}
func NewExistingUserBuilder(user registration.User) *UserBuilder {
	return &UserBuilder{
		user: &user,
	}
}
func (builder *UserBuilder) Build() *registration.User {
	return builder.user
}
func (builder *UserBuilder) Details() *UserDetailsBuilder {
	return &UserDetailsBuilder{
		*builder,
	}
}

func (builder *UserBuilder) ResidenceDetails() *ResidenceDetailsBuilder {
	return &ResidenceDetailsBuilder{
		*builder,
	}
}

func (builder *UserBuilder) Income() *UserIncomeBuilder {
	return &UserIncomeBuilder{
		*builder,
	}
}

func (builder *UserBuilder) Employment() *UserEmploymentBuilder {
	return &UserEmploymentBuilder{
		*builder,
	}
}

func (builder *UserBuilder) Identification() *UserIdentificationBuilder {
	return &UserIdentificationBuilder{
		*builder,
	}
}

func (builder *UserBuilder) Contact() *UserContactBuilder {
	return &UserContactBuilder{
		*builder,
	}
}

type UserDetailsBuilder struct {
	UserBuilder
}

func (b *UserDetailsBuilder) Are(title string, dob string, maritalStatus string, gender string) *UserDetailsBuilder {
	b.user.UserDetails.Title = title
	b.user.UserDetails.DOB = dob
	b.user.UserDetails.MaritalStatus = maritalStatus
	b.user.UserDetails.Gender = gender
	return b
}

type ResidenceDetailsBuilder struct {
	UserBuilder
}

func (b *ResidenceDetailsBuilder) Are(proofOfResidence string, countryOfBirth int, nationality int, residenceStatus string) *ResidenceDetailsBuilder {
	b.user.ResidenceDetails.CreatedAt = time.Now()
	b.user.ResidenceDetails.ProofOfResidency = proofOfResidence
	b.user.ResidenceDetails.CountryOfBirthID = countryOfBirth
	b.user.ResidenceDetails.NationalityID = nationality
	b.user.ResidenceDetails.ResidentialStatus = residenceStatus
	b.user.ResidenceDetails.UpdatedAt = time.Now()
	return b
}

type UserIdentificationBuilder struct {
	UserBuilder
}

func (b *UserIdentificationBuilder) Is(idNumber, typeOfId string, issuingDate, expiryDate time.Time) *UserIdentificationBuilder {
	b.user.UserIdentification.CreatedAt = time.Now()
	b.user.UserIdentification.UpdatedAt = time.Now()
	b.user.UserIdentification.ExpiryDate = expiryDate
	b.user.UserIdentification.IssueDate = issuingDate
	b.user.UserIdentification.IDNumber = idNumber
	b.user.UserIdentification.TypeOfID = typeOfId
	return b
}

type UserIncomeBuilder struct {
	UserBuilder
}

func (b *UserIncomeBuilder) Is(monthlyIncome decimal.Decimal, source string, proofOfSource string) *UserIncomeBuilder {
	b.user.UserIncome.MonthlyIncome = monthlyIncome
	b.user.UserIncome.ProofOfSource = proofOfSource
	b.user.UserIncome.CreatedAt = time.Now()
	b.user.UserIncome.UpdatedAt = time.Now()
	b.user.UserIncome.Source = source
	return b
}

type UserEmploymentBuilder struct {
	UserBuilder
}

func (b *UserEmploymentBuilder) At(employer string, industry string) *UserEmploymentBuilder {
	b.user.UserEmployment.CreatedAt = time.Now()
	b.user.UserEmployment.UpdatedAt = time.Now()
	b.user.UserEmployment.NameOfEmployer = employer
	b.user.UserEmployment.Industry = industry
	return b
}

type UserContactBuilder struct {
	UserBuilder
}

func (b *UserContactBuilder) Add(typeId int, value string) *UserContactBuilder {
	b.user.Contacts = append(b.user.Contacts, registration.Contact{
		TypeID: typeId,
		Value:  value,
	})
	return b
}
