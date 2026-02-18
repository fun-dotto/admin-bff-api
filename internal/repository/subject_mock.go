package repository

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/service"
)

type MockSubjectRepository struct {
	ListFunc   func(ctx context.Context) ([]domain.Subject, error)
	DetailFunc func(ctx context.Context, id string) (*domain.Subject, error)
	CreateFunc func(ctx context.Context, req *domain.SubjectRequest) (*domain.Subject, error)
	UpdateFunc func(ctx context.Context, id string, req *domain.SubjectRequest) (*domain.Subject, error)
	DeleteFunc func(ctx context.Context, id string) error
}

func NewMockSubjectRepository() service.SubjectRepository {
	return &MockSubjectRepository{}
}

func (r *MockSubjectRepository) List(ctx context.Context) ([]domain.Subject, error) {
	if r.ListFunc != nil {
		return r.ListFunc(ctx)
	}
	return []domain.Subject{}, nil
}

func (r *MockSubjectRepository) Detail(ctx context.Context, id string) (*domain.Subject, error) {
	if r.DetailFunc != nil {
		return r.DetailFunc(ctx, id)
	}
	return nil, nil
}

func (r *MockSubjectRepository) Create(ctx context.Context, req *domain.SubjectRequest) (*domain.Subject, error) {
	if r.CreateFunc != nil {
		return r.CreateFunc(ctx, req)
	}
	return nil, nil
}

func (r *MockSubjectRepository) Update(ctx context.Context, id string, req *domain.SubjectRequest) (*domain.Subject, error) {
	if r.UpdateFunc != nil {
		return r.UpdateFunc(ctx, id, req)
	}
	return nil, nil
}

func (r *MockSubjectRepository) Delete(ctx context.Context, id string) error {
	if r.DeleteFunc != nil {
		return r.DeleteFunc(ctx, id)
	}
	return nil
}
