package repository

import (
	"context"
	"fmt"
	"net/http"

	"github.com/fun-dotto/api-template/generated/external/announcement_api"
	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/external"
	"github.com/fun-dotto/api-template/internal/middleware"
	"github.com/fun-dotto/api-template/internal/service"
)

type announcementRepository struct {
	client *announcement_api.ClientWithResponses
}

func NewAnnouncementRepository(client *announcement_api.ClientWithResponses) service.AnnouncementRepository {
	return &announcementRepository{client: client}
}

// withAuth はコンテキストからIDトークンを取得し、Authorizationヘッダーを追加するRequestEditorFnを返す
func withAuth(ctx context.Context) announcement_api.RequestEditorFn {
	return func(_ context.Context, req *http.Request) error {
		if token, ok := middleware.GetRawIDTokenFromContext(ctx); ok {
			req.Header.Set("Authorization", "Bearer "+token)
		}
		return nil
	}
}

// List 一覧を取得する
func (r *announcementRepository) List(ctx context.Context) ([]domain.Announcement, error) {
	response, err := r.client.AnnouncementsV1ListWithResponse(ctx, nil, withAuth(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to get announcements: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get announcements: status %d", response.StatusCode())
	}

	// 外部API形式 → ドメイン形式に変換
	result := make([]domain.Announcement, len(response.JSON200.Announcements))
	for i, a := range response.JSON200.Announcements {
		result[i] = external.ToDomainAnnouncement(a)
	}

	return result, nil
}

// Create 新規作成する
func (r *announcementRepository) Create(ctx context.Context, req *domain.AnnouncementRequest) (*domain.Announcement, error) {
	body := external.ToExternalAnnouncementRequest(req)

	response, err := r.client.AnnouncementsV1CreateWithResponse(ctx, body, withAuth(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to create announcement: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to create announcement: status %d", response.StatusCode())
	}

	result := external.ToDomainAnnouncement(response.JSON200.Announcement)
	return &result, nil
}

// Update 更新する
func (r *announcementRepository) Update(ctx context.Context, id string, req *domain.AnnouncementRequest) (*domain.Announcement, error) {
	body := external.ToExternalAnnouncementRequest(req)

	response, err := r.client.AnnouncementsV1UpdateWithResponse(ctx, id, body, withAuth(ctx))
	if err != nil {
		return nil, fmt.Errorf("failed to update announcement: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to update announcement: status %d", response.StatusCode())
	}

	result := external.ToDomainAnnouncement(response.JSON200.Announcement)
	return &result, nil
}

// Delete 削除する
func (r *announcementRepository) Delete(ctx context.Context, id string) error {
	response, err := r.client.AnnouncementsV1DeleteWithResponse(ctx, id, withAuth(ctx))
	if err != nil {
		return fmt.Errorf("failed to delete announcement: %w", err)
	}

	if response.StatusCode() != 204 {
		return fmt.Errorf("failed to delete announcement: status %d", response.StatusCode())
	}

	return nil
}
