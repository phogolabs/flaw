package flaw

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	_ error          = &Error{}
	_ json.Marshaler = &Error{}
)

// ErrorData represents the error's data
type ErrorData struct {
	Code    int             `json:"code,omitempty"`
	Message string          `json:"message,omitempty"`
	Details []string        `json:"details,omitempty"`
	Reason  json.RawMessage `json:"reason,omitempty"`
}

// Error represents a wrapped error
type Error struct {
	code    int
	status  int
	msg     string
	details []string
	stack   StackTrace
	reason  error
}

// New creates a new error
func New(msg string, details ...string) *Error {
	return &Error{
		msg:     msg,
		stack:   NewStackTrace(),
		details: details,
	}
}

// Wrap wraps an error
func Wrap(err error) *Error {
	var errx *Error

	if !errors.As(err, &errx) {
		errx = &Error{
			reason: err,
			stack:  NewStackTrace(),
		}
	}

	return errx
}

// WithError creates an error copy with given error wrapped
func (x Error) WithError(err error) *Error {
	x.reason = err
	x.stack = NewStackTrace()
	return &x
}

// WithMessage creates an error copy with given message
func (x Error) WithMessage(text string) *Error {
	x.msg = text
	return &x
}

// WithStatus creates an error copy with given status
func (x Error) WithStatus(status int) *Error {
	x.status = status
	return &x
}

// WithCode creates an error copy with given status
func (x Error) WithCode(code int) *Error {
	x.code = code
	return &x
}

// Code returns the error code
func (x *Error) Code() int {
	return x.code
}

// Status returns the error status
func (x *Error) Status() int {
	return x.status
}

// Wrap wraps the given error
func (x *Error) Wrap(err error) {
	x.stack = NewStackTrace()
	x.reason = err
}

// Unwrap unwraps the underlying error
func (x *Error) Unwrap() error {
	return x.reason
}

// StackTrace returns the stack trace where the error occurred
func (x *Error) StackTrace() StackTrace {
	return x.stack
}

// Error returns the error message
func (x *Error) Error() string {
	return fmt.Sprintf("%v", x)
}

// Format the error as string
func (x *Error) Format(state fmt.State, verb rune) {
	separate := func() {
		separator := " "

		if verb == 'v' && state.Flag('+') {
			separator = "\n"
		}

		fmt.Fprint(state, separator)
	}

	if x.code != 0 {
		fmt.Fprintf(state, "code: %d", x.code)
		separate()
	}

	if x.msg != "" {
		fmt.Fprintf(state, "message: %s", x.msg)
		separate()
	}

	if x.reason != nil {
		fmt.Fprint(state, "reason: ")

		if _, ok := x.reason.(ErrorCollector); ok {
			separate()
		}

		if formatter, ok := x.reason.(fmt.Formatter); ok {
			formatter.Format(state, verb)
		} else {
			fmt.Fprint(state, x.reason.Error())
		}

		if _, ok := x.reason.(ErrorCollector); !ok {
			separate()
		}
	}

	if x.stack != nil {
		fmt.Fprint(state, "stack:")
		separate()

		x.stack.Format(state, verb)
	}
}

// MarshalJSON marshals the error as json
func (x *Error) MarshalJSON() ([]byte, error) {
	data := &ErrorData{
		Code:    x.code,
		Message: x.msg,
		Details: x.details,
	}

	if x.reason != nil {
		var input interface{} = x.reason

		if _, ok := x.reason.(json.Marshaler); !ok {
			input = x.reason.Error()
		}

		buffer, err := json.Marshal(input)
		if err != nil {
			return nil, err
		}

		data.Reason = buffer
	}

	return json.Marshal(data)
}

var _ error = ErrorCollector{}

// ErrorCollector is a slice of errors
type ErrorCollector []error

// Error returns the error message
func (errs ErrorCollector) Error() string {
	return fmt.Sprintf("%v", errs)
}

// MarshalJSON marshals the error as json
func (errs ErrorCollector) MarshalJSON() ([]byte, error) {
	input := make([]interface{}, len(errs))

	for index, err := range errs {
		if _, ok := err.(json.Marshaler); ok {
			input[index] = err
		} else {
			input[index] = err.Error()
		}
	}

	return json.Marshal(input)
}

// Format the error as string
func (errs ErrorCollector) Format(state fmt.State, verb rune) {
	// write as singleline string
	var (
		separator string
		prefix    string
		count     = len(errs)
	)

	if count > 1 {
		separator = "; "
	}

	if verb == 'v' && state.Flag('+') {
		separator = "\n"
		prefix = " --- "
	}

	for _, item := range errs {
		fmt.Fprint(state, prefix)

		if formatter, ok := item.(fmt.Formatter); ok {
			formatter.Format(state, verb)
		} else {
			fmt.Fprint(state, item.Error())
		}

		fmt.Fprint(state, separator)
	}
}

// Is reports whether any error in err's chain matches target.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error is considered to match a target if it is equal to that target or if
// it implements a method Is(error) bool such that Is(target) returns true.
func (errs ErrorCollector) Is(target error) bool {
	items, ok := target.(ErrorCollector)

	if !ok {
		items = ErrorCollector{target}
	}

	if len(errs) != len(items) {
		return false
	}

	for index, child := range errs {
		if !errors.Is(child, items[index]) {
			return false
		}
	}

	return true
}

// As finds the first error in err's chain that matches target, and if so, sets
// target to that error value and returns true.
//
// The chain consists of err itself followed by the sequence of errors obtained by
// repeatedly calling Unwrap.
//
// An error matches target if the error's concrete value is assignable to the value
// pointed to by target, or if the error has a method As(interface{}) bool such that
// As(target) returns true. In the latter case, the As method is responsible for
// setting target.
//
// As will panic if target is not a non-nil pointer to either a type that implements
// error, or to any interface type. As returns false if err is nil.
func (errs ErrorCollector) As(err interface{}) bool {
	for _, child := range errs {
		if errors.As(child, err) {
			return true
		}
	}

	return false
}

// Unwrap unwraps the underlying error
func (errs ErrorCollector) Unwrap() error {
	count := len(errs)

	switch {
	case count == 1:
		return errs[0]
	default:
		return nil
	}
}
