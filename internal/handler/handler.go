package handler

import (
	"context"

	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/generated/external/announcement_api"
	"github.com/fun-dotto/api-template/internal/domain"
)


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

type SubjectService interface {
	List(ctx context.Context) ([]domain.Subject, error)
	Detail(ctx context.Context, id string) (*domain.Subject, error)
	Create(ctx context.Context, req *domain.SubjectRequest) (*domain.Subject, error)
	Update(ctx context.Context, id string, req *domain.SubjectRequest) (*domain.Subject, error)
	Delete(ctx context.Context, id string) error
}

type Handler struct {
	announcementClient            *announcement_api.ClientWithResponses
	facultyService                FacultyService
	courseService                 CourseService
	dayOfWeekTimetableSlotService DayOfWeekTimetableSlotService
	subjectCategoryService        SubjectCategoryService
	subjectService                SubjectService
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) WithAnnouncementClient(c *announcement_api.ClientWithResponses) *Handler {
	h.announcementClient = c
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

func (h *Handler) WithSubjectService(s SubjectService) *Handler {
	h.subjectService = s
	return h
}

var _ api.ServerInterface = (*Handler)(nil)
