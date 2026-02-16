package domain

// Faculty 教員のドメインモデル
type Faculty struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// FacultyRequest 教員のリクエストモデル
type FacultyRequest struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}
