package repository

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/service"
)

type mockCourseRepository struct{}

// NewMockCourseRepository モックリポジトリを作成する
func NewMockCourseRepository() service.CourseRepository {
	return &mockCourseRepository{}
}

// List 一覧を取得する（モック）
func (r *mockCourseRepository) List(ctx context.Context) ([]domain.Course, error) {
	return []domain.Course{
		{
			ID:   "1",
			Name: "コース1",
		},
	}, nil
}

// Detail 詳細を取得する（モック）
func (r *mockCourseRepository) Detail(ctx context.Context, id string) (*domain.Course, error) {
	return &domain.Course{
		ID:   id,
		Name: "コース" + id,
	}, nil
}

// Create 新規作成する（モック）
func (r *mockCourseRepository) Create(ctx context.Context, req *domain.CourseRequest) (*domain.Course, error) {
	return &domain.Course{
		ID:   "created-id",
		Name: req.Name,
	}, nil
}

// Update 更新する（モック）
func (r *mockCourseRepository) Update(ctx context.Context, id string, req *domain.CourseRequest) (*domain.Course, error) {
	return &domain.Course{
		ID:   id,
		Name: req.Name,
	}, nil
}

// Delete 削除する（モック）
func (r *mockCourseRepository) Delete(ctx context.Context, id string) error {
	return nil
}
