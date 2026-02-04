package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"

	api "github.com/fun-dotto/api-template/generated"
)

// listAnnouncementsResponse represents the response from listing announcements
type listAnnouncementsResponse struct {
	Announcements []api.Announcement `json:"announcements"`
}

// announcementResponse represents the response containing a single announcement
type announcementResponse struct {
	Announcement api.Announcement `json:"announcement"`
}

// AnnouncementsV1List returns a list of announcements
// (GET /v1/announcements)
func (h *Handler) AnnouncementsV1List(c *gin.Context) {
	ctx := c.Request.Context()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, h.announcementAPIURL+"/v1/announcements", nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("upstream error: %s", string(body))})
		return
	}

	var result listAnnouncementsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"announcements": result.Announcements,
	})
}

// AnnouncementsV1Create creates a new announcement
// (POST /v1/announcements)
func (h *Handler) AnnouncementsV1Create(c *gin.Context) {
	ctx := c.Request.Context()

	var reqBody api.AnnouncementRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.createAnnouncement(ctx, &reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"announcement": result.Announcement,
	})
}

func (h *Handler) createAnnouncement(ctx context.Context, reqBody *api.AnnouncementRequest) (*announcementResponse, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, h.announcementAPIURL+"/v1/announcements", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result announcementResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// AnnouncementsV1Delete deletes an announcement by ID
// (DELETE /v1/announcements/{id})
func (h *Handler) AnnouncementsV1Delete(c *gin.Context, id string) {
	ctx := c.Request.Context()

	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, h.announcementAPIURL+"/v1/announcements/"+id, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.httpClient.Do(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		c.JSON(resp.StatusCode, gin.H{"error": fmt.Sprintf("upstream error: %s", string(body))})
		return
	}

	c.Status(http.StatusNoContent)
}

// AnnouncementsV1Update updates an announcement by ID
// (PUT /v1/announcements/{id})
func (h *Handler) AnnouncementsV1Update(c *gin.Context, id string) {
	ctx := c.Request.Context()

	var reqBody api.AnnouncementRequest
	if err := c.ShouldBindJSON(&reqBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.updateAnnouncement(ctx, id, &reqBody)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"announcement": result.Announcement,
	})
}

func (h *Handler) updateAnnouncement(ctx context.Context, id string, reqBody *api.AnnouncementRequest) (*announcementResponse, error) {
	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, h.announcementAPIURL+"/v1/announcements/"+id, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}

	var result announcementResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
