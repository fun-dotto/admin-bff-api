package repository

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/service"
)

type MockSubjectCategoryRepository struct {
	ListFunc   func(ctx context.Context) ([]domain.SubjectCategory, error)
	DetailFunc func(ctx context.Context, id string) (*domain.SubjectCategory, error)
	CreateFunc func(ctx context.Context, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error)
	UpdateFunc func(ctx context.Context, id string, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error)
	DeleteFunc func(ctx context.Context, id string) error
}

func NewMockSubjectCategoryRepository() service.SubjectCategoryRepository {
	return &MockSubjectCategoryRepository{}
}

func (r *MockSubjectCategoryRepository) List(ctx context.Context) ([]domain.SubjectCategory, error) {
	if r.ListFunc != nil {
		return r.ListFunc(ctx)
	}
	return []domain.SubjectCategory{}, nil
}

func (r *MockSubjectCategoryRepository) Detail(ctx context.Context, id string) (*domain.SubjectCategory, error) {
	if r.DetailFunc != nil {
		return r.DetailFunc(ctx, id)
	}
	return nil, nil
}

func (r *MockSubjectCategoryRepository) Create(ctx context.Context, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error) {
	if r.CreateFunc != nil {
		return r.CreateFunc(ctx, req)
	}
	return nil, nil
}

func (r *MockSubjectCategoryRepository) Update(ctx context.Context, id string, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error) {
	if r.UpdateFunc != nil {
		return r.UpdateFunc(ctx, id, req)
	}
	return nil, nil
}

func (r *MockSubjectCategoryRepository) Delete(ctx context.Context, id string) error {
	if r.DeleteFunc != nil {
		return r.DeleteFunc(ctx, id)
	}
	return nil
}
