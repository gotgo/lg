package er

import (
	"encoding/json"

	"github.com/gotgo/lg"

	"runtime"
)

type Error struct {
	Message  string          `json:"message,omitempty"`
	Location *lg.CallContext `json:"location,omitempty"`
	Details  lg.KV           `json:"details,omitempty"`
	Nested   *Error          `json:"nested,omitempty"`
	External error           `json:"external,omitempty"` // inner or source error
	Stack    string          `json:"stack,omitempty"`
}

func (e *Error) Error() string {
	if bytes, err := json.MarshalIndent(e, "", "\t"); err != nil {
		return e.Message + " Error Failed to Marshal: " + err.Error()
	} else {
		return string(bytes)
	}
}

func stackBuffer() []byte {
	const size = 1 << 12
	return make([]byte, size)
}

func NewErr(message string, details ...lg.KV) *Error {
	buf := stackBuffer()
	n := runtime.Stack(buf, false)
	stackTrace := string(buf[:n])
	ctx, _ := lg.CallerContext(1)
	return &Error{
		Message:  message,
		Location: ctx,
		Details:  lg.CollapseKV(details),
		Nested:   nil,
		External: nil,
		Stack:    stackTrace,
	}
}

func Err(err error, message string, details ...lg.KV) *Error {
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
	ctx, _ := lg.CallerContext(1)

	return &Error{
		Message:  message,
		Location: ctx,
		Details:  lg.CollapseKV(details),
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

func NewApiErr(message string, err error, statusCode int, statusMessage string, details ...lg.KV) *ApiError {
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
	ctx, _ := lg.CallerContext(1)

	r := &ApiError{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
	}
	r.Message = message
	r.Location = ctx
	r.Details = lg.CollapseKV(details)
	r.Nested = nested
	r.External = external
	r.Stack = stackTrace
	return r
}
