package repository

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/service"
)

type mockSubjectCategoryRepository struct{}

// NewMockSubjectCategoryRepository モックリポジトリを作成する
func NewMockSubjectCategoryRepository() service.SubjectCategoryRepository {
	return &mockSubjectCategoryRepository{}
}

// List 一覧を取得する（モック）
func (r *mockSubjectCategoryRepository) List(ctx context.Context) ([]domain.SubjectCategory, error) {
	return []domain.SubjectCategory{
		{
			ID:   "1",
			Name: "カテゴリ1",
		},
	}, nil
}

// Detail 詳細を取得する（モック）
func (r *mockSubjectCategoryRepository) Detail(ctx context.Context, id string) (*domain.SubjectCategory, error) {
	return &domain.SubjectCategory{
		ID:   id,
		Name: "カテゴリ" + id,
	}, nil
}

// Create 新規作成する（モック）
func (r *mockSubjectCategoryRepository) Create(ctx context.Context, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error) {
	return &domain.SubjectCategory{
		ID:   "created-id",
		Name: req.Name,
	}, nil
}

// Update 更新する（モック）
func (r *mockSubjectCategoryRepository) Update(ctx context.Context, id string, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error) {
	return &domain.SubjectCategory{
		ID:   id,
		Name: req.Name,
	}, nil
}

// Delete 削除する（モック）
func (r *mockSubjectCategoryRepository) Delete(ctx context.Context, id string) error {
	return nil
}
