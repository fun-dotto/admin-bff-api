package handler

import (
	"context"

	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/domain"
)

type AnnouncementService interface {
	List(ctx context.Context) ([]domain.Announcement, error)
	Detail(ctx context.Context, id string) (*domain.Announcement, error)
	Create(ctx context.Context, req *domain.AnnouncementRequest) (*domain.Announcement, error)
	Update(ctx context.Context, id string, req *domain.AnnouncementRequest) (*domain.Announcement, error)
	Delete(ctx context.Context, id string) error
}

type Handler struct {
	announcementService AnnouncementService
}

func NewHandler(announcementService AnnouncementService) *Handler {
	return &Handler{
		announcementService: announcementService,
	}
}

var _ api.ServerInterface = (*Handler)(nil)
