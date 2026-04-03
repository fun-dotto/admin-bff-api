package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/admin-bff-api/generated"
	"github.com/fun-dotto/admin-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/admin-bff-api/internal/middleware"
)

// NotifyIrregularitiesV1Notify 時間割変更を通知する
func (h *Handler) NotifyIrregularitiesV1Notify(c *gin.Context, params api.NotifyIrregularitiesV1NotifyParams) {
	if !middleware.RequireAnyClaim(c, "admin", "developer") {
		return
	}

	response, err := h.academicClient.NotifyIrregularitiesV1NotifyWithResponse(c.Request.Context(), &academic_api.NotifyIrregularitiesV1NotifyParams{
		UserIds: params.UserIds,
		Date:    params.Date,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if response.StatusCode() != http.StatusNoContent {
		c.JSON(response.StatusCode(), gin.H{"error": "unexpected response from upstream"})
		return
	}

	c.Status(http.StatusNoContent)
}
