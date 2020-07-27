package utils

import (
	"github.com/nsnikhil/stories/pkg/blog/dao"
	"github.com/nsnikhil/stories/pkg/blog/dto"
)

func StoryDaoToDto(s *dao.Story) *dto.Story {
	return dto.NewStory()
}

func StoryDtoToDao(s *dto.Story) *dao.Story {
	return dao.NewStory()
}
