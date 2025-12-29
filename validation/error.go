package validation

import (
	"bytes"
	"errors"
	"strings"

	xerr "github.com/starme/go-zero/httpx/errors"
)

// ValidateError aggregates validation failures for later reporting.
type ValidateError []error

// Code returns a consistent HTTP code for validation failures.
func (e ValidateError) Code() xerr.HttpCode {
	return 100
}

// Error builds a newline-separated message containing all validation error strings.
func (e ValidateError) Error() string {
	buff := bytes.NewBufferString("")

	for i := 0; i < len(e); i++ {

		buff.WriteString(e[i].Error())
		buff.WriteString("\n")
	}

	return strings.TrimSpace(buff.String())
}

// AddString appends a new validation error message.
func (e ValidateError) AddString(msg string) ValidateError {
	return append(e, errors.New(msg))
}
