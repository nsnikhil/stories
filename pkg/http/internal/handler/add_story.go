package handler

import (
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/http/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/liberr"
	"github.com/nsnikhil/stories/pkg/http/internal/util"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/nsnikhil/stories/pkg/story/service"
	"net/http"
)

type AddHandler struct {
	cfg config.StoryConfig
	svc service.StoryService
}

func (ash *AddHandler) Add(resp http.ResponseWriter, req *http.Request) error {
	var data contract.AddStoryRequest
	err := util.ParseRequest(req, &data)
	if err != nil {
		return liberr.ValidationError(err.Error())
	}

	st, err := model.NewStoryBuilder().
		SetTitle(ash.cfg.TitleMaxLength(), data.Title).
		SetBody(ash.cfg.BodyMaxLength(), data.Body).
		Build()

	if err != nil {
		return liberr.ValidationError(err.Error())
	}

	err = ash.svc.AddStory(st)
	if err != nil {
		return liberr.InternalError(err.Error())
	}

	util.WriteSuccessResponse(http.StatusCreated, contract.AddStoryResponse{Success: true}, resp)
	return nil
}

func NewAddHandler(cfg config.StoryConfig, svc service.StoryService) *AddHandler {
	return &AddHandler{
		cfg: cfg,
		svc: svc,
	}
}
