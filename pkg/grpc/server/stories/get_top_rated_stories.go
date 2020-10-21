package stories

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/liberr"
)

func (ss *Server) GetTopRatedStories(ctx context.Context, req *proto.TopRatedStoriesRequest) (*proto.TopRatedStoriesResponse, error) {
	stories, err := ss.svc.GetTopRatedStories(int(req.GetOffset()), int(req.GetLimit()))
	if err != nil {
		return nil, liberr.WithArgs(liberr.Operation("Server.GetTopRatedStories"), err)
	}

	sz := len(stories)
	resp := make([]*proto.Story, sz)

	for i := 0; i < sz; i++ {
		resp[i] = toProtoStory(&stories[i])
	}

	//TODO: ADD SUCCESS LOG
	return &proto.TopRatedStoriesResponse{Stories: resp}, nil
}
