package util

import (
	"encoding/json"
	"errors"
	"github.com/nsnikhil/stories/pkg/liberr"
	"io/ioutil"
	"net/http"
)

func ParseRequest(req *http.Request, data interface{}) error {
	if req == nil {
		return liberr.WithArgs(liberr.SeverityError, liberr.ValidationError, liberr.Operation("ParseRequest"), errors.New("request is nil"))
	}

	if req.Body == nil {
		return liberr.WithArgs(liberr.SeverityError, liberr.ValidationError, liberr.Operation("ParseRequest"), errors.New("request body is nil"))
	}

	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return liberr.WithArgs(liberr.SeverityError, liberr.ValidationError, liberr.Operation("ParseRequest.ioutil.ReadAll"), err)
	}

	err = json.Unmarshal(b, &data)
	if err != nil {
		return liberr.WithArgs(liberr.SeverityError, liberr.ValidationError, liberr.Operation("ParseRequest.json.Unmarshal"), err)
	}

	return nil
}
