package resperr

type HTTPResponseError struct {
	statusCode  int
	description string
}

func (re HTTPResponseError) StatusCode() int {
	return re.statusCode
}

func (re HTTPResponseError) Description() string {
	return re.description
}

func NewHTTPResponseError(statusCode int, description string) HTTPResponseError {
	return HTTPResponseError{
		statusCode:  statusCode,
		description: description,
	}
}
