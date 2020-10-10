package handler

import (
	"github.com/nsnikhil/stories/pkg/http/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/util"
	"github.com/nsnikhil/stories/pkg/story/service"
	"net/http"
)

type GetStoryHandler struct {
	svc service.StoryService
}

func (gs *GetStoryHandler) GetStory(resp http.ResponseWriter, req *http.Request) error {
	var data contract.GetStoryRequest
	err := util.ParseRequest(req, &data)
	if err != nil {
		return err
	}

	st, err := gs.svc.GetStory(data.StoryID)
	if err != nil {
		return err
	}

	util.WriteSuccessResponse(http.StatusOK, util.ConvertToDTO(st), resp)
	return nil
}

func NewGetStoryHandler(svc service.StoryService) *GetStoryHandler {
	return &GetStoryHandler{
		svc: svc,
	}
}
