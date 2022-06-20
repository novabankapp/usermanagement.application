package registration

import (
	"bytes"
	"context"
	"encoding/json"
	baseService "github.com/novabankapp/common.application/services/base"
	"github.com/novabankapp/common.application/services/message_queue"
	"github.com/novabankapp/common.application/utilities"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	registrationDtos "github.com/novabankapp/usermanagement.application/dtos/registration"
	accDomain "github.com/novabankapp/usermanagement.data/domain/account"
	loginDomain "github.com/novabankapp/usermanagement.data/domain/login"
	regDomain "github.com/novabankapp/usermanagement.data/domain/registration"

	"strconv"
)

type KycService interface {
	SaveUserDetails(ctx context.Context, details registrationDtos.UserDetailsDto) (bool, error)
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
	loginRepo              baseService.NoSqlService[loginDomain.UserLogin]
	accountActivityRepo    baseService.NoSqlService[accDomain.UserAccountActivity]
}
type kycService struct {
	kycRepos     KycRepositories
	messageQueue message_queue.MessageQueue
	topics       *kafkaClient.KafkaTopics
}

func NewKycService(kycRepos KycRepositories, messageQueue message_queue.MessageQueue, topics *kafkaClient.KafkaTopics) KycService {
	return &kycService{
		kycRepos:     kycRepos,
		messageQueue: messageQueue,
		topics:       topics,
	}
}
func (k kycService) DeleteContact(ctx context.Context, id uint) (bool, error) {

	_, err := k.kycRepos.contactRepo.Delete(ctx, id)
	if err != nil {
		return false, err
	}
	res := new(bytes.Buffer)
	e := json.NewEncoder(res).Encode(id)
	if e == nil {
		msgBytes := res.Bytes()
		_, _ = k.messageQueue.PublishMessage(ctx, msgBytes, strconv.FormatUint(uint64(id), 10), k.topics.ContactDeleted.TopicName)
	}
	return true, nil
}
func (k kycService) SaveUserDetails(ctx context.Context, details registrationDtos.UserDetailsDto) (bool, error) {
	res, _ := k.kycRepos.userDetailsRepo.GetByCondition(ctx, &regDomain.UserDetails{
		UserID: details.UserID,
	})
	var results regDomain.UserDetails
	if res == nil {
		r, err := k.kycRepos.userDetailsRepo.Create(ctx, regDomain.UserDetails{
			UserID:        details.UserID,
			DOB:           details.DOB,
			Title:         details.Title,
			MaritalStatus: details.MaritalStatus,
			Gender:        details.Gender,
		})

		if err != nil {
			return false, err
		}
		results = *r
	} else {
		result := *res
		data := regDomain.UserDetails{
			UserID:        details.UserID,
			DOB:           details.DOB,
			Title:         details.Title,
			MaritalStatus: details.MaritalStatus,
			Gender:        details.Gender,
		}
		_, err := k.kycRepos.userDetailsRepo.Update(ctx, data, result.ID)
		if err != nil {
			return false, err
		}

		results = data
	}
	queries := make([]map[string]string, 1)
	queries = utilities.MakeQueries(queries, "UserId", "=", details.UserID)

	r, _ := k.kycRepos.kycRepo.GetByCondition(ctx, queries)
	if r == nil {
		_, _ = k.kycRepos.kycRepo.Create(ctx, accDomain.KycCompliant{
			UserId:         details.UserID,
			HasUserDetails: true,
		})
	} else {
		re := *r
		re.HasUserDetails = true
		_, _ = k.kycRepos.kycRepo.Update(ctx, re, re.ID)
	}
	val := new(bytes.Buffer)
	e := json.NewEncoder(val).Encode(results)
	if e == nil {
		msgBytes := val.Bytes()
		_, _ = k.messageQueue.PublishMessage(ctx, msgBytes, details.UserID, k.topics.UserUpdated.TopicName)
	}
	return true, nil
}
func (k kycService) SaveUserIncome(ctx context.Context, details registrationDtos.UserIncomeDto) (bool, error) {
	res, _ := k.kycRepos.userIncomeRepo.GetByCondition(ctx, &regDomain.UserIncome{
		UserID: details.UserID,
	})
	var results regDomain.UserIncome
	if res == nil {
		r, err := k.kycRepos.userIncomeRepo.Create(ctx, regDomain.UserIncome{
			UserID:        details.UserID,
			Source:        details.Source,
			MonthlyIncome: details.MonthlyIncome,
			ProofOfSource: details.ProofOfSource,
		})
		if err != nil {
			return false, err
		}
		results = *r
	} else {
		result := *res
		data := regDomain.UserIncome{
			UserID:        details.UserID,
			Source:        details.Source,
			MonthlyIncome: details.MonthlyIncome,
			ProofOfSource: details.ProofOfSource,
		}
		_, err := k.kycRepos.userIncomeRepo.Update(ctx, data, result.ID)
		if err != nil {
			return false, err
		}
		results = data
	}
	queries := make([]map[string]string, 1)
	queries = utilities.MakeQueries(queries, "UserId", "=", details.UserID)

	r, _ := k.kycRepos.kycRepo.GetByCondition(ctx, queries)
	if r == nil {
		_, _ = k.kycRepos.kycRepo.Create(ctx, accDomain.KycCompliant{
			UserId:        details.UserID,
			HasUserIncome: true,
		})
	} else {
		re := *r
		re.HasUserIncome = true
		_, _ = k.kycRepos.kycRepo.Update(ctx, re, re.ID)
	}
	val := new(bytes.Buffer)
	e := json.NewEncoder(val).Encode(results)
	if e == nil {
		msgBytes := val.Bytes()
		_, _ = k.messageQueue.PublishMessage(ctx, msgBytes, details.UserID, k.topics.UserUpdated.TopicName)
	}
	return true, nil
}
func (k kycService) SaveUserIdentification(ctx context.Context, details registrationDtos.UserIdentificationDto) (bool, error) {
	res, _ := k.kycRepos.userIdentificationRepo.GetByCondition(ctx, &regDomain.UserIdentification{
		UserID: details.UserID,
	})
	var results regDomain.UserIdentification
	if res == nil {
		r, err := k.kycRepos.userIdentificationRepo.Create(ctx, regDomain.UserIdentification{
			UserID:     details.UserID,
			TypeOfID:   details.TypeOfID,
			IDNumber:   details.IDNumber,
			IssueDate:  details.IssueDate,
			ExpiryDate: details.ExpiryDate,
		})
		if err != nil {
			return false, err
		}
		results = *r
	} else {
		result := *res
		data := regDomain.UserIdentification{
			UserID:     details.UserID,
			TypeOfID:   details.TypeOfID,
			IDNumber:   details.IDNumber,
			IssueDate:  details.IssueDate,
			ExpiryDate: details.ExpiryDate,
		}
		_, err := k.kycRepos.userIdentificationRepo.Update(ctx, data, result.ID)
		if err != nil {
			return false, err
		}
		results = data
	}
	queries := make([]map[string]string, 1)
	queries = utilities.MakeQueries(queries, "UserId", "=", details.UserID)

	r, _ := k.kycRepos.kycRepo.GetByCondition(ctx, queries)
	if r == nil {
		_, _ = k.kycRepos.kycRepo.Create(ctx, accDomain.KycCompliant{
			UserId:                details.UserID,
			HasUserIdentification: true,
		})
	} else {
		re := *r
		re.HasUserIdentification = true
		_, _ = k.kycRepos.kycRepo.Update(ctx, re, re.ID)
	}
	val := new(bytes.Buffer)
	e := json.NewEncoder(val).Encode(results)
	if e == nil {
		msgBytes := val.Bytes()
		_, _ = k.messageQueue.PublishMessage(ctx, msgBytes, details.UserID, k.topics.UserUpdated.TopicName)
	}
	return true, nil
}
func (k kycService) SaveResidenceDetails(ctx context.Context, details registrationDtos.ResidenceDetailsDto) (bool, error) {
	res, _ := k.kycRepos.residenceDetailsRepo.GetByCondition(ctx, &regDomain.ResidenceDetails{
		UserID: details.UserID,
	})
	var results regDomain.ResidenceDetails
	if res == nil {
		r, err := k.kycRepos.residenceDetailsRepo.Create(ctx, regDomain.ResidenceDetails{
			UserID:            details.UserID,
			ResidentialStatus: details.ResidentialStatus,
			ProofOfResidency:  details.ProofOfResidency,
			NationalityID:     details.NationalityCountryID,
			CountryOfBirthID:  details.CountryOfBirthID,
		})
		if err != nil {
			return false, err
		}
		results = *r
	} else {
		result := *res
		data := regDomain.ResidenceDetails{
			UserID:            details.UserID,
			ResidentialStatus: details.ResidentialStatus,
			ProofOfResidency:  details.ProofOfResidency,
			NationalityID:     details.NationalityCountryID,
			CountryOfBirthID:  details.CountryOfBirthID,
		}
		_, err := k.kycRepos.residenceDetailsRepo.Update(ctx, data, result.ID)
		if err != nil {
			return false, err
		}
		results = data
	}
	queries := make([]map[string]string, 1)
	queries = utilities.MakeQueries(queries, "UserId", "=", details.UserID)

	r, _ := k.kycRepos.kycRepo.GetByCondition(ctx, queries)
	if r == nil {
		_, _ = k.kycRepos.kycRepo.Create(ctx, accDomain.KycCompliant{
			UserId:              details.UserID,
			HasResidenceDetails: true,
		})
	} else {
		re := *r
		re.HasResidenceDetails = true
		_, _ = k.kycRepos.kycRepo.Update(ctx, re, re.ID)
	}
	val := new(bytes.Buffer)
	e := json.NewEncoder(val).Encode(results)
	if e == nil {
		msgBytes := val.Bytes()
		_, _ = k.messageQueue.PublishMessage(ctx, msgBytes, details.UserID, k.topics.UserUpdated.TopicName)
	}
	return true, nil
}
func (k kycService) SaveUserEmployment(ctx context.Context, details registrationDtos.UserEmploymentDto) (bool, error) {
	res, _ := k.kycRepos.userEmploymentRepo.GetByCondition(ctx, &regDomain.UserEmployment{
		UserID: details.UserID,
	})
	var results regDomain.UserEmployment
	if res == nil {
		r, err := k.kycRepos.userEmploymentRepo.Create(ctx, regDomain.UserEmployment{
			UserID:         details.UserID,
			NameOfEmployer: details.NameOfEmployer,
			Industry:       details.Industry,
		})
		if err != nil {
			return false, err
		}
		results = *r
	} else {
		result := *res
		data := regDomain.UserEmployment{
			UserID:         details.UserID,
			NameOfEmployer: details.NameOfEmployer,
			Industry:       details.Industry,
		}
		_, err := k.kycRepos.userEmploymentRepo.Update(ctx, data, result.ID)
		if err != nil {
			return false, err
		}
		results = data
	}
	queries := make([]map[string]string, 1)
	queries = utilities.MakeQueries(queries, "UserId", "=", details.UserID)

	r, _ := k.kycRepos.kycRepo.GetByCondition(ctx, queries)
	if r == nil {
		_, _ = k.kycRepos.kycRepo.Create(ctx, accDomain.KycCompliant{
			UserId:            details.UserID,
			HasUserEmployment: true,
		})
	} else {
		re := *r
		re.HasUserEmployment = true
		_, _ = k.kycRepos.kycRepo.Update(ctx, re, re.ID)
	}
	val := new(bytes.Buffer)
	e := json.NewEncoder(val).Encode(results)
	if e == nil {
		msgBytes := val.Bytes()
		_, _ = k.messageQueue.PublishMessage(ctx, msgBytes, details.UserID, k.topics.UserUpdated.TopicName)
	}
	return true, nil
}

