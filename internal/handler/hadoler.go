package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	api "github.com/fun-dotto/api-template/generated"
)

// Handler implements api.ServerInterface
type Handler struct {
	// TODO: Add dependencies like DB, services, etc.
}

// NewHandler creates a new Handler instance
func NewHandler() *Handler {
	return &Handler{}
}

// Ensure Handler implements api.ServerInterface
var _ api.ServerInterface = (*Handler)(nil)

// AnnouncementsV1List returns a list of announcements
// (GET /v1/announcements)
func (h *Handler) AnnouncementsV1List(c *gin.Context) {
	// TODO: Implement actual database query
	announcements := []api.Announcement{
		{
			Id:            "1",
			Title:         "Sample Announcement",
			Url:           "https://example.com/announcement/1",
			AvailableFrom: time.Now(),
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"announcements": announcements,
	})
}

// AnnouncementsV1Create creates a new announcement
// (POST /v1/announcements)
func (h *Handler) AnnouncementsV1Create(c *gin.Context) {
	var req api.AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement actual database insert
	announcement := api.Announcement{
		Id:             uuid.New().String(),
		Title:          req.Title,
		Url:            req.Url,
		AvailableFrom:  req.AvailableFrom,
		AvailableUntil: req.AvailableUntil,
	}

	c.JSON(http.StatusOK, gin.H{
		"announcement": announcement,
	})
}

// AnnouncementsV1Delete deletes an announcement by ID
// (DELETE /v1/announcements/{id})
func (h *Handler) AnnouncementsV1Delete(c *gin.Context, id string) {
	// TODO: Implement actual database delete
	// For now, just return 204 No Content

	c.Status(http.StatusNoContent)
}

// AnnouncementsV1Update updates an announcement by ID
// (PUT /v1/announcements/{id})
func (h *Handler) AnnouncementsV1Update(c *gin.Context, id string) {
	var req api.AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// TODO: Implement actual database update
	announcement := api.Announcement{
		Id:             id,
		Title:          req.Title,
		Url:            req.Url,
		AvailableFrom:  req.AvailableFrom,
		AvailableUntil: req.AvailableUntil,
	}

	c.JSON(http.StatusOK, gin.H{
		"announcement": announcement,
	})
}
