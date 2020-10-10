package liberr

type Kind string

const (
	InternalError    Kind = "internalError"
	ResourceNotFound Kind = "resourceNotFound"
	ValidationError  Kind = "validationError"
)
