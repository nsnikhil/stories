package liberr

import (
	"fmt"
	"strings"
)

type Error struct {
	kind      Kind
	operation Operation
	severity  Severity
	cause     error
}

func (e *Error) Kind() Kind {
	k := e.kind
	if len(k) != 0 {
		return k
	}

	t, ok := e.cause.(*Error)
	if !ok {
		return k
	}

	return t.Kind()
}

//TODO: DECIDE THE FINAL IMPLEMENTATION OF THE DETAILS METHOD
func (e *Error) Details() string {
	b := &strings.Builder{}

	if len(e.kind) != 0 {
		b.WriteString(fmt.Sprintf(" kind: %v ", e.kind))
	}

	if len(e.operation) != 0 {
		b.WriteString(fmt.Sprintf(" operations: %v ", e.operation))
	}

	if len(e.severity) != 0 {
		b.WriteString(fmt.Sprintf(" severity: %v ", e.severity))
	}

	if e.cause != nil {
		b.WriteString(fmt.Sprintf(" cause: %v ", e.cause))
	}

	return strings.TrimSpace(b.String())
}

func (e *Error) Operation() Operation {
	return e.operation
}

func (e *Error) Operations() []Operation {
	ops := []Operation{e.operation}

	t, ok := e.cause.(*Error)
	if ok {
		ops = append(ops, t.Operations()...)
	}

	return ops
}

func (e *Error) Severity() Severity {
	s := e.severity
	if len(s) != 0 {
		return s
	}

	t, ok := e.cause.(*Error)
	if !ok {
		return s
	}

	return t.Severity()
}

func (e *Error) Unwrap() error {
	return e.cause
}

//TODO: DECIDE THE FINAL IMPLEMENTATION OF THE ERROR METHOD
func (e *Error) Error() string {
	if e.cause == nil {
		return ""
	}

	return e.cause.Error()
}

func NewError(kind Kind, operation Operation, severity Severity, cause error) *Error {
	return &Error{
		kind:      kind,
		operation: operation,
		severity:  severity,
		cause:     cause,
	}
}

func WithArgs(args ...interface{}) *Error {
	e := &Error{}

	for _, arg := range args {
		switch t := arg.(type) {
		case Operation:
			e.operation = t
		case Kind:
			e.kind = t
		case Severity:
			e.severity = t
		case error:
			e.cause = t
		default:
			return nil
		}
	}

	return e
}
