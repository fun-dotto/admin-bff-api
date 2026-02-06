package repository

import (
	"context"
	"time"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/service"
)

type mockAnnouncementRepository struct{}

// NewMockAnnouncementRepository モックリポジトリを作成する
func NewMockAnnouncementRepository() service.AnnouncementRepository {
	return &mockAnnouncementRepository{}
}

// List 一覧を取得する（モック）
func (r *mockAnnouncementRepository) List(ctx context.Context) ([]domain.Announcement, error) {
	now := time.Now()
	until := now.Add(24 * time.Hour)
	return []domain.Announcement{
		{
			ID:             "1",
			Title:          "お知らせ1",
			URL:            "https://example.com/1",
			AvailableFrom:  now,
			AvailableUntil: &until,
		},
	}, nil
}

// Detail 詳細を取得する（モック）
func (r *mockAnnouncementRepository) Detail(ctx context.Context, id string) (*domain.Announcement, error) {
	now := time.Now()
	until := now.Add(24 * time.Hour)
	return &domain.Announcement{
		ID:             id,
		Title:          "お知らせ" + id,
		URL:            "https://example.com/" + id,
		AvailableFrom:  now,
		AvailableUntil: &until,
	}, nil
}

// Create 新規作成する（モック）
func (r *mockAnnouncementRepository) Create(ctx context.Context, req *domain.AnnouncementRequest) (*domain.Announcement, error) {
	return &domain.Announcement{
		ID:             "created-id",
		Title:          req.Title,
		URL:            req.URL,
		AvailableFrom:  req.AvailableFrom,
		AvailableUntil: req.AvailableUntil,
	}, nil
}

// Update 更新する（モック）
func (r *mockAnnouncementRepository) Update(ctx context.Context, id string, req *domain.AnnouncementRequest) (*domain.Announcement, error) {
	return &domain.Announcement{
		ID:             id,
		Title:          req.Title,
		URL:            req.URL,
		AvailableFrom:  req.AvailableFrom,
		AvailableUntil: req.AvailableUntil,
	}, nil
}

// Delete 削除する（モック）
func (r *mockAnnouncementRepository) Delete(ctx context.Context, id string) error {
	return nil
}
