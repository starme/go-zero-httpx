package response

import (
	"context"
	"net/http"

	"github.com/starme/go-zero/httpx/errors"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Body defines the standard envelope returned for HTTP responses.
type Body struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// Success writes a default HTTP 200 response with the provided payload.
func Success(w http.ResponseWriter, data ...any) {
	SuccessCtx(context.Background(), w, data...)
}

// SuccessCtx writes a success response using the provided context for tracing/logging.
func SuccessCtx(ctx context.Context, w http.ResponseWriter, data ...any) {
	responseCtx(ctx, w, http.StatusOK, 0, data, nil)
}

// Error writes an HttpError payload using the default context.
func Error(w http.ResponseWriter, err errors.HttpError) {
	ErrorCtx(context.Background(), w, err)
}

// ErrorCtx writes an HttpError payload while carrying the supplied context.
func ErrorCtx(ctx context.Context, w http.ResponseWriter, err errors.HttpError) {
	responseCtx(ctx, w, http.StatusBadRequest, 0, nil, err)
}

// Response writes a response with explicit status, code, optional data, and error payloads.
func Response(w http.ResponseWriter, status, code int, data any, err errors.HttpError) {
	ResponseCtx(context.Background(), w, status, code, data, err)
}

// ResponseCtx writes a response using the provided context.
func ResponseCtx(ctx context.Context, w http.ResponseWriter, status, code int, data any, err errors.HttpError) {
	responseCtx(ctx, w, status, code, data, err)
}

func responseCtx(ctx context.Context, w http.ResponseWriter, status, code int, data any, err errors.HttpError) {
	httpx.WriteJsonCtx(ctx, w, status, wrapResponse(code, data, err))
}

func wrapResponse(code int, data any, err errors.HttpError) Body {
	body := Body{Code: code, Data: formatData(data), Msg: "success"}

	if err != nil {
		body.Code = int(err.Code())
		body.Msg = err.Error()
	}

	return body
}

func formatData(data any) any {
	if data == nil {
		return []any{}
	}

	if slice, ok := data.([]any); ok && slice == nil {
		return []any{}
	}

	return data
}
