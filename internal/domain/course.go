package domain

// Course コースのドメインモデル
type Course struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CourseRequest コースのリクエストモデル
type CourseRequest struct {
	Name string `json:"name"`
}
