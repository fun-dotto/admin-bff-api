package handler

import (
	api "github.com/fun-dotto/admin-bff-api/generated"
	"github.com/fun-dotto/admin-bff-api/generated/external/announcement_api"
	"github.com/fun-dotto/admin-bff-api/generated/external/subject_api"
)

type Handler struct {
	announcementClient *announcement_api.ClientWithResponses
	subjectClient      *subject_api.ClientWithResponses
}

func NewHandler(
	announcementClient *announcement_api.ClientWithResponses,
	subjectClient *subject_api.ClientWithResponses,
) *Handler {
	if announcementClient == nil {
		panic("announcementClient is required")
	}
	if subjectClient == nil {
		panic("subjectClient is required")
	}
	return &Handler{
		announcementClient: announcementClient,
		subjectClient:      subjectClient,
	}
}

var _ api.ServerInterface = (*Handler)(nil)
