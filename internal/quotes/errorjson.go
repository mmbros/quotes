package quotes

import "encoding/json"

// ErrorJsonizable is an error that can be JSON-serialized as the Error() string.
type ErrorJsonizable struct {
	err error
}

func (e *ErrorJsonizable) Error() string {
	if e == nil || e.err == nil {
		return ""
	}
	return e.err.Error()
}

func (e *ErrorJsonizable) Unwrap() error {
	if e == nil {
		return nil
	}
	return e.err
}

func (e *ErrorJsonizable) MarshalJSON() ([]byte, error) {
	if e == nil || e.err == nil {
		return nil, nil
	}
	// eerr := e.err
	// if errors.Is(eerr, context.Canceled) {
	// 	eerr = context.Canceled
	// }
	// return json.Marshal(eerr.Error())
	return json.Marshal(e.err.Error())
}
