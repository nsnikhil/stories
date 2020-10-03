package stories

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
)

func (ss *StoriesServer) GetMostViewedStories(ctx context.Context, req *proto.MostViewedStoriesRequest) (*proto.MostViewedStoriesResponse, error) {
	stories, err := ss.svc.GetMostViewsStories(int(req.GetOffset()), int(req.GetLimit()))
	if err != nil {
		return nil, err
	}

	sz := len(stories)
	resp := make([]*proto.Story, sz)

	for i := 0; i < sz; i++ {
		resp[i] = toProtoStory(&stories[i])
	}

	return &proto.MostViewedStoriesResponse{Stories: resp}, nil
}
