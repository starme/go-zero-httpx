package errors

// HttpCode represents the HTTP status code associated with an HttpError.
type HttpCode int

// HttpError augments the standard error interface with an HTTP response code.
type HttpError interface {
	error
	Code() HttpCode
}
