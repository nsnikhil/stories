package resperr_test

import (
	"github.com/bmizerany/assert"
	"github.com/nsnikhil/stories/pkg/http/internal/resperr"
	"net/http"
	"testing"
)

func TestGenericErrorGetErrorCode(t *testing.T) {
	ge := resperr.NewResponseError(http.StatusBadRequest, "some reason")

	assert.Equal(t, http.StatusBadRequest, ge.StatusCode())
	assert.Equal(t, "some reason", ge.Description())
}
