package domain

import "time"

// Announcement お知らせのドメインモデル
type Announcement struct {
	ID             string     `json:"id"`
	Title          string     `json:"title"`
	URL            string     `json:"url"`
	AvailableFrom  time.Time  `json:"availableFrom"`
	AvailableUntil *time.Time `json:"availableUntil,omitempty"`
}

// AnnouncementRequest お知らせのリクエストモデル
type AnnouncementRequest struct {
	Title          string     `json:"title"`
	URL            string     `json:"url"`
	AvailableFrom  time.Time  `json:"availableFrom"`
	AvailableUntil *time.Time `json:"availableUntil,omitempty"`
}
