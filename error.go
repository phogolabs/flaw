package flaw

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/phogolabs/flaw/format"
)

const (
	keyCode    = "error_code"
	keyMessage = "error_message"
	keyDetails = "error_details"
	keyCause   = "error_cause"
	keyStack   = "error_stack"
)

var (
	_ error          = &Error{}
	_ json.Marshaler = &Error{}
)

// Map is an alias to map[string]interface{}
type Map = map[string]interface{}

// ErrorConstant represents an error that can create a constant / sentinel
// error such as io.EOF
type ErrorConstant string

// Error returns the error message
func (x ErrorConstant) Error() string {
	return fmt.Sprintf("%v", x)
}

// Error represents a wrapped error
type Error struct {
	code    int
	status  int
	msg     string
	details format.StringSlice
	stack   StackTrace
	context map[string]interface{}
	reason  error
}

// New creates a new error
func New(msg string, details ...string) *Error {
	return &Error{
		status:  500,
		msg:     msg,
		details: details,
		context: Map{},
		stack:   NewStackTrace(),
	}
}

// Wrap wraps an error
func Wrap(err error) *Error {
	var errx *Error

	if !errors.As(err, &errx) {
		errx = &Error{
			status:  500,
			reason:  err,
			context: Map{},
			stack:   NewStackTrace(),
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

// WithCode creates an error copy with given status
func (x Error) WithCode(code int) *Error {
	x.code = code
	return &x
}

// WithStatus creates an error copy with given status
func (x Error) WithStatus(status int) *Error {
	x.status = status
	return &x
}

// WithContext creates an error copy with given map
func (x Error) WithContext(context Map) *Error {
	if context == nil {
		context = Map{}
	}

	x.context = context
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

// Error returns the error message
func (x *Error) Error() string {
	return fmt.Sprintf("%v", x)
}

// Cause returns the underlying error
func (x *Error) Cause() error {
	return x.reason
}

// StackTrace returns the stack trace where the error occurred
func (x *Error) StackTrace() StackTrace {
	return x.stack
}

// Context returns the error's context
func (x *Error) Context() Map {
	return x.data()
}

// Format formats the frame according to the fmt.Formatter interface.
//
//    %m    error message
//    %c    error code
//    %r    error reason
//    %v    code: %d message: %s reason: %w
//
// Format accepts flags that alter the printing of some verbs, as follows:
//
//    %+s   stack trace
//    %+v   equivalent
func (x *Error) Format(state fmt.State, verb rune) {
	switch verb {
	case 'c':
		fmt.Fprintf(state, "%d", x.code)
	case 'm':
		fmt.Fprintf(state, "%s", x.msg)
	case 'r':
		fmt.Fprintf(state, "%v", x.reason)
	case 'd':
		x.details.Format(state, 'v')
	case 's':
		x.stack.Format(state, 'v')
	case 'v':
		formatter := format.NewState(state)
		defer formatter.Flush()

		title := func(text string) {
			if formatter.Size() > 0 {
				if formatter.Flag('+') {
					fmt.Fprint(formatter, "\n")
				} else {
					fmt.Fprint(formatter, " ")
				}
			}

			fmt.Fprint(formatter, text)

			if formatter.Flag('+') {
				fmt.Fprint(formatter, "\t")
			}

			fmt.Fprint(formatter, " ")
		}

		newline := func() {
			if formatter.Flag('+') {
				fmt.Fprint(formatter, "\n")
			}
		}

		if x.code != 0 {
			title("code:")
			x.Format(formatter, 'c')
		}

		if x.msg != "" {
			title("message:")
			x.Format(formatter, 'm')
		}

		if x.details != nil {
			title("details:")
			newline()
			x.Format(formatter, 'd')
		}

		if x.reason != nil {
			title("cause:")
			x.Format(formatter, 'r')
		}

		if x.stack != nil && state.Flag('+') {
			title("stack:")
			newline()
			x.Format(formatter, 's')
		}
	}
}

// MarshalJSON marshals the error as json
func (x *Error) MarshalJSON() ([]byte, error) {
	data := x.data(keyStack)

	if x.reason != nil {
		if _, ok := x.reason.(json.Marshaler); ok {
			data[keyCause] = x.reason
		}
	}

	return json.Marshal(data)
}

func (x *Error) data(keys ...string) Map {
	m := Map{}

	set := func(field string, value interface{}) {
		for _, key := range keys {
			if strings.EqualFold(key, field) {
				return
			}
		}

		m[field] = value
	}

	if x.code > 0 {
		set(keyCode, x.code)
	}

	if x.msg != "" {
		set(keyMessage, x.msg)
	}

	if len(x.details) > 0 {
		set(keyDetails, x.details)
	}

	if x.reason != nil {
		set(keyCause, x.reason.Error())
	}

	if x.stack != nil {
		set(keyStack, x.stack)
	}

	for k, v := range x.context {
		set(k, v)
	}

	return m
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
	switch verb {
	case 's':
		fallthrough
	case 'v':
		switch {
		case state.Flag('+'):
			errs.formatBullet(state, verb)
		case state.Flag('#'):
			fmt.Fprintf(state, "%#v", []error(errs))
		default:
			errs.formatSlice(state, verb)
		}
	}
}

func (errs ErrorCollector) formatBullet(state fmt.State, verb rune) {
	count := len(errs)

	for index, err := range errs {
		fmt.Fprint(state, " --- ")
		fmt.Fprintf(state, "%v", err)

		if index < count-1 {
			fmt.Fprint(state, "\n")
		}
	}
}

func (errs ErrorCollector) formatSlice(state fmt.State, verb rune) {
	fmt.Fprint(state, "[")

	for index, err := range errs {
		if index > 0 {
			fmt.Fprint(state, ", ")
		}

		fmt.Fprintf(state, "%v", err)
	}

	fmt.Fprint(state, "]")
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

// Wrap appends an error to the slice
func (errs *ErrorCollector) Wrap(err error) {
	*errs = append(*errs, err)
}

// Unwrap unwraps the underlying error it's only one
func (errs ErrorCollector) Unwrap() error {
	count := len(errs)

	switch {
	case count == 1:
		return errs[0]
	default:
		return nil
	}
}

// Code returns the code from an error
func Code(err error) int {
	type Coder interface {
		Code() int
	}

	if coder, ok := err.(Coder); ok {
		return coder.Code()
	}

	return 0
}

// Status returns the status from an error
func Status(err error) int {
	type Statuser interface {
		Status() int
	}

	if status, ok := err.(Statuser); ok {
		return status.Status()
	}

	return 0
}
