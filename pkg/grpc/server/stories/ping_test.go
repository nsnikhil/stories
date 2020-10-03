package stories_test

import (
	"context"
	"github.com/nsnikhil/stories-proto/proto"
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/grpc/server/stories"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestStoriesPing(t *testing.T) {
	cfg := config.NewConfig()

	server := stories.NewStoriesServer(cfg.StoryConfig(), nil)
	resp, err := server.Ping(context.Background(), &proto.PingRequest{})

	require.NoError(t, err)
	assert.Equal(t, &proto.PingResponse{Message: "pong"}, resp)
}
