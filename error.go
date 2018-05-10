package errx

import (
	"bytes"
	"fmt"
	"strings"
)

// Error wraps an error and has a message and stack trace associated with it.
type Error struct {
	Inner      error       `json:"inner,omitempty"`
	Message    string      `json:"msg,omitempty"`
	StackTrace *StackTrace `json:"stack,omitempty"`
}

// Error returns a formatted error string, including all inner errors.
func (e *Error) Error() string {
	return e.error(0)
}

func (e *Error) isZero() bool {
	return e.Inner == nil && e.Message == "" && e.StackTrace == nil
}

func (e *Error) error(depth int) string {
	b := new(bytes.Buffer)
	if e.Message != "" {
		pad(b, padding)
		b.WriteString(e.Message)
	}

	if e.Inner != nil {
		if inner, ok := e.Inner.(*Error); ok {
			separator := fmt.Sprintf("\n%s", strings.Repeat(separator, depth+1))
			if !inner.isZero() {
				pad(b, separator)
				b.WriteString(inner.error(depth + 1))
			}
		} else {
			pad(b, padding)
			b.WriteString(e.Inner.Error())
		}
	}

	if e.StackTrace != nil {
		b.WriteString("\n")
		if depth > 0 {
			prefix := strings.Repeat(separator, depth+1)
			b.WriteString(strings.Replace(e.StackTrace.String(), prefix[depth*2:], prefix, -1))
		} else {
			b.WriteString(e.StackTrace.String())
		}
	}

	return b.String()
}

func pad(b *bytes.Buffer, str string) {
	if b.Len() == 0 {
		return
	}
	b.WriteString(str)
}

const (
	padding   = " "
	separator = "  "
)

// New creates a new error with a stack trace at the point which New was called,
// a message, and a nil inner error.
func New(message string) error {
	return newErr(message)
}

// Errorf creates a new error with a stack trace at the point which Errorf was called,
// a formatted message, and a nil inner error.
func Errorf(format string, args ...interface{}) error {
	return newErr(fmt.Sprintf(format, args...))
}

// Wrap wraps an existing error with a message. If the inner error is an errx.Error, then
// no stack trace is added, otherwise a stack trace is captured at the point which Wrap
// was called.
func Wrap(err error, message string) error {
	return wrapErr(err, message)
}

// Wrapf wraps an existing error with a formatted message. If the inner error is an
// errx.Error, then no stack trace is added, otherwise a stack trace is captured at
// the point which Wrapf was called.
func Wrapf(err error, format string, args ...interface{}) error {
	return wrapErr(err, fmt.Sprintf(format, args...))
}

func newErr(message string) error {
	return &Error{
		Inner:      nil,
		Message:    message,
		StackTrace: getStack(),
	}
}

func wrapErr(err error, message string) error {
	e := &Error{Message: message}

	if inner, ok := err.(*Error); ok {
		copy := *inner
		e.Inner = &copy
		return e
	}

	e.Inner = err
	e.StackTrace = getStack()
	return e
}

// GetInner returns the inner error from a message if it has one or nil otherwise.
func GetInner(err error) error {
	if e, ok := err.(*Error); ok {
		return e.Inner
	}

	return nil
}

// GetStackTrace returns the StackTrace from an error message if it has one or nil otherwise.
func GetStackTrace(err error) *StackTrace {
	if e, ok := err.(*Error); ok {
		return e.StackTrace
	}

	return nil
}
