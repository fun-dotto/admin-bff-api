package repository

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/service"
)

type mockFacultyRepository struct{}

// NewMockFacultyRepository モックリポジトリを作成する
func NewMockFacultyRepository() service.FacultyRepository {
	return &mockFacultyRepository{}
}

// List 一覧を取得する（モック）
func (r *mockFacultyRepository) List(ctx context.Context) ([]domain.Faculty, error) {
	return []domain.Faculty{
		{
			ID:    "1",
			Name:  "教員1",
			Email: "faculty1@example.com",
		},
	}, nil
}

// Detail 詳細を取得する（モック）
func (r *mockFacultyRepository) Detail(ctx context.Context, id string) (*domain.Faculty, error) {
	return &domain.Faculty{
		ID:    id,
		Name:  "教員" + id,
		Email: "faculty" + id + "@example.com",
	}, nil
}

// Create 新規作成する（モック）
func (r *mockFacultyRepository) Create(ctx context.Context, req *domain.FacultyRequest) (*domain.Faculty, error) {
	return &domain.Faculty{
		ID:    "created-id",
		Name:  req.Name,
		Email: req.Email,
	}, nil
}

// Update 更新する（モック）
func (r *mockFacultyRepository) Update(ctx context.Context, id string, req *domain.FacultyRequest) (*domain.Faculty, error) {
	return &domain.Faculty{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}, nil
}

// Delete 削除する（モック）
func (r *mockFacultyRepository) Delete(ctx context.Context, id string) error {
	return nil
}
