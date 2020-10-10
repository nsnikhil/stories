package handler

import (
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/http/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/util"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/nsnikhil/stories/pkg/story/service"
	"net/http"
)

type AddStoryHandler struct {
	cfg config.StoryConfig
	svc service.StoryService
}

func (ash *AddStoryHandler) AddStory(resp http.ResponseWriter, req *http.Request) error {
	var data contract.AddStoryRequest
	err := util.ParseRequest(req, &data)
	if err != nil {
		return err
	}

	st, err := model.NewStoryBuilder().
		SetTitle(ash.cfg.TitleMaxLength(), data.Title).
		SetBody(ash.cfg.BodyMaxLength(), data.Body).
		Build()

	if err != nil {
		return err
	}

	err = ash.svc.AddStory(st)
	if err != nil {
		return err
	}

	util.WriteSuccessResponse(http.StatusCreated, contract.AddStoryResponse{Success: true}, resp)
	return nil
}

func NewAddHandler(cfg config.StoryConfig, svc service.StoryService) *AddStoryHandler {
	return &AddStoryHandler{
		cfg: cfg,
		svc: svc,
	}
}
