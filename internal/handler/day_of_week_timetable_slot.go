package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/internal/middleware"
)

// DayOfWeekTimetableSlotsV1List 曜日・時限一覧を取得する
func (h *Handler) DayOfWeekTimetableSlotsV1List(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	slots, err := h.dayOfWeekTimetableSlotService.List(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.DayOfWeekTimetableSlotsV1List200JSONResponse{
		DayOfWeekTimetableSlots: ToAPIDayOfWeekTimetableSlots(slots),
	})
}

// DayOfWeekTimetableSlotsV1Detail 曜日・時限を詳細取得する
func (h *Handler) DayOfWeekTimetableSlotsV1Detail(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	slot, err := h.dayOfWeekTimetableSlotService.Detail(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.DayOfWeekTimetableSlotsV1Detail200JSONResponse{
		DayOfWeekTimetableSlot: ToAPIDayOfWeekTimetableSlot(slot),
	})
}

// DayOfWeekTimetableSlotsV1Create 曜日・時限を作成する
func (h *Handler) DayOfWeekTimetableSlotsV1Create(c *gin.Context) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req api.SubjectServiceDayOfWeekTimetableSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slot, err := h.dayOfWeekTimetableSlotService.Create(c.Request.Context(), ToDomainDayOfWeekTimetableSlotRequest(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, api.DayOfWeekTimetableSlotsV1Create201JSONResponse{
		DayOfWeekTimetableSlot: ToAPIDayOfWeekTimetableSlot(slot),
	})
}

// DayOfWeekTimetableSlotsV1Update 曜日・時限を更新する
func (h *Handler) DayOfWeekTimetableSlotsV1Update(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	var req api.SubjectServiceDayOfWeekTimetableSlotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	slot, err := h.dayOfWeekTimetableSlotService.Update(c.Request.Context(), id, ToDomainDayOfWeekTimetableSlotRequest(&req))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, api.DayOfWeekTimetableSlotsV1Update200JSONResponse{
		DayOfWeekTimetableSlot: ToAPIDayOfWeekTimetableSlot(slot),
	})
}

// DayOfWeekTimetableSlotsV1Delete 曜日・時限を削除する
func (h *Handler) DayOfWeekTimetableSlotsV1Delete(c *gin.Context, id string) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	if err := h.dayOfWeekTimetableSlotService.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
	c.Writer.WriteHeaderNow()
}
