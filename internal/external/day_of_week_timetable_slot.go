package external

import (
	"github.com/fun-dotto/api-template/generated/external/subject_api"
	"github.com/fun-dotto/api-template/internal/domain"
)

// ToDomainDayOfWeekTimetableSlot 外部API形式からドメイン形式に変換する
func ToDomainDayOfWeekTimetableSlot(d subject_api.DayOfWeekTimetableSlot) domain.DayOfWeekTimetableSlot {
	return domain.DayOfWeekTimetableSlot{
		ID:            domain.DayOfWeekTimetableSlotID(d.Id),
		DayOfWeek:     domain.DayOfWeek(d.DayOfWeek),
		TimetableSlot: domain.TimetableSlot(d.TimetableSlot),
	}
}

// ToExternalDayOfWeekTimetableSlotRequest ドメイン形式から外部API形式に変換する
func ToExternalDayOfWeekTimetableSlotRequest(req *domain.DayOfWeekTimetableSlotRequest) subject_api.DayOfWeekTimetableSlotRequest {
	return subject_api.DayOfWeekTimetableSlotRequest{
		DayOfWeek:     subject_api.DottoFoundationV1DayOfWeek(req.DayOfWeek),
		TimetableSlot: subject_api.DottoFoundationV1TimetableSlot(req.TimetableSlot),
	}
}
