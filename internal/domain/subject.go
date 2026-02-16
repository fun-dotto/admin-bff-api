package domain

// SubjectTargetClass 対象学年・クラス
type SubjectTargetClass struct {
	Grade Grade  `json:"grade"`
	Class *Class `json:"class,omitempty"` // 修士課程・博士課程対象の場合はnil
}

// SubjectRequirement 科目の必修・選択要件
type SubjectRequirement struct {
	Course          Course                 `json:"course"`
	RequirementType SubjectRequirementType `json:"requirementType"`
}

// SubjectRequirementRequest 科目の必修・選択要件のリクエストモデル
type SubjectRequirementRequest struct {
	CourseID        string                 `json:"courseId"`
	RequirementType SubjectRequirementType `json:"requirementType"`
}

// Subject 科目のドメインモデル
type Subject struct {
	ID                      string                   `json:"id"`
	Name                    string                   `json:"name"`
	Faculty                 Faculty                  `json:"faculty"`
	Semester                CourseSemester           `json:"semester"`
	DayOfWeekTimetableSlots []DayOfWeekTimetableSlot `json:"dayOfWeekTimetableSlots"`
	EligibleAttributes      []SubjectTargetClass     `json:"eligibleAttributes"`
	Requirements            []SubjectRequirement     `json:"requirements"`
	Categories              []SubjectCategory        `json:"categories"`
	SyllabusID              string                   `json:"syllabusId"`
}

// SubjectRequest 科目のリクエストモデル
type SubjectRequest struct {
	Name                      string                      `json:"name"`
	FacultyID                 string                      `json:"facultyId"`
	Semester                  CourseSemester              `json:"semester"`
	DayOfWeekTimetableSlotIDs []string                    `json:"dayOfWeekTimetableSlotIds"`
	EligibleAttributes        []SubjectTargetClass        `json:"eligibleAttributes"`
	Requirements              []SubjectRequirementRequest `json:"requirements"`
	CategoryIDs               []string                    `json:"categoryIds"`
	SyllabusID                string                      `json:"syllabusId"`
}
