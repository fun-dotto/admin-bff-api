package domain

// SubjectCategory 科目群・科目区分のドメインモデル
type SubjectCategory struct {
	ID   string
	Name string
}

// SubjectCategoryRequest 科目群・科目区分のリクエストモデル
type SubjectCategoryRequest struct {
	Name string
}
