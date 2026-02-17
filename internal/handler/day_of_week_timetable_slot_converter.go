package handler

import (
	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/domain"
)

// ToAPIDayOfWeekTimetableSlot ドメインモデルをAPIモデルに変換する
func ToAPIDayOfWeekTimetableSlot(d *domain.DayOfWeekTimetableSlot) api.SubjectServiceDayOfWeekTimetableSlot {
	return api.SubjectServiceDayOfWeekTimetableSlot{
		Id:            d.ID,
		DayOfWeek:     api.DottoFoundationV1DayOfWeek(d.DayOfWeek),
		TimetableSlot: api.DottoFoundationV1TimetableSlot(d.TimetableSlot),
	}
}

// ToAPIDayOfWeekTimetableSlots ドメインモデルの配列をAPIモデルの配列に変換する
func ToAPIDayOfWeekTimetableSlots(slots []domain.DayOfWeekTimetableSlot) []api.SubjectServiceDayOfWeekTimetableSlot {
	result := make([]api.SubjectServiceDayOfWeekTimetableSlot, len(slots))
	for i, d := range slots {
		result[i] = ToAPIDayOfWeekTimetableSlot(&d)
	}
	return result
}

// ToDomainDayOfWeekTimetableSlotRequest APIモデルをドメインモデルに変換する
func ToDomainDayOfWeekTimetableSlotRequest(req *api.SubjectServiceDayOfWeekTimetableSlotRequest) *domain.DayOfWeekTimetableSlotRequest {
	return &domain.DayOfWeekTimetableSlotRequest{
		DayOfWeek:     domain.DayOfWeek(req.DayOfWeek),
		TimetableSlot: domain.TimetableSlot(req.TimetableSlot),
	}
}
