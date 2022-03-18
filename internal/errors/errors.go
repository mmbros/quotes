package errors

import (
	"fmt"
)

// wrappedError is an error object with message and underlying error.
type wrappedError struct {
	msg string
	err error
}

func (e *wrappedError) Unwrap() error { return e.err }
func (e *wrappedError) Error() string { return e.msg }

// // wrapNameError returns an error object with inner error and "name: err" message.
// // If err is nil, wrapNameError returns nil.
// func wrapNameError(err error, name string) error {
// 	if err == nil {
// 		return nil
// 	}
// 	return WrapErrorf(err, "%s: %s", name, err.Error())
// }

// // wrapNameErrorString returns an error object with inner error and "name: err str" message.
// // If err is nil, wrapNameErrorString returns nil.
// func wrapNameErrorString(err error, name, str string) error {
// 	if err == nil {
// 		return nil
// 	}
// 	return WrapErrorf(err, "%s: %s %q", name, err.Error(), str)
// }

// WrapErrorf returns an error object with inner error and formatted message.
// If err is nil, WrapErrorf returns nil.
func WrapErrorf(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}
	return &wrappedError{fmt.Sprintf(format, a...), err}
}
