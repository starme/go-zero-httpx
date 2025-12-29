package errors

// DownloadError wraps an underlying download failure with an HTTP error code.
type DownloadError struct {
	code HttpCode
	Err  error
}

// NewDownloadError creates an HttpError that reports download failures.
func NewDownloadError(err error) HttpError {
	return &DownloadError{200, err}
}

// Code returns the HTTP status code that should be sent to the caller.
func (e DownloadError) Code() HttpCode {
	return e.code
}

// Error exposes the underlying download error message.
func (e DownloadError) Error() string {
	return e.Err.Error()
}
