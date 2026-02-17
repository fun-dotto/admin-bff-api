package handler

import (
	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/domain"
)

// ToAPIAnnouncement ドメインモデルをAPIモデルに変換する
func ToAPIAnnouncement(a *domain.Announcement) api.AnnouncementServiceAnnouncement {
	return api.AnnouncementServiceAnnouncement{
		Id:             a.ID,
		Title:          a.Title,
		Url:            a.URL,
		AvailableFrom:  a.AvailableFrom,
		AvailableUntil: a.AvailableUntil,
	}
}

// ToAPIAnnouncements ドメインモデルの配列をAPIモデルの配列に変換する
func ToAPIAnnouncements(announcements []domain.Announcement) []api.AnnouncementServiceAnnouncement {
	result := make([]api.AnnouncementServiceAnnouncement, len(announcements))
	for i, a := range announcements {
		result[i] = ToAPIAnnouncement(&a)
	}
	return result
}

// ToDomainAnnouncementRequest APIモデルをドメインモデルに変換する
func ToDomainAnnouncementRequest(req *api.AnnouncementServiceAnnouncementRequest) *domain.AnnouncementRequest {
	return &domain.AnnouncementRequest{
		Title:          req.Title,
		URL:            req.Url,
		AvailableFrom:  req.AvailableFrom,
		AvailableUntil: req.AvailableUntil,
	}
}

// ToAPIFaculty ドメインモデルをAPIモデルに変換する
func ToAPIFaculty(f *domain.Faculty) api.SubjectServiceFaculty {
	return api.SubjectServiceFaculty{
		Id:    f.ID,
		Name:  f.Name,
		Email: f.Email,
	}
}

// ToAPIFaculties ドメインモデルの配列をAPIモデルの配列に変換する
func ToAPIFaculties(faculties []domain.Faculty) []api.SubjectServiceFaculty {
	result := make([]api.SubjectServiceFaculty, len(faculties))
	for i, f := range faculties {
		result[i] = ToAPIFaculty(&f)
	}
	return result
}

// ToDomainFacultyRequest APIモデルをドメインモデルに変換する
func ToDomainFacultyRequest(req *api.SubjectServiceFacultyRequest) *domain.FacultyRequest {
	return &domain.FacultyRequest{
		Name:  req.Name,
		Email: req.Email,
	}
}

// ToAPICourse ドメインモデルをAPIモデルに変換する
func ToAPICourse(c *domain.Course) api.SubjectServiceCourse {
	return api.SubjectServiceCourse{
		Id:   c.ID,
		Name: c.Name,
	}
}

// ToAPICourses ドメインモデルの配列をAPIモデルの配列に変換する
func ToAPICourses(courses []domain.Course) []api.SubjectServiceCourse {
	result := make([]api.SubjectServiceCourse, len(courses))
	for i, c := range courses {
		result[i] = ToAPICourse(&c)
	}
	return result
}

// ToDomainCourseRequest APIモデルをドメインモデルに変換する
func ToDomainCourseRequest(req *api.SubjectServiceCourseRequest) *domain.CourseRequest {
	return &domain.CourseRequest{
		Name: req.Name,
	}
}

// ToAPIDayOfWeekTimetableSlot ドメインモデルをAPIモデルに変換する
func ToAPIDayOfWeekTimetableSlot(d *domain.DayOfWeekTimetableSlot) api.SubjectServiceDayOfWeekTimetableSlot {
	return api.SubjectServiceDayOfWeekTimetableSlot{
		Id:            d.ID,
		DayOfWeek:     api.DottoFoundationV1DayOfWeek(d.DayOfWeek),
		TimetableSlot: api.DottoFoundationV1TimetableSlot(d.TimetableSlot),
	}
}

// ToAPIDayOfWeekTimetableSlots ドメインモデルの配列をAPIモデルの配列に変換する
func ToAPIDayOfWeekTimetableSlots(slots []domain.DayOfWeekTimetableSlot) []api.SubjectServiceDayOfWeekTimetableSlot {
	result := make([]api.SubjectServiceDayOfWeekTimetableSlot, len(slots))
	for i, d := range slots {
		result[i] = ToAPIDayOfWeekTimetableSlot(&d)
	}
	return result
}

// ToDomainDayOfWeekTimetableSlotRequest APIモデルをドメインモデルに変換する
func ToDomainDayOfWeekTimetableSlotRequest(req *api.SubjectServiceDayOfWeekTimetableSlotRequest) *domain.DayOfWeekTimetableSlotRequest {
	return &domain.DayOfWeekTimetableSlotRequest{
		DayOfWeek:     domain.DayOfWeek(req.DayOfWeek),
		TimetableSlot: domain.TimetableSlot(req.TimetableSlot),
	}
}

// ToAPISubjectCategory ドメインモデルをAPIモデルに変換する
func ToAPISubjectCategory(c *domain.SubjectCategory) api.SubjectServiceSubjectCategory {
	return api.SubjectServiceSubjectCategory{
		Id:   c.ID,
		Name: c.Name,
	}
}

// ToAPISubjectCategories ドメインモデルの配列をAPIモデルの配列に変換する
func ToAPISubjectCategories(categories []domain.SubjectCategory) []api.SubjectServiceSubjectCategory {
	result := make([]api.SubjectServiceSubjectCategory, len(categories))
	for i, c := range categories {
		result[i] = ToAPISubjectCategory(&c)
	}
	return result
}

// ToDomainSubjectCategoryRequest APIモデルをドメインモデルに変換する
func ToDomainSubjectCategoryRequest(req *api.SubjectServiceSubjectCategoryRequest) *domain.SubjectCategoryRequest {
	return &domain.SubjectCategoryRequest{
		Name: req.Name,
	}
}

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
		Id:                      s.ID,
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
