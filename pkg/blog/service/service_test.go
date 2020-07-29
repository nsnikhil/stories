package service

import (
	"github.com/nsnikhil/stories/pkg/blog/store"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"testing"
)

func TestCreateNewService(t *testing.T) {
	ss := NewDefaultStoriesService(store.NewStore(&store.MockStoriesStore{}, &store.MockStoriesCache{}), zap.NewExample())
	actualResult := NewService(ss)
	expectedResult := &Service{ss: ss}
	assert.Equal(t, expectedResult, actualResult)
}

func TestServiceGetStoriesService(t *testing.T) {
	ss := NewDefaultStoriesService(store.NewStore(&store.MockStoriesStore{}, &store.MockStoriesCache{}), zap.NewExample())
	s := NewService(ss)
	assert.Equal(t, ss, s.GetStoriesService())
}
