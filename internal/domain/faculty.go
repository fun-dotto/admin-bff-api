package domain

// Faculty 教員のドメインモデル
type Faculty struct {
	ID    string
	Name  string
	Email string
}

// FacultyRequest 教員のリクエストモデル
type FacultyRequest struct {
	Name  string
	Email string
}
