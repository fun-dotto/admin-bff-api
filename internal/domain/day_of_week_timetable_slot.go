package domain

// DayOfWeekTimetableSlotID 曜日・時限IDの型
type DayOfWeekTimetableSlotID string

func (d DayOfWeekTimetableSlotID) String() string {
	return string(d)
}

// DayOfWeekTimetableSlot 曜日・時限のドメインモデル
type DayOfWeekTimetableSlot struct {
	ID            DayOfWeekTimetableSlotID
	DayOfWeek     DayOfWeek
	TimetableSlot TimetableSlot
}

// DayOfWeekTimetableSlotRequest 曜日・時限のリクエストモデル
type DayOfWeekTimetableSlotRequest struct {
	DayOfWeek     DayOfWeek
	TimetableSlot TimetableSlot
}
