package server

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/cmd/config"
	"github.com/nsnikhil/stories/pkg/blog/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"testing"
)

func TestStoriesPing(t *testing.T) {
	server := newStoriesServer(newDeps(service.NewService(&service.MockStoriesService{}), config.LoadConfigs(), zap.NewExample()))
	resp, err := server.Ping(context.Background(), &proto.PingRequest{})
	require.NoError(t, err)
	assert.Equal(t, &proto.PingResponse{Message: "pong"}, resp)
}
