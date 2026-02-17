package service

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/handler"
)

// SubjectRepository Serviceが必要とするRepositoryのインターフェース
type SubjectRepository interface {
	List(ctx context.Context) ([]domain.Subject, error)
	Detail(ctx context.Context, id string) (*domain.Subject, error)
	Create(ctx context.Context, req *domain.SubjectRequest) (*domain.Subject, error)
	Update(ctx context.Context, id string, req *domain.SubjectRequest) (*domain.Subject, error)
	Delete(ctx context.Context, id string) error
}

type subjectService struct {
	repo SubjectRepository
}

// NewSubjectService 新規作成する
func NewSubjectService(repo SubjectRepository) handler.SubjectService {
	return &subjectService{
		repo: repo,
	}
}

// List 一覧を取得する
func (s *subjectService) List(ctx context.Context) ([]domain.Subject, error) {
	return s.repo.List(ctx)
}

// Detail 詳細を取得する
func (s *subjectService) Detail(ctx context.Context, id string) (*domain.Subject, error) {
	return s.repo.Detail(ctx, id)
}

// Create 新規作成する
func (s *subjectService) Create(ctx context.Context, req *domain.SubjectRequest) (*domain.Subject, error) {
	return s.repo.Create(ctx, req)
}

// Update 更新する
func (s *subjectService) Update(ctx context.Context, id string, req *domain.SubjectRequest) (*domain.Subject, error) {
	return s.repo.Update(ctx, id, req)
}

// Delete 削除する
func (s *subjectService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
