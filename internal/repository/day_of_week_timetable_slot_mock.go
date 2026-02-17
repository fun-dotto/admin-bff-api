package repository

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/service"
)

type mockDayOfWeekTimetableSlotRepository struct{}

// NewMockDayOfWeekTimetableSlotRepository モックリポジトリを作成する
func NewMockDayOfWeekTimetableSlotRepository() service.DayOfWeekTimetableSlotRepository {
	return &mockDayOfWeekTimetableSlotRepository{}
}

// List 一覧を取得する（モック）
func (r *mockDayOfWeekTimetableSlotRepository) List(ctx context.Context) ([]domain.DayOfWeekTimetableSlot, error) {
	return []domain.DayOfWeekTimetableSlot{
		{
			ID:            "1",
			DayOfWeek:     domain.DayOfWeekMonday,
			TimetableSlot: domain.TimetableSlotSlot1,
		},
	}, nil
}

// Detail 詳細を取得する（モック）
func (r *mockDayOfWeekTimetableSlotRepository) Detail(ctx context.Context, id string) (*domain.DayOfWeekTimetableSlot, error) {
	return &domain.DayOfWeekTimetableSlot{
		ID:            id,
		DayOfWeek:     domain.DayOfWeekMonday,
		TimetableSlot: domain.TimetableSlotSlot1,
	}, nil
}

// Create 新規作成する（モック）
func (r *mockDayOfWeekTimetableSlotRepository) Create(ctx context.Context, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error) {
	return &domain.DayOfWeekTimetableSlot{
		ID:            "created-id",
		DayOfWeek:     req.DayOfWeek,
		TimetableSlot: req.TimetableSlot,
	}, nil
}

// Update 更新する（モック）
func (r *mockDayOfWeekTimetableSlotRepository) Update(ctx context.Context, id string, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error) {
	return &domain.DayOfWeekTimetableSlot{
		ID:            id,
		DayOfWeek:     req.DayOfWeek,
		TimetableSlot: req.TimetableSlot,
	}, nil
}

// Delete 削除する（モック）
func (r *mockDayOfWeekTimetableSlotRepository) Delete(ctx context.Context, id string) error {
	return nil
}
