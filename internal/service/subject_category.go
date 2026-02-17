package service

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/handler"
)

// SubjectCategoryRepository Serviceが必要とするRepositoryのインターフェース
type SubjectCategoryRepository interface {
	List(ctx context.Context) ([]domain.SubjectCategory, error)
	Detail(ctx context.Context, id string) (*domain.SubjectCategory, error)
	Create(ctx context.Context, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error)
	Update(ctx context.Context, id string, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error)
	Delete(ctx context.Context, id string) error
}

type subjectCategoryService struct {
	repo SubjectCategoryRepository
}

// NewSubjectCategoryService 新規作成する
func NewSubjectCategoryService(repo SubjectCategoryRepository) handler.SubjectCategoryService {
	return &subjectCategoryService{
		repo: repo,
	}
}

// List 一覧を取得する
func (s *subjectCategoryService) List(ctx context.Context) ([]domain.SubjectCategory, error) {
	return s.repo.List(ctx)
}

// Detail 詳細を取得する
func (s *subjectCategoryService) Detail(ctx context.Context, id string) (*domain.SubjectCategory, error) {
	return s.repo.Detail(ctx, id)
}

// Create 新規作成する
func (s *subjectCategoryService) Create(ctx context.Context, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error) {
	return s.repo.Create(ctx, req)
}

// Update 更新する
func (s *subjectCategoryService) Update(ctx context.Context, id string, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error) {
	return s.repo.Update(ctx, id, req)
}

// Delete 削除する
func (s *subjectCategoryService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
