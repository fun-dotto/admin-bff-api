package service

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/handler"
)

// CourseRepository Serviceが必要とするRepositoryのインターフェース
type CourseRepository interface {
	List(ctx context.Context) ([]domain.Course, error)
	Detail(ctx context.Context, id string) (*domain.Course, error)
	Create(ctx context.Context, req *domain.CourseRequest) (*domain.Course, error)
	Update(ctx context.Context, id string, req *domain.CourseRequest) (*domain.Course, error)
	Delete(ctx context.Context, id string) error
}

type courseService struct {
	repo CourseRepository
}

// NewCourseService 新規作成する
func NewCourseService(repo CourseRepository) handler.CourseService {
	return &courseService{
		repo: repo,
	}
}

// List 一覧を取得する
func (s *courseService) List(ctx context.Context) ([]domain.Course, error) {
	return s.repo.List(ctx)
}

// Detail 詳細を取得する
func (s *courseService) Detail(ctx context.Context, id string) (*domain.Course, error) {
	return s.repo.Detail(ctx, id)
}

// Create 新規作成する
func (s *courseService) Create(ctx context.Context, req *domain.CourseRequest) (*domain.Course, error) {
	return s.repo.Create(ctx, req)
}

// Update 更新する
func (s *courseService) Update(ctx context.Context, id string, req *domain.CourseRequest) (*domain.Course, error) {
	return s.repo.Update(ctx, id, req)
}

// Delete 削除する
func (s *courseService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
