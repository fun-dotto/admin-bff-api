package domain

// CourseID コースIDの型
type CourseID string

func (c CourseID) String() string {
	return string(c)
}

// Course コースのドメインモデル
type Course struct {
	ID   CourseID
	Name string
}

// CourseRequest コースのリクエストモデル
type CourseRequest struct {
	Name string
}
