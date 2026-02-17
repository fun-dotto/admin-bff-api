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

type CourseService interface {
	List(ctx context.Context) ([]domain.Course, error)
	Detail(ctx context.Context, id string) (*domain.Course, error)
	Create(ctx context.Context, req *domain.CourseRequest) (*domain.Course, error)
	Update(ctx context.Context, id string, req *domain.CourseRequest) (*domain.Course, error)
	Delete(ctx context.Context, id string) error
}

type DayOfWeekTimetableSlotService interface {
	List(ctx context.Context) ([]domain.DayOfWeekTimetableSlot, error)
	Detail(ctx context.Context, id string) (*domain.DayOfWeekTimetableSlot, error)
	Create(ctx context.Context, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error)
	Update(ctx context.Context, id string, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error)
	Delete(ctx context.Context, id string) error
}

type SubjectCategoryService interface {
	List(ctx context.Context) ([]domain.SubjectCategory, error)
	Detail(ctx context.Context, id string) (*domain.SubjectCategory, error)
	Create(ctx context.Context, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error)
	Update(ctx context.Context, id string, req *domain.SubjectCategoryRequest) (*domain.SubjectCategory, error)
	Delete(ctx context.Context, id string) error
}

type Handler struct {
	announcementService           AnnouncementService
	facultyService                FacultyService
	courseService                 CourseService
	dayOfWeekTimetableSlotService DayOfWeekTimetableSlotService
	subjectCategoryService        SubjectCategoryService
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

func (h *Handler) WithCourseService(s CourseService) *Handler {
	h.courseService = s
	return h
}

func (h *Handler) WithDayOfWeekTimetableSlotService(s DayOfWeekTimetableSlotService) *Handler {
	h.dayOfWeekTimetableSlotService = s
	return h
}

func (h *Handler) WithSubjectCategoryService(s SubjectCategoryService) *Handler {
	h.subjectCategoryService = s
	return h
}

var _ api.ServerInterface = (*Handler)(nil)
