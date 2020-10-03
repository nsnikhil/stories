package stories

import (
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/story/model"
	"time"
)

func toProtoStory(st *model.Story) *proto.Story {
	return &proto.Story{
		Id:            st.GetID(),
		Title:         st.GetTitle(),
		Body:          st.GetBody(),
		Views:         st.GetViewCount(),
		UpVotes:       st.GetUpVotes(),
		DownVotes:     st.GetDownVotes(),
		CreatedAtUnix: st.GetCreatedAt().Unix(),
		UpdatedAtUnix: st.GetUpdatedAt().Unix(),
	}
}

func toDomainStory(cfg config.StoryConfig, st *proto.Story) (*model.Story, error) {
	return model.NewStoryBuilder().
		SetID(st.GetId()).
		SetTitle(cfg.TitleMaxLength(), st.GetTitle()).
		SetBody(cfg.BodyMaxLength(), st.GetBody()).
		SetViewCount(st.GetViews()).
		SetUpVotes(st.GetUpVotes()).
		SetDownVotes(st.GetDownVotes()).
		SetCreatedAt(time.Unix(st.GetCreatedAtUnix(), 0).UTC()).
		SetUpdatedAt(time.Unix(st.GetUpdatedAtUnix(), 0).UTC()).
		Build()
}
