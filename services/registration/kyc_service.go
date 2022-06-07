package registration

import (
	baseService "github.com/novabankapp/usermanagement.application/services/base"
	accDomain "github.com/novabankapp/usermanagement.data/domain/account"
	regDomain "github.com/novabankapp/usermanagement.data/domain/registration"
)

type KycService interface {
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
