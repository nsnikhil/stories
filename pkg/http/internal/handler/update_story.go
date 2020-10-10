package handler

import (
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/http/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/util"
	"github.com/nsnikhil/stories/pkg/story/service"
	"net/http"
)

type UpdateStoryHandler struct {
	cfg config.StoryConfig
	svc service.StoryService
}

func (ush *UpdateStoryHandler) UpdateStory(resp http.ResponseWriter, req *http.Request) error {
	var data contract.UpdateStoryRequest
	err := util.ParseRequest(req, &data)
	if err != nil {
		return err
	}

	st, err := util.ConvertToDAO(ush.cfg.TitleMaxLength(), ush.cfg.BodyMaxLength(), data.Story)
	if err != nil {
		return err
	}

	_, err = ush.svc.UpdateStory(st)
	if err != nil {
		return err
	}

	util.WriteSuccessResponse(http.StatusOK, contract.UpdateStoryResponse{Success: true}, resp)
	return nil
}

func NewUpdateStoryHandler(cfg config.StoryConfig, svc service.StoryService) *UpdateStoryHandler {
	return &UpdateStoryHandler{
		cfg: cfg,
		svc: svc,
	}
}
