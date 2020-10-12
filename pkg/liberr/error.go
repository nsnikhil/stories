package liberr

import (
	"fmt"
	"strings"
)

type Error struct {
	kind       Kind
	operations []Operation
	//TODO: WHY IS SEVERITY NEEDED ?
	severity Severity
	err      error
}

func (e *Error) Kind() Kind {
	return e.kind
}

//TODO: DECIDE THE FINAL IMPLEMENTATION OF THE DETAILS METHOD
func (e *Error) Details() string {
	b := &strings.Builder{}

	if len(e.kind) != 0 {
		b.WriteString(fmt.Sprintf(" kind: %v ", e.kind))
	}

	if e.operations != nil && len(e.operations) != 0 {
		b.WriteString(fmt.Sprintf(" operations: %v ", e.operations))
	}

	if len(e.severity) != 0 {
		b.WriteString(fmt.Sprintf(" severity: %v ", e.severity))
	}

	if e.err != nil {
		b.WriteString(fmt.Sprintf(" cause: %v ", e.err))
	}

	return strings.TrimSpace(b.String())
}

func (e *Error) Operations() []Operation {
	return e.operations
}

func (e *Error) Severity() Severity {
	return e.severity
}

func (e *Error) Unwrap() error {
	return e.err
}

//TODO: DECIDE THE FINAL IMPLEMENTATION OF THE ERROR METHOD
func (e *Error) Error() string {
	if e.err == nil {
		return ""
	}

	return e.err.Error()
}

func WithKind(kind Kind, err error) *Error {
	if err == nil {
		return nil
	}

	t, ok := err.(*Error)
	if !ok {
		return WithArgs(kind, err)
	}

	//TODO: SHOULD OVERRIDING BE ALLOWED
	t.kind = kind
	return t
}

func WithOperation(operation Operation, err error) *Error {
	if err == nil {
		return nil
	}

	t, ok := err.(*Error)
	if !ok {
		return WithArgs(operation, err)
	}

	t.operations = append(t.operations, operation)
	return t
}

func WithSeverity(severity Severity, err error) *Error {
	if err == nil {
		return nil
	}

	t, ok := err.(*Error)
	if !ok {
		return WithArgs(severity, err)
	}

	//TODO: SHOULD OVERRIDING BE ALLOWED
	t.severity = severity
	return t
}

func WithCause(err error) *Error {
	if err == nil {
		return nil
	}

	t, ok := err.(*Error)
	if !ok {
		return WithArgs(err)
	}

	//TODO: SHOULD OVERRIDING BE ALLOWED
	t.err = err
	return t
}

// TODO: DECIDE WHICH OF THE THREE FUNC WOULD YOU BE USING TO CREATE NEW ERROR
//TODO: CREATE A SIMILAR FUNCTIONS WHICH DOES NOT CREATE NEW ERROR BUT RATHER INJECTS ARGS
func WithArgs(args ...interface{}) *Error {
	e := &Error{}
	ok := injectArgs(e, args...)
	if !ok {
		return nil
	}

	return e
}

//TODO: RENAME THIS TO SOMETHING MEANINGFUL
func WithInjectArgs(err error, args ...interface{}) *Error {
	if err == nil {
		return nil
	}

	t, ok := err.(*Error)
	if !ok {
		return WithArgs(err, args)
	}

	//TODO: SHOULD IT RETURN NIL IF THE ARGS IS INVALID AS WILL LOSE THE ERROR?
	ok = injectArgs(t, args...)
	if !ok {
		return nil
	}

	return t
}

func injectArgs(e *Error, args ...interface{}) bool {
	for _, arg := range args {
		switch t := arg.(type) {
		case Operation:
			e.operations = append(e.operations, t)
		case Kind:
			e.kind = t
		case Severity:
			e.severity = t
		case error:
			e.err = t
		default:
			return false
		}
	}

	return true
}

// TODO: DECIDE WHICH OF THE THREE FUNC WOULD YOU BE USING TO CREATE NEW ERROR
func NewError(kind Kind, operations []Operation, severity Severity, err error) *Error {
	return &Error{
		kind:       kind,
		operations: operations,
		severity:   severity,
		err:        err,
	}
}

// TODO: DECIDE WHICH OF THE THREE FUNC WOULD YOU BE USING TO CREATE NEW ERROR
type ErrorBuilder struct {
	kind       Kind
	operations []Operation
	severity   Severity
	err        error
}

func Builder() *ErrorBuilder {
	return &ErrorBuilder{}
}

func (eb *ErrorBuilder) WithKind(kind Kind) *ErrorBuilder {
	eb.kind = kind
	return eb
}

func (eb *ErrorBuilder) WithSeverity(severity Severity) *ErrorBuilder {
	eb.severity = severity
	return eb
}

func (eb *ErrorBuilder) WithOperation(operations Operation) *ErrorBuilder {
	eb.operations = append(eb.operations, operations)
	return eb
}

func (eb *ErrorBuilder) WithCause(err error) *ErrorBuilder {
	eb.err = err
	return eb
}

func (eb *ErrorBuilder) Build() *Error {
	return NewError(eb.kind, eb.operations, eb.severity, eb.err)
}
