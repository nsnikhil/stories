package resperr

import (
	"github.com/nsnikhil/stories/pkg/liberr"
	"net/http"
)

const (
	defaultStatusCode = http.StatusInternalServerError
	defaultMessage    = "internal server error"
	notFoundMessage   = "requested resource was not found"
)

func MapError(err error) ResponseError {
	t, ok := err.(*liberr.Error)
	if !ok {
		return NewResponseError(defaultStatusCode, defaultMessage)
	}

	k := t.Kind()
	switch k {
	case liberr.ValidationError:
		return NewResponseError(http.StatusBadRequest, t.Error())
	case liberr.ResourceNotFound:
		return NewResponseError(http.StatusNotFound, notFoundMessage)
	default:
		return NewResponseError(defaultStatusCode, defaultMessage)
	}
}
