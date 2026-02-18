package repository

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/service"
)

type MockFacultyRepository struct {
	ListFunc   func(ctx context.Context) ([]domain.Faculty, error)
	DetailFunc func(ctx context.Context, id string) (*domain.Faculty, error)
	CreateFunc func(ctx context.Context, req *domain.FacultyRequest) (*domain.Faculty, error)
	UpdateFunc func(ctx context.Context, id string, req *domain.FacultyRequest) (*domain.Faculty, error)
	DeleteFunc func(ctx context.Context, id string) error
}

func NewMockFacultyRepository() service.FacultyRepository {
	return &MockFacultyRepository{}
}

func (r *MockFacultyRepository) List(ctx context.Context) ([]domain.Faculty, error) {
	if r.ListFunc != nil {
		return r.ListFunc(ctx)
	}
	return []domain.Faculty{}, nil
}

func (r *MockFacultyRepository) Detail(ctx context.Context, id string) (*domain.Faculty, error) {
	if r.DetailFunc != nil {
		return r.DetailFunc(ctx, id)
	}
	return nil, nil
}

func (r *MockFacultyRepository) Create(ctx context.Context, req *domain.FacultyRequest) (*domain.Faculty, error) {
	if r.CreateFunc != nil {
		return r.CreateFunc(ctx, req)
	}
	return nil, nil
}

func (r *MockFacultyRepository) Update(ctx context.Context, id string, req *domain.FacultyRequest) (*domain.Faculty, error) {
	if r.UpdateFunc != nil {
		return r.UpdateFunc(ctx, id, req)
	}
	return nil, nil
}

func (r *MockFacultyRepository) Delete(ctx context.Context, id string) error {
	if r.DeleteFunc != nil {
		return r.DeleteFunc(ctx, id)
	}
	return nil
}
