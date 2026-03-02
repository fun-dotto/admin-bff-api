package handler

import (
	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/generated/external/announcement_api"
	"github.com/fun-dotto/api-template/generated/external/faculty_api"
	"github.com/fun-dotto/api-template/generated/external/subject_api"
)

type Handler struct {
	announcementClient *announcement_api.ClientWithResponses
	facultyClient      *faculty_api.ClientWithResponses
	subjectClient      *subject_api.ClientWithResponses
}

func NewHandler(
	announcementClient *announcement_api.ClientWithResponses,
	facultyClient *faculty_api.ClientWithResponses,
	subjectClient *subject_api.ClientWithResponses,
) *Handler {
	if announcementClient == nil {
		panic("announcementClient is required")
	}
	if facultyClient == nil {
		panic("facultyClient is required")
	}
	if subjectClient == nil {
		panic("subjectClient is required")
	}
	return &Handler{
		announcementClient: announcementClient,
		facultyClient:      facultyClient,
		subjectClient:      subjectClient,
	}
}

var _ api.ServerInterface = (*Handler)(nil)
