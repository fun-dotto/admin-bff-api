package domain

// DayOfWeekTimetableSlot 曜日・時限のドメインモデル
type DayOfWeekTimetableSlot struct {
	ID            string
	DayOfWeek     DayOfWeek
	TimetableSlot TimetableSlot
}

// DayOfWeekTimetableSlotRequest 曜日・時限のリクエストモデル
type DayOfWeekTimetableSlotRequest struct {
	DayOfWeek     DayOfWeek
	TimetableSlot TimetableSlot
}
