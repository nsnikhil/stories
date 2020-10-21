package stories

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/liberr"
)

func (ss *Server) GetMostViewedStories(ctx context.Context, req *proto.MostViewedStoriesRequest) (*proto.MostViewedStoriesResponse, error) {
	stories, err := ss.svc.GetMostViewsStories(int(req.GetOffset()), int(req.GetLimit()))
	if err != nil {
		return nil, liberr.WithArgs(liberr.Operation("Server.GetMostViewedStories"), err)
	}

	sz := len(stories)
	resp := make([]*proto.Story, sz)

	for i := 0; i < sz; i++ {
		resp[i] = toProtoStory(&stories[i])
	}

	//TODO: ADD SUCCESS LOG
	return &proto.MostViewedStoriesResponse{Stories: resp}, nil
}
