package svc

import (
	"context"
	"encoding/json"
	svc_http "github.com/boram-gong/service/svc/http"
	"net/http"
)

var (
	ErrInvalidPara    = NewError(http.StatusBadRequest, "invalid para")
	ErrNotFound       = NewError(http.StatusNotFound, "not found")
	ErrInternalServer = NewError(http.StatusInternalServerError, "internal main-server error")
)

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) StatusCode() int {
	return e.Code
}

func NewError(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

type errorWrapper struct {
	Error string `json:"error"`
}

func ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	body, _ := json.Marshal(errorWrapper{Error: err.Error()})
	if marshal, ok := err.(json.Marshaler); ok {
		if jsonBody, marshalErr := marshal.MarshalJSON(); marshalErr == nil {
			body = jsonBody
		}
	}
	w.Header().Set("Content-Type", ContentType)
	if head, ok := err.(svc_http.Headerer); ok {
		for k := range head.Headers() {
			w.Header().Set(k, head.Headers().Get(k))
		}
	}
	code := http.StatusInternalServerError
	if sc, ok := err.(svc_http.StatusCoder); ok {
		code = sc.StatusCode()
	}
	w.WriteHeader(code)
	_, _ = w.Write(body)
}
