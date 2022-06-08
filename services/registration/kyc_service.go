package registration

import (
	"context"
	registration_dtos "github.com/novabankapp/usermanagement.application/dtos/registration"
	baseService "github.com/novabankapp/usermanagement.application/services/base"
	accDomain "github.com/novabankapp/usermanagement.data/domain/account"
	regDomain "github.com/novabankapp/usermanagement.data/domain/registration"
	noSql "github.com/novabankapp/usermanagement.data/repositories/base/cassandra"
)

type KycService interface {
	SaveUserDetails(ctx context.Context, details registration_dtos.UserDetailsDto) (bool, error)
}
type KycRepositories struct {
	contactRepo            baseService.RdbmsService[regDomain.Contact]
	userDetailsRepo        baseService.RdbmsService[regDomain.UserDetails]
	residenceDetailsRepo   baseService.RdbmsService[regDomain.ResidenceDetails]
	userIdentificationRepo baseService.RdbmsService[regDomain.UserIdentification]
	userIncomeRepo         baseService.RdbmsService[regDomain.UserIncome]
	userEmploymentRepo     baseService.RdbmsService[regDomain.UserEmployment]
	kycRepo                baseService.NoSqlService[accDomain.KycCompliant]
	accountRepo            baseService.NoSqlService[accDomain.UserAccount]
	accountActivityRepo    baseService.NoSqlService[accDomain.UserAccountActivity]
}
type kycService struct {
	kycRepos KycRepositories
}

func NewKycService(kycRepos KycRepositories) KycService {
	return &kycService{
		kycRepos: kycRepos,
	}
}
func (k kycService) SaveUserDetails(ctx context.Context, details registration_dtos.UserDetailsDto) (bool, error) {
	res, _ := k.kycRepos.userDetailsRepo.GetByCondition(ctx, &regDomain.UserDetails{
		UserID: details.UserID,
	})
	if res == nil {
		_, error := k.kycRepos.userDetailsRepo.Create(ctx, regDomain.UserDetails{
			UserID:        details.UserID,
			DOB:           details.DOB,
			Title:         details.Title,
			MaritalStatus: details.MaritalStatus,
			Gender:        details.Gender,
		})
		if error != nil {
			return false, error
		}
	} else {
		result := *res
		_, error := k.kycRepos.userDetailsRepo.Update(ctx, regDomain.UserDetails{
			UserID:        details.UserID,
			DOB:           details.DOB,
			Title:         details.Title,
			MaritalStatus: details.MaritalStatus,
			Gender:        details.Gender,
		}, result.ID)
		if error != nil {
			return false, error
		}
	}
	queries := make([]map[string]string, 1)
	queries = makeQueries(queries, "UserId", "=", details.UserID)

	r, _ := k.kycRepos.kycRepo.GetByCondition(ctx, queries)
	if r == nil {
		_, _ = k.kycRepos.kycRepo.Create(ctx, accDomain.KycCompliant{
			UserId:         details.UserID,
			HasUserDetails: true,
		})
	} else {
		re := *r
		k.kycRepos.kycRepo.Update(ctx, accDomain.KycCompliant{
			HasUserDetails: true,
		}, re.ID)
	}
	return true, nil
}
func (k kycService) SaveUserIncome(ctx context.Context, details registration_dtos.UserIncomeDto) (bool, error) {
	res, _ := k.kycRepos.userIncomeRepo.GetByCondition(ctx, &regDomain.UserIncome{
		UserID: details.UserID,
	})
	if res == nil {
		_, error := k.kycRepos.userIncomeRepo.Create(ctx, regDomain.UserIncome{
			UserID:        details.UserID,
			Source:        details.Source,
			MonthlyIncome: details.MonthlyIncome,
			ProofOfSource: details.ProofOfSource,
		})
		if error != nil {
			return false, error
		}
	} else {
		result := *res
		_, error := k.kycRepos.userIncomeRepo.Update(ctx, regDomain.UserIncome{
			UserID:        details.UserID,
			Source:        details.Source,
			MonthlyIncome: details.MonthlyIncome,
			ProofOfSource: details.ProofOfSource,
		}, result.ID)
		if error != nil {
			return false, error
		}
	}
	queries := make([]map[string]string, 1)
	queries = makeQueries(queries, "UserId", "=", details.UserID)

	r, _ := k.kycRepos.kycRepo.GetByCondition(ctx, queries)
	if r == nil {
		_, _ = k.kycRepos.kycRepo.Create(ctx, accDomain.KycCompliant{
			UserId:        details.UserID,
			HasUserIncome: true,
		})
	} else {
		re := *r
		k.kycRepos.kycRepo.Update(ctx, accDomain.KycCompliant{
			HasUserIncome: true,
		}, re.ID)
	}
	return true, nil
}
func (k kycService) SaveUserIdentification(ctx context.Context, details registration_dtos.UserIdentificationDto) (bool, error) {
	res, _ := k.kycRepos.userIdentificationRepo.GetByCondition(ctx, &regDomain.UserIdentification{
		UserID: details.UserID,
	})
	if res == nil {
		_, error := k.kycRepos.userIdentificationRepo.Create(ctx, regDomain.UserIdentification{
			UserID:     details.UserID,
			TypeOfID:   details.TypeOfID,
			IDNumber:   details.IDNumber,
			IssueDate:  details.IssueDate,
			ExpiryDate: details.ExpiryDate,
		})
		if error != nil {
			return false, error
		}
	} else {
		result := *res
		_, error := k.kycRepos.userIdentificationRepo.Update(ctx, regDomain.UserIdentification{
			UserID:     details.UserID,
			TypeOfID:   details.TypeOfID,
			IDNumber:   details.IDNumber,
			IssueDate:  details.IssueDate,
			ExpiryDate: details.ExpiryDate,
		}, result.ID)
		if error != nil {
			return false, error
		}
	}
	queries := make([]map[string]string, 1)
	queries = makeQueries(queries, "UserId", "=", details.UserID)

	r, _ := k.kycRepos.kycRepo.GetByCondition(ctx, queries)
	if r == nil {
		_, _ = k.kycRepos.kycRepo.Create(ctx, accDomain.KycCompliant{
			UserId:                details.UserID,
			HasUserIdentification: true,
		})
	} else {
		re := *r
		k.kycRepos.kycRepo.Update(ctx, accDomain.KycCompliant{
			HasUserIdentification: true,
		}, re.ID)
	}
	return true, nil
}
func (k kycService) SaveResidenceDetails(ctx context.Context, details registration_dtos.ResidenceDetailsDto) (bool, error) {
	res, _ := k.kycRepos.residenceDetailsRepo.GetByCondition(ctx, &regDomain.ResidenceDetails{
		UserID: details.UserID,
	})
	if res == nil {
		_, error := k.kycRepos.residenceDetailsRepo.Create(ctx, regDomain.ResidenceDetails{
			UserID:            details.UserID,
			ResidentialStatus: details.ResidentialStatus,
			ProofOfResidency:  details.ProofOfResidency,
			NationalityID:     details.NationalityCountryID,
			CountryOfBirthID:  details.CountryOfBirthID,
		})
		if error != nil {
			return false, error
		}
	} else {
		result := *res
		_, error := k.kycRepos.residenceDetailsRepo.Update(ctx, regDomain.ResidenceDetails{
			UserID:            details.UserID,
			ResidentialStatus: details.ResidentialStatus,
			ProofOfResidency:  details.ProofOfResidency,
			NationalityID:     details.NationalityCountryID,
			CountryOfBirthID:  details.CountryOfBirthID,
		}, result.ID)
		if error != nil {
			return false, error
		}
	}
	queries := make([]map[string]string, 1)
	queries = makeQueries(queries, "UserId", "=", details.UserID)

	r, _ := k.kycRepos.kycRepo.GetByCondition(ctx, queries)
	if r == nil {
		_, _ = k.kycRepos.kycRepo.Create(ctx, accDomain.KycCompliant{
			UserId:              details.UserID,
			HasResidenceDetails: true,
		})
	} else {
		re := *r
		k.kycRepos.kycRepo.Update(ctx, accDomain.KycCompliant{
			HasResidenceDetails: true,
		}, re.ID)
	}
	return true, nil
}
func (k kycService) SaveUserEmployment(ctx context.Context, details registration_dtos.UserEmploymentDto) (bool, error) {
	res, _ := k.kycRepos.userEmploymentRepo.GetByCondition(ctx, &regDomain.UserEmployment{
		UserID: details.UserID,
	})
	if res == nil {
		_, error := k.kycRepos.userEmploymentRepo.Create(ctx, regDomain.UserEmployment{
			UserID:         details.UserID,
			NameOfEmployer: details.NameOfEmployer,
			Industry:       details.Industry,
		})
		if error != nil {
			return false, error
		}
	} else {
		result := *res
		_, error := k.kycRepos.userEmploymentRepo.Update(ctx, regDomain.UserEmployment{
			UserID:         details.UserID,
			NameOfEmployer: details.NameOfEmployer,
			Industry:       details.Industry,
		}, result.ID)
		if error != nil {
			return false, error
		}
	}
	queries := make([]map[string]string, 1)
	queries = makeQueries(queries, "UserId", "=", details.UserID)

	r, _ := k.kycRepos.kycRepo.GetByCondition(ctx, queries)
	if r == nil {
		_, _ = k.kycRepos.kycRepo.Create(ctx, accDomain.KycCompliant{
			UserId:            details.UserID,
			HasUserEmployment: true,
		})
	} else {
		re := *r
		k.kycRepos.kycRepo.Update(ctx, accDomain.KycCompliant{
			HasUserEmployment: true,
		}, re.ID)
	}
	return true, nil
}

