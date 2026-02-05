package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/api-template/generated"
)

// AnnouncementsV1List 一覧を取得する
func (h *Handler) AnnouncementsV1List(c *gin.Context) {
	announcements, err := h.announcementService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"announcements": ToAPIAnnouncements(announcements),
	})
}

// AnnouncementsV1Create 新規作成する
func (h *Handler) AnnouncementsV1Create(c *gin.Context) {
	var req api.AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	announcement, err := h.announcementService.Create(c.Request.Context(), ToDomainAnnouncementRequest(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"announcement": ToAPIAnnouncement(announcement),
	})
}

// AnnouncementsV1Delete 削除する
func (h *Handler) AnnouncementsV1Delete(c *gin.Context, id string) {
	if err := h.announcementService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// AnnouncementsV1Update 更新する
func (h *Handler) AnnouncementsV1Update(c *gin.Context, id string) {
	var req api.AnnouncementRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	announcement, err := h.announcementService.Update(c.Request.Context(), id, ToDomainAnnouncementRequest(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"announcement": ToAPIAnnouncement(announcement),
	})
}
