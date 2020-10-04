package handler

import (
	"errors"
	"github.com/nsnikhil/stories/pkg/story/service"
	"net/http"
)

type SearchStoriesHandler struct {
	svc service.StoryService
}

func (ssh *SearchStoriesHandler) SearchStories(resp http.ResponseWriter, req *http.Request) error {
	return errors.New("UNIMPLEMENTED")
}

func NewSearchStoriesHandler(svc service.StoryService) *SearchStoriesHandler {
	return &SearchStoriesHandler{
		svc: svc,
	}
}
