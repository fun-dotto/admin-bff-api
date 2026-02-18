package domain

import "time"

// AnnouncementID お知らせIDの型
type AnnouncementID string

func (a AnnouncementID) String() string {
	return string(a)
}

// Announcement お知らせのドメインモデル
type Announcement struct {
	ID             AnnouncementID
	Title          string
	URL            string
	AvailableFrom  time.Time
	AvailableUntil *time.Time
}

// AnnouncementRequest お知らせのリクエストモデル
type AnnouncementRequest struct {
	Title          string
	URL            string
	AvailableFrom  time.Time
	AvailableUntil *time.Time
}
