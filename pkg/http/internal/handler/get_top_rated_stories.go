package handler

import (
	"github.com/nsnikhil/stories/pkg/http/internal/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/util"
	"github.com/nsnikhil/stories/pkg/liberr"
	"github.com/nsnikhil/stories/pkg/story/service"
	"net/http"
)

type GetTopRatedStoriesHandler struct {
	svc service.StoryService
}

func (gmh *GetTopRatedStoriesHandler) GetTopRatedStories(resp http.ResponseWriter, req *http.Request) error {
	var data contract.TopRatedStoriesRequest
	err := util.ParseRequest(req, &data)
	if err != nil {
		return liberr.WithOperation("GetTopRatedStoriesHandler.GetTopRatedStories", err)
	}

	dss, err := gmh.svc.GetTopRatedStories(data.OffSet, data.Limit)
	if err != nil {
		return liberr.WithOperation("GetTopRatedStoriesHandler.GetTopRatedStories", err)
	}

	// TODO: UNIFY BELOW LOGIC WITH MOST VIEWED HANDLER
	sz := len(dss)
	res := make([]contract.Story, sz)

	for i := 0; i < sz; i++ {
		res[i] = util.ConvertToDTO(&dss[i])
	}

	//TODO: ADD SUCCESS LOG
	util.WriteSuccessResponse(http.StatusOK, res, resp)
	return nil
}

func NewGetTopRatedStoriesHandler(svc service.StoryService) *GetTopRatedStoriesHandler {
	return &GetTopRatedStoriesHandler{
		svc: svc,
	}
}
