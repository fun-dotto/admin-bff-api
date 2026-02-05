package external

import (
	"github.com/fun-dotto/api-template/generated/external/announcement_api"
	"github.com/fun-dotto/api-template/internal/domain"
)

// ToDomainAnnouncement 外部API形式からドメイン形式に変換する
func ToDomainAnnouncement(a announcement_api.Announcement) domain.Announcement {
	return domain.Announcement{
		ID:             a.Id,
		Title:          a.Title,
		URL:            a.Url,
		AvailableFrom:  a.AvailableFrom,
		AvailableUntil: a.AvailableUntil,
	}
}

// ToExternalAnnouncementRequest ドメイン形式から外部API形式に変換する
func ToExternalAnnouncementRequest(req *domain.AnnouncementRequest) announcement_api.AnnouncementRequest {
	return announcement_api.AnnouncementRequest{
		Title:          req.Title,
		Url:            req.URL,
		AvailableFrom:  req.AvailableFrom,
		AvailableUntil: req.AvailableUntil,
	}
}
