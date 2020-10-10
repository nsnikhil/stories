package liberr

import (
	"fmt"
	"strings"
)

type Error struct {
	kind       Kind
	operations []Operation
	severity   Severity
	err        error
}

func (e *Error) Kind() Kind {
	return e.kind
}

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

	t.err = err
	return t
}

func WithArgs(args ...interface{}) *Error {
	e := &Error{}
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
			return nil
		}
	}
	return e
}

func NewError(kind Kind, operations []Operation, severity Severity, err error) *Error {
	return &Error{
		kind:       kind,
		operations: operations,
		severity:   severity,
		err:        err,
	}
}

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
