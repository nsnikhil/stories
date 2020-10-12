package handler

import (
	"github.com/nsnikhil/stories/pkg/http/internal/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/util"
	"github.com/nsnikhil/stories/pkg/liberr"
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
		return liberr.WithOperation("DeleteStoryHandler.DeleteStory", err)
	}

	_, err = dsh.svc.DeleteStory(data.StoryID)
	if err != nil {
		return liberr.WithOperation("DeleteStoryHandler.DeleteStory", err)
	}

	//TODO: ADD SUCCESS LOG
	util.WriteSuccessResponse(http.StatusOK, contract.DeleteStoryResponse{Success: true}, resp)
	return nil
}

func NewDeleteStoryHandler(svc service.StoryService) *DeleteStoryHandler {
	return &DeleteStoryHandler{
		svc: svc,
	}
}
