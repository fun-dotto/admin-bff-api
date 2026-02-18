package domain

import "time"

// Announcement お知らせのドメインモデル
type Announcement struct {
	ID             string
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
