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
