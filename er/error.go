package er

import (
	"encoding/json"

	"runtime"
)

type KV map[string]interface{}

type Error struct {
	Message  string       `json:"message,omitempty"`
	Location *CallContext `json:"location,omitempty"`
	Details  KV           `json:"details,omitempty"`
	Nested   *Error       `json:"nested,omitempty"`
	External error        `json:"external,omitempty"` // inner or source error
	Stack    string       `json:"stack,omitempty"`
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

func NewErr(message string, details ...KV) *Error {
	buf := stackBuffer()
	n := runtime.Stack(buf, false)
	stackTrace := string(buf[:n])
	ctx, _ := CallerContext(1)
	return &Error{
		Message:  message,
		Location: ctx,
		Details:  collapse(details),
		Nested:   nil,
		External: nil,
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
	ctx, _ := CallerContext(1)

	return &Error{
		Message:  message,
		Location: ctx,
		Details:  collapse(details),
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

func NewApiErr(message string, err error, statusCode int, statusMessage string, details ...KV) *ApiError {
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
	ctx, _ := CallerContext(1)

	r := &ApiError{
		StatusCode:    statusCode,
		StatusMessage: statusMessage,
	}
	r.Message = message
	r.Location = ctx
	r.Details = collapse(details)
	r.Nested = nested
	r.External = external
	r.Stack = stackTrace
	return r
}

func collapse(kvs []KV, values ...KV) KV {
	if kvs == nil {
		kvs = append([]KV{}, values...)
	}
	if len(kvs) > 1 {
		data := KV{}
		for _, mp := range kvs {
			if mp == nil {
				continue
			}
			for k, v := range mp {
				data[k] = v
			}
		}
		return data
	} else if len(kvs) == 1 {
		return kvs[0]
	} else {
		return nil
	}
}
