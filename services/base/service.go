package base

import (
	"context"
	domainBase "github.com/novabankapp/usermanagement.data/domain/base"
	baseRepository "github.com/novabankapp/usermanagement.data/repositories/base/postgres"
)

type Service[E domainBase.Entity] struct {
	repo baseRepository.PostgresRepository[E]
}

func NewService[E domainBase.Entity](repo baseRepository.PostgresRepository[E]) Service[E] {
	return Service[E]{
		repo: repo,
	}
}
func (s *Service[E]) Create(ctx context.Context, entity E) (*E, error) {
	return s.repo.Create(ctx, entity)

}
func (s *Service[E]) GetById(ctx context.Context, id string) (*E, error) {
	return s.repo.GetById(ctx, id)
}

func (s *Service[E]) Update(ctx context.Context, entity E, id string) (bool, error) {
	return s.repo.Update(ctx, entity, id)
}
func (s *Service[E]) Delete(ctx context.Context, id string) (bool, error) {
	return s.repo.Delete(ctx, id)
}
func (s *Service[E]) Get(ctx context.Context,
	page int, pageSize int, query *E, orderBy string) (*[]E, error) {
	return s.repo.Get(ctx, page, pageSize, query, orderBy)

}
