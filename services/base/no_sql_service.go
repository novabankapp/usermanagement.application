package base

import (
	"context"
	domainBase "github.com/novabankapp/usermanagement.data/domain/base"
	baseRepository "github.com/novabankapp/usermanagement.data/repositories/base/cassandra"
)

type NoSqlService[E domainBase.NoSqlEntity] struct {
	repo baseRepository.CassandraRepository[E]
}

func NewDocumentDatabaseService[E domainBase.NoSqlEntity](repo baseRepository.CassandraRepository[E]) NoSqlService[E] {
	return NoSqlService[E]{
		repo: repo,
	}
}
func (s *NoSqlService[E]) Create(ctx context.Context, entity E) (bool, error) {
	return s.repo.Create(ctx, entity)

}
func (s *NoSqlService[E]) GetById(ctx context.Context, id string) (*E, error) {
	return s.repo.GetById(ctx, id)
}

func (s *NoSqlService[E]) Update(ctx context.Context, entity E, id string) (bool, error) {
	return s.repo.Update(ctx, entity, id)
}
func (s *NoSqlService[E]) Delete(ctx context.Context, id string) (bool, error) {
	return s.repo.Delete(ctx, id)
}
func (s *NoSqlService[E]) GetGet(ctx context.Context,
	page []byte, pageSize int, queries []map[string]string, orderBy string) (*[]E, error) {
	return s.repo.Get(ctx, page, pageSize, queries, orderBy)

}
