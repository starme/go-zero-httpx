package request

import (
	"net/http"

	"github.com/starme/go-zero/httpx/validation"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// Parse decodes the incoming request into v and validates the resulting struct.
func Parse(r *http.Request, v any) error {
	if err := httpx.Parse(r, v); err != nil {
		return err
	}

	return validation.Validate(r.Context(), v)
}

// ParseForm reads form values from the request body or query string into v and validates it.
func ParseForm(r *http.Request, v any) error {
	if err := httpx.ParseForm(r, v); err != nil {
		return err
	}

	return validation.Validate(r.Context(), v)
}

// ParseJsonBody decodes a JSON payload from the request body into v and validates it.
func ParseJsonBody(r *http.Request, v any) error {
	if err := httpx.ParseJsonBody(r, v); err != nil {
		return err
	}

	return validation.Validate(r.Context(), v)
}

// ParsePath binds URI path parameters into v and validates the result.
// For example: http://localhost/bag/:name.
func ParsePath(r *http.Request, v any) error {
	if err := httpx.ParsePath(r, v); err != nil {
		return err
	}

	return validation.Validate(r.Context(), v)
}
