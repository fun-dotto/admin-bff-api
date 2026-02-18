package external

import (
	"github.com/fun-dotto/api-template/generated/external/subject_api"
	"github.com/fun-dotto/api-template/internal/domain"
)

// ToDomainSubject 外部API形式からドメイン形式に変換する
func ToDomainSubject(s subject_api.Subject) domain.Subject {
	// DayOfWeekTimetableSlots の変換
	dayOfWeekTimetableSlots := make([]domain.DayOfWeekTimetableSlot, len(s.DayOfWeekTimetableSlots))
	for i, d := range s.DayOfWeekTimetableSlots {
		dayOfWeekTimetableSlots[i] = ToDomainDayOfWeekTimetableSlot(d)
	}

	// EligibleAttributes の変換
	eligibleAttributes := make([]domain.SubjectTargetClass, len(s.EligibleAttributes))
	for i, e := range s.EligibleAttributes {
		eligibleAttributes[i] = ToDomainSubjectTargetClass(e)
	}

	// Requirements の変換
	requirements := make([]domain.SubjectRequirement, len(s.Requirements))
	for i, r := range s.Requirements {
		requirements[i] = ToDomainSubjectRequirement(r)
	}

	// Categories の変換
	categories := make([]domain.SubjectCategory, len(s.Categories))
	for i, c := range s.Categories {
		categories[i] = ToDomainSubjectCategory(c)
	}

	return domain.Subject{
		ID:                      domain.SubjectID(s.Id),
		Name:                    s.Name,
		Faculty:                 ToDomainFaculty(s.Faculty),
		Semester:                domain.CourseSemester(s.Semester),
		DayOfWeekTimetableSlots: dayOfWeekTimetableSlots,
		EligibleAttributes:      eligibleAttributes,
		Requirements:            requirements,
		Categories:              categories,
		SyllabusID:              s.SyllabusId,
	}
}

// ToDomainSubjectTargetClass 外部API形式からドメイン形式に変換する
func ToDomainSubjectTargetClass(e subject_api.SubjectTargetClass) domain.SubjectTargetClass {
	var class *domain.Class
	if e.Class != nil {
		c := domain.Class(*e.Class)
		class = &c
	}
	return domain.SubjectTargetClass{
		Grade: domain.Grade(e.Grade),
		Class: class,
	}
}

// ToDomainSubjectRequirement 外部API形式からドメイン形式に変換する
func ToDomainSubjectRequirement(r subject_api.SubjectRequirement) domain.SubjectRequirement {
	return domain.SubjectRequirement{
		Course:          ToDomainCourse(r.Course),
		RequirementType: domain.SubjectRequirementType(r.RequirementType),
	}
}

// ToExternalSubjectRequest ドメイン形式から外部API形式に変換する
func ToExternalSubjectRequest(req *domain.SubjectRequest) subject_api.SubjectRequest {
	// EligibleAttributes の変換
	eligibleAttributes := make([]subject_api.SubjectTargetClass, len(req.EligibleAttributes))
	for i, e := range req.EligibleAttributes {
		eligibleAttributes[i] = ToExternalSubjectTargetClass(e)
	}

	// Requirements の変換
	requirements := make([]subject_api.SubjectRequirementRequest, len(req.Requirements))
	for i, r := range req.Requirements {
		requirements[i] = ToExternalSubjectRequirementRequest(r)
	}

	return subject_api.SubjectRequest{
		Name:                      req.Name,
		FacultyId:                 req.FacultyID,
		Semester:                  subject_api.DottoFoundationV1CourseSemester(req.Semester),
		DayOfWeekTimetableSlotIds: req.DayOfWeekTimetableSlotIDs,
		EligibleAttributes:        eligibleAttributes,
		Requirements:              requirements,
		CategoryIds:               req.CategoryIDs,
		SyllabusId:                req.SyllabusID,
	}
}

// ToExternalSubjectTargetClass ドメイン形式から外部API形式に変換する
func ToExternalSubjectTargetClass(e domain.SubjectTargetClass) subject_api.SubjectTargetClass {
	var class *subject_api.DottoFoundationV1Class
	if e.Class != nil {
		c := subject_api.DottoFoundationV1Class(*e.Class)
		class = &c
	}
	return subject_api.SubjectTargetClass{
		Grade: subject_api.DottoFoundationV1Grade(e.Grade),
		Class: class,
	}
}

// ToExternalSubjectRequirementRequest ドメイン形式から外部API形式に変換する
func ToExternalSubjectRequirementRequest(r domain.SubjectRequirementRequest) subject_api.SubjectRequirementRequest {
	return subject_api.SubjectRequirementRequest{
		CourseId:        r.CourseID,
		RequirementType: subject_api.DottoFoundationV1SubjectRequirementType(r.RequirementType),
	}
}
