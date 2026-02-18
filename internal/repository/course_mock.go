package repository

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/service"
)

type MockCourseRepository struct {
	ListFunc   func(ctx context.Context) ([]domain.Course, error)
	DetailFunc func(ctx context.Context, id string) (*domain.Course, error)
	CreateFunc func(ctx context.Context, req *domain.CourseRequest) (*domain.Course, error)
	UpdateFunc func(ctx context.Context, id string, req *domain.CourseRequest) (*domain.Course, error)
	DeleteFunc func(ctx context.Context, id string) error
}

func NewMockCourseRepository() service.CourseRepository {
	return &MockCourseRepository{}
}

func (r *MockCourseRepository) List(ctx context.Context) ([]domain.Course, error) {
	if r.ListFunc != nil {
		return r.ListFunc(ctx)
	}
	return []domain.Course{}, nil
}

func (r *MockCourseRepository) Detail(ctx context.Context, id string) (*domain.Course, error) {
	if r.DetailFunc != nil {
		return r.DetailFunc(ctx, id)
	}
	return nil, nil
}

func (r *MockCourseRepository) Create(ctx context.Context, req *domain.CourseRequest) (*domain.Course, error) {
	if r.CreateFunc != nil {
		return r.CreateFunc(ctx, req)
	}
	return nil, nil
}

func (r *MockCourseRepository) Update(ctx context.Context, id string, req *domain.CourseRequest) (*domain.Course, error) {
	if r.UpdateFunc != nil {
		return r.UpdateFunc(ctx, id, req)
	}
	return nil, nil
}

func (r *MockCourseRepository) Delete(ctx context.Context, id string) error {
	if r.DeleteFunc != nil {
		return r.DeleteFunc(ctx, id)
	}
	return nil
}
