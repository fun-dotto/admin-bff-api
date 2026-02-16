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
