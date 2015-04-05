package er

import (
	"encoding/json"

	"runtime"
)

type Error struct {
	Message  string      `json:"message,omitempty"`
	Location CallContext `json:"location,omitempty"`
	Details  *KV         `json:"details,omitempty"`
	Nested   *Error      `json:"nested,omitempty"`
	External error       `json:"external,omitempty"` // inner or source error
	Stack    string      `json:"stack,omitempty"`
}

func (e *Error) Error() string {
	if bytes, err := json.MarshalIndent(e, "", "\t"); err != nil {
		return Message + " Error Failed to Marshal: " + err.Error()
	}
	return string(bytes)
}

func stackBuffer() []byte {
	const size = 1 << 12
	return make([]byte, size)
}

func NewErr(message string, details ...KV) *Error {
	buf := stackBuffer()
	n := runtime.Stack(buf, false)
	stackTrace := string(buf[:n])

	return &Error{
		Message:  message,
		Location: CallContext(1),
		Details:  details,
		Nested:   nested,
		External: external,
		Stack:    stackTrace,
	}
}

func Err(err error, message string, details ...KV) *Error {
	stackTrace := ""
	var external error = nil

	nested, ok := err.(*Error)
	if !ok {
		external = err

		//stack trace
		buf := stackBuffer()
		n := runtime.Stack(buf, false)
		stackTrace = string(buf[:n])
	}

	return &Error{
		Message:  message,
		Location: CallContext(1),
		Details:  details,
		Nested:   nested,
		External: external,
		Stack:    stackTrace,
	}
}

type ApiError struct {
	Error
	StatusCode    int    `json:"statusCode"`
	StatusMessage string `json:"statusMessage,omitempty"`
}

func NewApiErr(message string, err error, statusCode int, statusMessage string) *ApiError {
	stackTrace := ""
	var external error = nil

	nested, ok := err.(*Error)
	if !ok {
		external = err

		//stack trace
		buf := stackBuffer()
		n := runtime.Stack(buf, false)
		stackTrace = string(buf[:n])
	}

	return &ApiError{
		Message:       message,
		Location:      CallContext(1),
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
		Details:       details,
		Nested:        nested,
		External:      external,
		Stack:         stackTrace,
	}

}
