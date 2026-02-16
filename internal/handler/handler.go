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

type FacultyService interface {
	List(ctx context.Context) ([]domain.Faculty, error)
	Detail(ctx context.Context, id string) (*domain.Faculty, error)
	Create(ctx context.Context, req *domain.FacultyRequest) (*domain.Faculty, error)
	Update(ctx context.Context, id string, req *domain.FacultyRequest) (*domain.Faculty, error)
	Delete(ctx context.Context, id string) error
}

type Handler struct {
	announcementService AnnouncementService
	facultyService      FacultyService
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) WithAnnouncementService(s AnnouncementService) *Handler {
	h.announcementService = s
	return h
}

func (h *Handler) WithFacultyService(s FacultyService) *Handler {
	h.facultyService = s
	return h
}

var _ api.ServerInterface = (*Handler)(nil)
