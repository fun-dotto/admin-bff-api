package domain

// SubjectTargetClass 対象学年・クラス
type SubjectTargetClass struct {
	Grade Grade
	Class *Class // 修士課程・博士課程対象の場合はnil
}

// SubjectRequirement 科目の必修・選択要件
type SubjectRequirement struct {
	Course          Course
	RequirementType SubjectRequirementType
}

// SubjectRequirementRequest 科目の必修・選択要件のリクエストモデル
type SubjectRequirementRequest struct {
	CourseID        string
	RequirementType SubjectRequirementType
}

// Subject 科目のドメインモデル
type Subject struct {
	ID                      string
	Name                    string
	Faculty                 Faculty
	Semester                CourseSemester
	DayOfWeekTimetableSlots []DayOfWeekTimetableSlot
	EligibleAttributes      []SubjectTargetClass
	Requirements            []SubjectRequirement
	Categories              []SubjectCategory
	SyllabusID              string
}

// SubjectRequest 科目のリクエストモデル
type SubjectRequest struct {
	Name                      string
	FacultyID                 string
	Semester                  CourseSemester
	DayOfWeekTimetableSlotIDs []string
	EligibleAttributes        []SubjectTargetClass
	Requirements              []SubjectRequirementRequest
	CategoryIDs               []string
	SyllabusID                string
}
