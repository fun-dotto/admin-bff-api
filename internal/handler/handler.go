package handler

import (
	api "github.com/fun-dotto/api-template/generated"
	"github.com/fun-dotto/api-template/generated/external/announcement_api"
	"github.com/fun-dotto/api-template/generated/external/subject_api"
)

type Handler struct {
	announcementClient *announcement_api.ClientWithResponses
	subjectClient      *subject_api.ClientWithResponses
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) WithAnnouncementClient(c *announcement_api.ClientWithResponses) *Handler {
	h.announcementClient = c
	return h
}

func (h *Handler) WithSubjectClient(c *subject_api.ClientWithResponses) *Handler {
	h.subjectClient = c
	return h
}

var _ api.ServerInterface = (*Handler)(nil)
