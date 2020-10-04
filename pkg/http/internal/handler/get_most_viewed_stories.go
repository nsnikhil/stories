package handler

import (
	"github.com/nsnikhil/stories/pkg/http/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/liberr"
	"github.com/nsnikhil/stories/pkg/http/internal/util"
	"github.com/nsnikhil/stories/pkg/story/service"
	"net/http"
)

type GetMostViewedStoriesHandler struct {
	svc service.StoryService
}

func (gmh *GetMostViewedStoriesHandler) GetMostViewedStories(resp http.ResponseWriter, req *http.Request) error {
	var data contract.MostViewedStoriesRequest
	err := util.ParseRequest(req, &data)
	if err != nil {
		return liberr.ValidationError(err.Error())
	}

	dss, err := gmh.svc.GetMostViewsStories(data.OffSet, data.Limit)
	if err != nil {
		return liberr.InternalError(err.Error())
	}

	// TODO: UNIFY BELOW LOGIC WITH TOP RATED HANDLER
	sz := len(dss)
	res := make([]contract.Story, sz)

	for i := 0; i < sz; i++ {
		res[i] = util.ConvertToDTO(&dss[i])
	}

	util.WriteSuccessResponse(http.StatusOK, res, resp)
	return nil
}

func NewGetMostViewedStoriesHandler(svc service.StoryService) *GetMostViewedStoriesHandler {
	return &GetMostViewedStoriesHandler{
		svc: svc,
	}
}
