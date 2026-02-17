package repository

import (
	"context"
	"fmt"

	"github.com/fun-dotto/api-template/generated/external/subject_api"
	"github.com/fun-dotto/api-template/internal/domain"
	"github.com/fun-dotto/api-template/internal/external"
	"github.com/fun-dotto/api-template/internal/service"
)

type dayOfWeekTimetableSlotRepository struct {
	client *subject_api.ClientWithResponses
}

// NewDayOfWeekTimetableSlotRepository 新規作成する
func NewDayOfWeekTimetableSlotRepository(client *subject_api.ClientWithResponses) service.DayOfWeekTimetableSlotRepository {
	return &dayOfWeekTimetableSlotRepository{client: client}
}

// List 一覧を取得する
func (r *dayOfWeekTimetableSlotRepository) List(ctx context.Context) ([]domain.DayOfWeekTimetableSlot, error) {
	response, err := r.client.DayOfWeekTimetableSlotsV1ListWithResponse(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get day of week timetable slots: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get day of week timetable slots: status %d", response.StatusCode())
	}

	result := make([]domain.DayOfWeekTimetableSlot, len(response.JSON200.DayOfWeekTimetableSlots))
	for i, d := range response.JSON200.DayOfWeekTimetableSlots {
		result[i] = external.ToDomainDayOfWeekTimetableSlot(d)
	}

	return result, nil
}

// Detail 詳細を取得する
func (r *dayOfWeekTimetableSlotRepository) Detail(ctx context.Context, id string) (*domain.DayOfWeekTimetableSlot, error) {
	response, err := r.client.DayOfWeekTimetableSlotsV1DetailWithResponse(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to get day of week timetable slot: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to get day of week timetable slot: status %d", response.StatusCode())
	}

	result := external.ToDomainDayOfWeekTimetableSlot(response.JSON200.DayOfWeekTimetableSlot)
	return &result, nil
}

// Create 新規作成する
func (r *dayOfWeekTimetableSlotRepository) Create(ctx context.Context, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error) {
	body := external.ToExternalDayOfWeekTimetableSlotRequest(req)

	response, err := r.client.DayOfWeekTimetableSlotsV1CreateWithResponse(ctx, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create day of week timetable slot: %w", err)
	}

	if response.JSON201 == nil {
		return nil, fmt.Errorf("failed to create day of week timetable slot: status %d", response.StatusCode())
	}

	result := external.ToDomainDayOfWeekTimetableSlot(response.JSON201.DayOfWeekTimetableSlot)
	return &result, nil
}

// Update 更新する
func (r *dayOfWeekTimetableSlotRepository) Update(ctx context.Context, id string, req *domain.DayOfWeekTimetableSlotRequest) (*domain.DayOfWeekTimetableSlot, error) {
	body := external.ToExternalDayOfWeekTimetableSlotRequest(req)

	response, err := r.client.DayOfWeekTimetableSlotsV1UpdateWithResponse(ctx, id, body)
	if err != nil {
		return nil, fmt.Errorf("failed to update day of week timetable slot: %w", err)
	}

	if response.JSON200 == nil {
		return nil, fmt.Errorf("failed to update day of week timetable slot: status %d", response.StatusCode())
	}

	result := external.ToDomainDayOfWeekTimetableSlot(response.JSON200.DayOfWeekTimetableSlot)
	return &result, nil
}

// Delete 削除する
func (r *dayOfWeekTimetableSlotRepository) Delete(ctx context.Context, id string) error {
	response, err := r.client.DayOfWeekTimetableSlotsV1DeleteWithResponse(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete day of week timetable slot: %w", err)
	}

	if response.StatusCode() != 204 {
		return fmt.Errorf("failed to delete day of week timetable slot: status %d", response.StatusCode())
	}

	return nil
}
