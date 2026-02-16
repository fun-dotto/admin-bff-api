package domain

// DayOfWeekTimetableSlot 曜日・時限のドメインモデル
type DayOfWeekTimetableSlot struct {
	ID            string        `json:"id"`
	DayOfWeek     DayOfWeek     `json:"dayOfWeek"`
	TimetableSlot TimetableSlot `json:"timetableSlot"`
}

// DayOfWeekTimetableSlotRequest 曜日・時限のリクエストモデル
type DayOfWeekTimetableSlotRequest struct {
	DayOfWeek     DayOfWeek     `json:"dayOfWeek"`
	TimetableSlot TimetableSlot `json:"timetableSlot"`
}
