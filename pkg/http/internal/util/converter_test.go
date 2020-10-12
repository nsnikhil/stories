package util_test

import (
	"github.com/nsnikhil/stories/pkg/http/internal/contract"
	"github.com/nsnikhil/stories/pkg/http/internal/util"
	"github.com/nsnikhil/stories/pkg/story/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConvertDomainToRequestStory(t *testing.T) {
	createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

	ds, err := model.NewStoryBuilder().
		SetID("adbca278-7e5c-4831-bf90-15fadfda0dd1").
		SetTitle(100, "title").
		SetBody(100, "test body").
		SetViewCount(25).
		SetUpVotes(10).
		SetDownVotes(2).
		SetCreatedAt(createdAt).
		SetUpdatedAt(updatedAt).
		Build()

	require.NoError(t, err)

	rs := contract.Story{
		ID:        "adbca278-7e5c-4831-bf90-15fadfda0dd1",
		Title:     "title",
		Body:      "test body",
		ViewCount: 25,
		UpVotes:   10,
		DownVotes: 2,
		CreatedAt: createdAt.Unix(),
		UpdatedAt: updatedAt.Unix(),
	}

	assert.Equal(t, rs, util.ConvertToDTO(ds))
}

func TestConvertRequestToDomainStory(t *testing.T) {
	createdAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)
	updatedAt := time.Date(2020, 07, 29, 16, 0, 0, 0, time.UTC)

	ds, err := model.NewStoryBuilder().
		SetID("adbca278-7e5c-4831-bf90-15fadfda0dd1").
		SetTitle(100, "title").
		SetBody(100, "test body").
		SetViewCount(25).
		SetUpVotes(10).
		SetDownVotes(2).
		SetCreatedAt(createdAt).
		SetUpdatedAt(updatedAt).
		Build()

	require.NoError(t, err)

	rs := contract.Story{
		ID:        "adbca278-7e5c-4831-bf90-15fadfda0dd1",
		Title:     "title",
		Body:      "test body",
		ViewCount: 25,
		UpVotes:   10,
		DownVotes: 2,
		CreatedAt: createdAt.Unix(),
		UpdatedAt: updatedAt.Unix(),
	}

	res, err := util.ConvertToDAO(100, 10000, rs)
	require.NoError(t, err)

	assert.Equal(t, ds, res)
}
