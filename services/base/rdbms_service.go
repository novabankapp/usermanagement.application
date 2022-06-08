package base

import (
	"context"
	domainBase "github.com/novabankapp/usermanagement.data/domain/base"
	baseRepository "github.com/novabankapp/usermanagement.data/repositories/base/postgres"
)

type RdbmsService[E domainBase.Entity] struct {
	repo baseRepository.PostgresRepository[E]
}

func NewRelationalDatabaseService[E domainBase.Entity](repo baseRepository.PostgresRepository[E]) RdbmsService[E] {
	return RdbmsService[E]{
		repo: repo,
	}
}
func (s *RdbmsService[E]) Create(ctx context.Context, entity E) (*E, error) {
	return s.repo.Create(ctx, entity)

}
func (s *RdbmsService[E]) GetById(ctx context.Context, id string) (*E, error) {
	return s.repo.GetById(ctx, id)
}

func (s *RdbmsService[E]) Update(ctx context.Context, entity E, id uint) (bool, error) {
	return s.repo.Update(ctx, entity, id)
}
func (s *RdbmsService[E]) Delete(ctx context.Context, id string) (bool, error) {
	return s.repo.Delete(ctx, id)
}
func (s *RdbmsService[E]) Get(ctx context.Context,
	page int, pageSize int, query *E, orderBy string) (*[]E, error) {
	return s.repo.Get(ctx, page, pageSize, query, orderBy)

}
func (s *RdbmsService[E]) GetByCondition(ctx context.Context, query *E) (*E, error) {
	return s.repo.GetByCondition(ctx, query)

}
