package repository

import (
	"context"

	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/service"
)

type MockDayOfWeekTimetableSlotRepository struct {
	ListFunc   func(ctx context.Context) ([]domain.DayOfWeekTimetableSlot, error)
	DetailFunc func(ctx context.Context, id string) (*domain.DayOfWeekTimetableSlot, error)
	CreateFunc func(ctx context.Context, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error)
	UpdateFunc func(ctx context.Context, id string, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error)
	DeleteFunc func(ctx context.Context, id string) error
}

func NewMockDayOfWeekTimetableSlotRepository() service.DayOfWeekTimetableSlotRepository {
	return &MockDayOfWeekTimetableSlotRepository{}
}

func (r *MockDayOfWeekTimetableSlotRepository) List(ctx context.Context) ([]domain.DayOfWeekTimetableSlot, error) {
	if r.ListFunc != nil {
		return r.ListFunc(ctx)
	}
	return []domain.DayOfWeekTimetableSlot{}, nil
}

func (r *MockDayOfWeekTimetableSlotRepository) Detail(ctx context.Context, id string) (*domain.DayOfWeekTimetableSlot, error) {
	if r.DetailFunc != nil {
		return r.DetailFunc(ctx, id)
	}
	return nil, nil
}

func (r *MockDayOfWeekTimetableSlotRepository) Create(ctx context.Context, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error) {
	if r.CreateFunc != nil {
		return r.CreateFunc(ctx, req)
	}
	return nil, nil
}

func (r *MockDayOfWeekTimetableSlotRepository) Update(ctx context.Context, id string, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error) {
	if r.UpdateFunc != nil {
		return r.UpdateFunc(ctx, id, req)
	}
	return nil, nil
}

func (r *MockDayOfWeekTimetableSlotRepository) Delete(ctx context.Context, id string) error {
	if r.DeleteFunc != nil {
		return r.DeleteFunc(ctx, id)
	}
	return nil
}
