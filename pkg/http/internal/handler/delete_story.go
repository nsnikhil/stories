package handler

import (
	"github.com/nsnikhil/stories/pkg/http/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/liberr"
	"github.com/nsnikhil/stories/pkg/http/internal/util"
	"github.com/nsnikhil/stories/pkg/story/service"
	"net/http"
)

type DeleteStoryHandler struct {
	svc service.StoryService
}

func (dsh *DeleteStoryHandler) DeleteStory(resp http.ResponseWriter, req *http.Request) error {
	var data contract.DeleteStoryRequest
	err := util.ParseRequest(req, &data)
	if err != nil {
		return liberr.ValidationError(err.Error())
	}

	_, err = dsh.svc.DeleteStory(data.StoryID)
	if err != nil {
		return liberr.InternalError(err.Error())
	}

	util.WriteSuccessResponse(http.StatusOK, contract.DeleteStoryResponse{Success: true}, resp)
	return nil
}

func NewDeleteStoryHandler(svc service.StoryService) *DeleteStoryHandler {
	return &DeleteStoryHandler{
		svc: svc,
	}
}
