package liberr_test

import (
	"github.com/bmizerany/assert"
	"github.com/nsnikhil/stories/pkg/http/internal/liberr"
	"net/http"
	"testing"
)

func TestGenericErrorGetErrorCode(t *testing.T) {
	ge := liberr.NewResponseError("us-def", http.StatusBadRequest, "some reason")

	assert.Equal(t, "us-def", ge.ErrorCode())
	assert.Equal(t, http.StatusBadRequest, ge.StatusCode())
	assert.Equal(t, "some reason", ge.Error())
}
