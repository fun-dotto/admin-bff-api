package repository

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/service"
)

type MockAnnouncementRepository struct {
	ListFunc   func(ctx context.Context) ([]domain.Announcement, error)
	DetailFunc func(ctx context.Context, id string) (*domain.Announcement, error)
	CreateFunc func(ctx context.Context, req *domain.AnnouncementRequest) (*domain.Announcement, error)
	UpdateFunc func(ctx context.Context, id string, req *domain.AnnouncementRequest) (*domain.Announcement, error)
	DeleteFunc func(ctx context.Context, id string) error
}

func NewMockAnnouncementRepository() service.AnnouncementRepository {
	return &MockAnnouncementRepository{}
}

func (r *MockAnnouncementRepository) List(ctx context.Context) ([]domain.Announcement, error) {
	if r.ListFunc != nil {
		return r.ListFunc(ctx)
	}
	return []domain.Announcement{}, nil
}

func (r *MockAnnouncementRepository) Detail(ctx context.Context, id string) (*domain.Announcement, error) {
	if r.DetailFunc != nil {
		return r.DetailFunc(ctx, id)
	}
	return nil, nil
}

func (r *MockAnnouncementRepository) Create(ctx context.Context, req *domain.AnnouncementRequest) (*domain.Announcement, error) {
	if r.CreateFunc != nil {
		return r.CreateFunc(ctx, req)
	}
	return nil, nil
}

func (r *MockAnnouncementRepository) Update(ctx context.Context, id string, req *domain.AnnouncementRequest) (*domain.Announcement, error) {
	if r.UpdateFunc != nil {
		return r.UpdateFunc(ctx, id, req)
	}
	return nil, nil
}

func (r *MockAnnouncementRepository) Delete(ctx context.Context, id string) error {
	if r.DeleteFunc != nil {
		return r.DeleteFunc(ctx, id)
	}
	return nil
}
