package resperr

import (
	"fmt"
	"github.com/nsnikhil/stories/pkg/liberr"
	"net/http"
)

const (
	defaultStatusCode = http.StatusInternalServerError
	defaultMessage    = "internal server error"
	notFoundMessage   = "requested resource was not found"
)

func MapError(err error) HTTPResponseError {
	t, ok := err.(*liberr.Error)
	if !ok {
		return NewHTTPResponseError(defaultStatusCode, defaultMessage)
	}

	fmt.Println(err)

	k := t.Kind()
	switch k {
	case liberr.ValidationError:
		return NewHTTPResponseError(http.StatusBadRequest, t.Error())
	case liberr.ResourceNotFound:
		return NewHTTPResponseError(http.StatusNotFound, notFoundMessage)
	default:
		return NewHTTPResponseError(defaultStatusCode, defaultMessage)
	}
}
