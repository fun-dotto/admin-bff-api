package handler

import (
	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/domain"
)

// ToAPISubject ドメインモデルをAPIモデルに変換する
func ToAPISubject(s *domain.Subject) api.SubjectServiceSubject {
	// DayOfWeekTimetableSlots
	slots := make([]api.SubjectServiceDayOfWeekTimetableSlot, len(s.DayOfWeekTimetableSlots))
	for i, slot := range s.DayOfWeekTimetableSlots {
		slots[i] = ToAPIDayOfWeekTimetableSlot(&slot)
	}

	// EligibleAttributes
	eligibleAttrs := make([]api.SubjectServiceSubjectTargetClass, len(s.EligibleAttributes))
	for i, attr := range s.EligibleAttributes {
		eligibleAttrs[i] = ToAPISubjectTargetClass(&attr)
	}

	// Requirements
	requirements := make([]api.SubjectServiceSubjectRequirement, len(s.Requirements))
	for i, req := range s.Requirements {
		requirements[i] = ToAPISubjectRequirement(&req)
	}

	// Categories
	categories := make([]api.SubjectServiceSubjectCategory, len(s.Categories))
	for i, cat := range s.Categories {
		categories[i] = ToAPISubjectCategory(&cat)
	}

	return api.SubjectServiceSubject{
		Id:                      s.ID.String(),
		Name:                    s.Name,
		Faculty:                 ToAPIFaculty(&s.Faculty),
		Semester:                api.DottoFoundationV1CourseSemester(s.Semester),
		DayOfWeekTimetableSlots: slots,
		EligibleAttributes:      eligibleAttrs,
		Requirements:            requirements,
		Categories:              categories,
		SyllabusId:              s.SyllabusID,
	}
}

// ToAPISubjects ドメインモデルの配列をAPIモデルの配列に変換する
func ToAPISubjects(subjects []domain.Subject) []api.SubjectServiceSubject {
	result := make([]api.SubjectServiceSubject, len(subjects))
	for i, s := range subjects {
		result[i] = ToAPISubject(&s)
	}
	return result
}

// ToAPISubjectTargetClass ドメインモデルをAPIモデルに変換する
func ToAPISubjectTargetClass(s *domain.SubjectTargetClass) api.SubjectServiceSubjectTargetClass {
	var class *api.DottoFoundationV1Class
	if s.Class != nil {
		c := api.DottoFoundationV1Class(*s.Class)
		class = &c
	}
	return api.SubjectServiceSubjectTargetClass{
		Grade: api.DottoFoundationV1Grade(s.Grade),
		Class: class,
	}
}

// ToAPISubjectRequirement ドメインモデルをAPIモデルに変換する
func ToAPISubjectRequirement(r *domain.SubjectRequirement) api.SubjectServiceSubjectRequirement {
	return api.SubjectServiceSubjectRequirement{
		Course:          ToAPICourse(&r.Course),
		RequirementType: api.DottoFoundationV1SubjectRequirementType(r.RequirementType),
	}
}

// ToDomainSubjectRequest APIモデルをドメインモデルに変換する
func ToDomainSubjectRequest(req *api.SubjectServiceSubjectRequest) *domain.SubjectRequest {
	// EligibleAttributes
	eligibleAttrs := make([]domain.SubjectTargetClass, len(req.EligibleAttributes))
	for i, attr := range req.EligibleAttributes {
		eligibleAttrs[i] = ToDomainSubjectTargetClass(&attr)
	}

	// Requirements
	requirements := make([]domain.SubjectRequirementRequest, len(req.Requirements))
	for i, r := range req.Requirements {
		requirements[i] = ToDomainSubjectRequirementRequest(&r)
	}

	return &domain.SubjectRequest{
		Name:                      req.Name,
		FacultyID:                 req.FacultyId,
		Semester:                  domain.CourseSemester(req.Semester),
		DayOfWeekTimetableSlotIDs: req.DayOfWeekTimetableSlotIds,
		EligibleAttributes:        eligibleAttrs,
		Requirements:              requirements,
		CategoryIDs:               req.CategoryIds,
		SyllabusID:                req.SyllabusId,
	}
}

// ToDomainSubjectTargetClass APIモデルをドメインモデルに変換する
func ToDomainSubjectTargetClass(s *api.SubjectServiceSubjectTargetClass) domain.SubjectTargetClass {
	var class *domain.Class
	if s.Class != nil {
		c := domain.Class(*s.Class)
		class = &c
	}
	return domain.SubjectTargetClass{
		Grade: domain.Grade(s.Grade),
		Class: class,
	}
}

// ToDomainSubjectRequirementRequest APIモデルをドメインモデルに変換する
func ToDomainSubjectRequirementRequest(r *api.SubjectServiceSubjectRequirementRequest) domain.SubjectRequirementRequest {
	return domain.SubjectRequirementRequest{
		CourseID:        r.CourseId,
		RequirementType: domain.SubjectRequirementType(r.RequirementType),
	}
}
