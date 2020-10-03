package handler

import (
	"github.com/nsnikhil/stories/pkg/http/internal/util"
	"net/http"
)

func PingHandler() http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		util.WriteSuccessResponse(http.StatusOK, "pong", resp)
	}
}
