package service

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/handler"
)

// DayOfWeekTimetableSlotRepository Serviceが必要とするRepositoryのインターフェース
type DayOfWeekTimetableSlotRepository interface {
	List(ctx context.Context) ([]domain.DayOfWeekTimetableSlot, error)
	Detail(ctx context.Context, id string) (*domain.DayOfWeekTimetableSlot, error)
	Create(ctx context.Context, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error)
	Update(ctx context.Context, id string, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error)
	Delete(ctx context.Context, id string) error
}

type dayOfWeekTimetableSlotService struct {
	repo DayOfWeekTimetableSlotRepository
}

// NewDayOfWeekTimetableSlotService 新規作成する
func NewDayOfWeekTimetableSlotService(repo DayOfWeekTimetableSlotRepository) handler.DayOfWeekTimetableSlotService {
	return &dayOfWeekTimetableSlotService{
		repo: repo,
	}
}

// List 一覧を取得する
func (s *dayOfWeekTimetableSlotService) List(ctx context.Context) ([]domain.DayOfWeekTimetableSlot, error) {
	return s.repo.List(ctx)
}

// Detail 詳細を取得する
func (s *dayOfWeekTimetableSlotService) Detail(ctx context.Context, id string) (*domain.DayOfWeekTimetableSlot, error) {
	return s.repo.Detail(ctx, id)
}

// Create 新規作成する
func (s *dayOfWeekTimetableSlotService) Create(ctx context.Context, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error) {
	return s.repo.Create(ctx, req)
}

// Update 更新する
func (s *dayOfWeekTimetableSlotService) Update(ctx context.Context, id string, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error) {
	return s.repo.Update(ctx, id, req)
}

// Delete 削除する
func (s *dayOfWeekTimetableSlotService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
