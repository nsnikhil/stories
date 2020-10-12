package liberr

type Kind string

//TODO: ADD NEW KINDS
const (
	InternalError    Kind = "internalError"
	ResourceNotFound Kind = "resourceNotFound"
	ValidationError  Kind = "validationError"
)
