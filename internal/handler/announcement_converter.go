package handler

import (
	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/domain"
)

// ToAPIAnnouncement ドメインモデルをAPIモデルに変換する
func ToAPIAnnouncement(a *domain.Announcement) api.AnnouncementServiceAnnouncement {
	return api.AnnouncementServiceAnnouncement{
		Id:             a.ID.String(),
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
