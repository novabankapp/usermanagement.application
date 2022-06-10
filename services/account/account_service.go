package account

import (
	"bytes"
	"context"
	"encoding/json"
	kafkaClient "github.com/novabankapp/common.infrastructure/kafka"
	baseService "github.com/novabankapp/usermanagement.application/services/base"
	"github.com/novabankapp/usermanagement.application/services/message_queue"
	accDomain "github.com/novabankapp/usermanagement.data/domain/account"
	loginDomain "github.com/novabankapp/usermanagement.data/domain/login"
)

type AccountService interface {
	LockAccount(ctx context.Context, id string) (bool, error)
	UnlockAccount(ctx context.Context, id string) (bool, error)
	DeactivateAccount(ctx context.Context, id string) (bool, error)
	DeleteAccount(ctx context.Context, id string) (bool, error)
}

type accountService struct {
	accountRepo  baseService.NoSqlService[accDomain.UserAccount]
	loginRepo    baseService.NoSqlService[loginDomain.UserLogin]
	messageQueue message_queue.MessageQueue
	topics       *kafkaClient.KafkaTopics
}

func NewAccountService(
	accountRepo baseService.NoSqlService[accDomain.UserAccount],
	loginRepo baseService.NoSqlService[loginDomain.UserLogin],
	messageQueue message_queue.MessageQueue,
	topics *kafkaClient.KafkaTopics,
) AccountService {
	return &accountService{
		accountRepo:  accountRepo,
		loginRepo:    loginRepo,
		messageQueue: messageQueue,
		topics:       topics,
	}
}
func (a accountService) UnlockAccount(ctx context.Context, id string) (bool, error) {
	res, err := a.accountRepo.GetById(ctx, id)
	if err != nil {
		return false, err
	}
	var result = *res
	result.IsLocked = false
	edited, err := a.accountRepo.Update(ctx, result)
	if edited {
		res := new(bytes.Buffer)
		e := json.NewEncoder(res).Encode(id)
		if e == nil {
			msgBytes := res.Bytes()
			_, _ = a.messageQueue.PublishMessage(ctx, msgBytes, id,
				a.topics.AccountLocked.TopicName)
		}
	}
	return edited, err
}
func (a accountService) LockAccount(ctx context.Context, id string) (bool, error) {
	res, err := a.accountRepo.GetById(ctx, id)
	if err != nil {
		return false, err
	}
	var result = *res
	result.IsLocked = true
	edited, err := a.accountRepo.Update(ctx, result)
	if edited {
		res := new(bytes.Buffer)
		e := json.NewEncoder(res).Encode(id)
		if e == nil {
			msgBytes := res.Bytes()
			_, _ = a.messageQueue.PublishMessage(ctx, msgBytes, id,
				a.topics.AccountLocked.TopicName)
		}
	}
	return edited, err
}
func (a accountService) ActivateAccount(ctx context.Context, id string) (bool, error) {
	res, err := a.accountRepo.GetById(ctx, id)
	if err != nil {
		return false, err
	}
	var result = *res
	result.IsActive = true
	edited, err := a.accountRepo.Update(ctx, result)
	if edited {
		res := new(bytes.Buffer)
		e := json.NewEncoder(res).Encode(id)
		if e == nil {
			msgBytes := res.Bytes()
			_, _ = a.messageQueue.PublishMessage(ctx, msgBytes, id,
				a.topics.AccountActivated.TopicName)
		}
	}
	return edited, err
}
func (a accountService) DeactivateAccount(ctx context.Context, id string) (bool, error) {
	res, err := a.accountRepo.GetById(ctx, id)
	if err != nil {
		return false, err
	}
	var result = *res
	result.IsActive = false
	edited, err := a.accountRepo.Update(ctx, result)
	if edited {
		res := new(bytes.Buffer)
		e := json.NewEncoder(res).Encode(id)
		if e == nil {
			msgBytes := res.Bytes()
			_, _ = a.messageQueue.PublishMessage(ctx, msgBytes, id,
				a.topics.AccountDeactivated.TopicName)
		}
	}
	return edited, err
}

func (a accountService) DeleteAccount(ctx context.Context, id string) (bool, error) {
	//TODO implement me
	panic("implement me")
}
