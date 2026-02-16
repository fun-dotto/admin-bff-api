package service

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/handler"
)

// FacultyRepository Serviceが必要とするRepositoryのインターフェース
type FacultyRepository interface {
	List(ctx context.Context) ([]domain.Faculty, error)
	Detail(ctx context.Context, id string) (*domain.Faculty, error)
	Create(ctx context.Context, req *domain.FacultyRequest) (*domain.Faculty, error)
	Update(ctx context.Context, id string, req *domain.FacultyRequest) (*domain.Faculty, error)
	Delete(ctx context.Context, id string) error
}

type facultyService struct {
	repo FacultyRepository
}

// NewFacultyService 新規作成する
func NewFacultyService(repo FacultyRepository) handler.FacultyService {
	return &facultyService{
		repo: repo,
	}
}

// List 一覧を取得する
func (s *facultyService) List(ctx context.Context) ([]domain.Faculty, error) {
	return s.repo.List(ctx)
}

// Detail 詳細を取得する
func (s *facultyService) Detail(ctx context.Context, id string) (*domain.Faculty, error) {
	return s.repo.Detail(ctx, id)
}

// Create 新規作成する
func (s *facultyService) Create(ctx context.Context, req *domain.FacultyRequest) (*domain.Faculty, error) {
	return s.repo.Create(ctx, req)
}

// Update 更新する
func (s *facultyService) Update(ctx context.Context, id string, req *domain.FacultyRequest) (*domain.Faculty, error) {
	return s.repo.Update(ctx, id, req)
}

// Delete 削除する
func (s *facultyService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
