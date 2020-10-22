package liberr

import (
	"errors"
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

func (e *Error) Is(target error) bool {
	return errors.Is(e.cause, target)
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

func (e *Error) EncodedStack() string {
	var encode func(d map[string]interface{}) string
	encode = func(d map[string]interface{}) string {
		b := new(strings.Builder)

		for k, v := range d {
			t, ok := v.(map[string]interface{})

			var val string

			if ok {
				val = encode(t)
			} else {
				val = fmt.Sprintf("%s", v)
			}

			b.WriteString(fmt.Sprintf(" %s:%s ", k, val))
		}

		return fmt.Sprintf("[%s]", strings.TrimSpace(b.String()))
	}

	return encode(e.Stack())
}

func (e *Error) Stack() map[string]interface{} {
	res := make(map[string]interface{})

	if len(e.kind) != 0 {
		res["kind"] = string(e.kind)
	}

	if len(e.operation) != 0 {
		res["operation"] = string(e.operation)
	}

	if len(e.severity) != 0 {
		res["severity"] = string(e.severity)
	}

	if e.cause != nil {
		t, ok := e.cause.(*Error)
		if ok {
			res["cause"] = t.Stack()
		} else {
			res["cause"] = e.cause.Error()
		}
	}

	return res
}