func (k kycService) SaveUserContact(ctx context.Context, details registration_dtos.ContactDto) (bool, error) {
	res, _ := k.kycRepos.contactRepo.GetByCondition(ctx, &regDomain.Contact{
		UserID: details.UserID,
	})
	if res == nil {
		_, error := k.kycRepos.contactRepo.Create(ctx, regDomain.Contact{
			UserID: details.UserID,
			TypeID: details.TypeID,
			Value:  details.Value,
		})
		if error != nil {
			return false, error
		}
	} else {
		result := *res
		_, error := k.kycRepos.contactRepo.Update(ctx, regDomain.Contact{
			UserID: details.UserID,
			TypeID: details.TypeID,
			Value:  details.Value,
		}, result.ID)
		if error != nil {
			return false, error
		}
	}
	queries := make([]map[string]string, 1)
	queries = makeQueries(queries, "UserId", "=", details.UserID)

	r, _ := k.kycRepos.kycRepo.GetByCondition(ctx, queries)
	if r == nil {
		_, _ = k.kycRepos.kycRepo.Create(ctx, accDomain.KycCompliant{
			UserId: details.UserID,
			//HasUserIncome: true,
		})
	} else {
		re := *r
		k.kycRepos.kycRepo.Update(ctx, accDomain.KycCompliant{
			//HasUserEmployment: true,
		}, re.ID)
	}
	return true, nil
}

func makeQueries(queries []map[string]string, field, compare, value string) []map[string]string {

	m := make(map[string]string)
	m[noSql.COLUMN] = field
	m[noSql.COMPARE] = compare
	m[noSql.VALUE] = value
	queries = append(queries, m)
	return queries
}
