package handler

import (
	api "github.com/fun-dotto/admin-bff-api/generated"
	"github.com/fun-dotto/admin-bff-api/generated/external/academic_api"
	"github.com/fun-dotto/admin-bff-api/generated/external/announcement_api"
)

type Handler struct {
	academicClient     *academic_api.ClientWithResponses
	announcementClient *announcement_api.ClientWithResponses
}

func NewHandler(
	academicClient *academic_api.ClientWithResponses,
	announcementClient *announcement_api.ClientWithResponses,
) *Handler {
	if academicClient == nil {
		panic("academicClient is required")
	}
	if announcementClient == nil {
		panic("announcementClient is required")
	}
	return &Handler{
		academicClient:     academicClient,
		announcementClient: announcementClient,
	}
}

var _ api.ServerInterface = (*Handler)(nil)
