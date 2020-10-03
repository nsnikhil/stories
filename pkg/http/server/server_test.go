package server_test

import (
	"github.com/nsnikhil/stories/pkg/config"
	"github.com/nsnikhil/stories/pkg/http/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"net/http"
	"testing"
	"time"
)

func TestServerStart(t *testing.T) {
	cfg := config.NewConfig()
	lgr := zap.NewNop()

	rt := http.NewServeMux()
	rt.HandleFunc("/ping", func(resp http.ResponseWriter, req *http.Request) {})

	srv := server.NewServer(cfg, lgr, rt)
	go srv.Start()

	//TODO REMOVE SLEEP
	time.Sleep(time.Millisecond)

	resp, err := http.Get("http://127.0.0.1:8081/ping")
	require.NoError(t, err)

	assert.Equal(t, http.StatusOK, resp.StatusCode)

}
