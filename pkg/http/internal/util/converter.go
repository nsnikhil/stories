package util

import (
	"github.com/nsnikhil/stories/pkg/http/contract"
	"github.com/nsnikhil/stories/pkg/story/model"
	"time"
)

func ConvertToDTO(st *model.Story) contract.Story {
	return contract.Story{
		ID:        st.GetID(),
		Title:     st.GetTitle(),
		Body:      st.GetBody(),
		ViewCount: st.GetViewCount(),
		UpVotes:   st.GetUpVotes(),
		DownVotes: st.GetDownVotes(),
		CreatedAt: st.GetCreatedAt().Unix(),
		UpdatedAt: st.GetUpdatedAt().Unix(),
	}
}

func ConvertToDAO(titleMaxLength, bodyMaxLength int, st contract.Story) (*model.Story, error) {
	return model.NewStoryBuilder().
		SetID(st.ID).
		SetTitle(titleMaxLength, st.Title).
		SetBody(bodyMaxLength, st.Body).
		SetViewCount(st.ViewCount).
		SetUpVotes(st.UpVotes).
		SetDownVotes(st.DownVotes).
		SetCreatedAt(time.Unix(st.CreatedAt, 0).UTC()).
		SetUpdatedAt(time.Unix(st.UpdatedAt, 0).UTC()).
		Build()
}