func (k kycService) SaveUserContact(ctx context.Context, details registrationDtos.ContactDto) (bool, error) {
	var results regDomain.Contact
	if details.ID == nil {

		r, err := k.kycRepos.contactRepo.Create(ctx, regDomain.Contact{
			UserID: details.UserID,
			TypeID: details.TypeID,
			Value:  details.Value,
		})
		if err != nil {
			return false, err
		}
		results = *r
	} else {
		dom := regDomain.Contact{
			UserID: details.UserID,
		}
		dom.ID = *details.ID
		res, _ := k.kycRepos.contactRepo.GetByCondition(ctx, &dom)
		result := *res
		data := regDomain.Contact{
			UserID: details.UserID,
			TypeID: details.TypeID,
			Value:  details.Value,
		}
		_, err := k.kycRepos.contactRepo.Update(ctx, data, result.ID)
		if err != nil {
			return false, err
		}
		results = data
	}
	queries := make([]map[string]string, 1)
	queries = utilities.MakeQueries(queries, "UserId", "=", details.UserID)

	r, _ := k.kycRepos.kycRepo.GetByCondition(ctx, queries)
	if r == nil {
		_, _ = k.kycRepos.kycRepo.Create(ctx, accDomain.KycCompliant{
			UserId: details.UserID,
			//HasUserIncome: true,
		})
	} else {
		re := *r
		_, _ = k.kycRepos.kycRepo.Update(ctx, accDomain.KycCompliant{
			//HasUserEmployment: true,
		}, re.ID)
	}
	val := new(bytes.Buffer)
	e := json.NewEncoder(val).Encode(results)
	if e == nil {
		msgBytes := val.Bytes()
		_, _ = k.messageQueue.PublishMessage(ctx, msgBytes, details.UserID, k.topics.ContactUpdated.TopicName)
	}
	return true, nil
}
