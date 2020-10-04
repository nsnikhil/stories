package util

import (
	"github.com/nsnikhil/stories/pkg/http/contract"
	"github.com/nsnikhil/stories/pkg/story/model"
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
