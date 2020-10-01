package grpc

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	config2 "github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/story"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestStoriesPing(t *testing.T) {
	server := newStoriesServer(newDeps(story.NewService(&story.MockStoriesService{}), config2.LoadConfigs(), zap.NewExample()))
	resp, err := server.Ping(context.Background(), &proto.PingRequest{})
	require.NoError(t, err)
	assert.Equal(t, &proto.PingResponse{Message: "pong"}, resp)
}
