package handler

import (
	"context"

	api "github.com/fun-dotto/api-template/generated"
)

type AnnouncementService interface {
	List(ctx context.Context) ([]api.Announcement, error)
	Create(ctx context.Context, req *api.AnnouncementRequest) (*api.Announcement, error)
	Update(ctx context.Context, id string, req *api.AnnouncementRequest) (*api.Announcement, error)
	Delete(ctx context.Context, id string) error
}

type Handler struct {
	announcementService AnnouncementService
}

// NewHandler 新規作成する
func NewHandler(announcementService AnnouncementService) *Handler {
	return &Handler{
		announcementService: announcementService,
	}
}

var _ api.ServerInterface = (*Handler)(nil)
