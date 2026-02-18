package domain

// SubjectCategoryID 科目群・科目区分IDの型
type SubjectCategoryID string

func (s SubjectCategoryID) String() string {
	return string(s)
}

// SubjectCategory 科目群・科目区分のドメインモデル
type SubjectCategory struct {
	ID   SubjectCategoryID
	Name string
}

// SubjectCategoryRequest 科目群・科目区分のリクエストモデル
type SubjectCategoryRequest struct {
	Name string
}
