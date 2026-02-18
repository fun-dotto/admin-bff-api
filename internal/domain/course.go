package domain

// Course コースのドメインモデル
type Course struct {
	ID   string
	Name string
}

// CourseRequest コースのリクエストモデル
type CourseRequest struct {
	Name string
}
