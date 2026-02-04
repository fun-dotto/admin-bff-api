package service

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/handler"
)

type AnnouncementRepository interface {
	List(ctx context.Context) ([]domain.Announcement, error)
	Create(ctx context.Context, req *domain.AnnouncementRequest) (*domain.Announcement, error)
	Update(ctx context.Context, id string, req *domain.AnnouncementRequest) (*domain.Announcement, error)
	Delete(ctx context.Context, id string) error
}

type announcementService struct {
	repo AnnouncementRepository
}

// NewAnnouncementService 新規作成する
func NewAnnouncementService(repo AnnouncementRepository) handler.AnnouncementService {
	return &announcementService{
		repo: repo,
	}
}

// List 一覧を取得する
func (s *announcementService) List(ctx context.Context) ([]domain.Announcement, error) {
	return s.repo.List(ctx)
}

// Create 新規作成する
func (s *announcementService) Create(ctx context.Context, req *domain.AnnouncementRequest) (*domain.Announcement, error) {
	return s.repo.Create(ctx, req)
}

// Update 更新する
func (s *announcementService) Update(ctx context.Context, id string, req *domain.AnnouncementRequest) (*domain.Announcement, error) {
	return s.repo.Update(ctx, id, req)
}

// Delete 削除する
func (s *announcementService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
