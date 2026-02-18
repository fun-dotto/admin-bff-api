package domain

// FacultyID 教員IDの型
type FacultyID string

func (f FacultyID) String() string {
	return string(f)
}

// Faculty 教員のドメインモデル
type Faculty struct {
	ID    FacultyID
	Name  string
	Email string
}

// FacultyRequest 教員のリクエストモデル
type FacultyRequest struct {
	Name  string
	Email string
}
