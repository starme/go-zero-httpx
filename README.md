# go-zero-httpx

Lightweight helpers built on top of `github.com/zeromicro/go-zero/rest/httpx` that
standardize request parsing, response formatting, validation, and safe file downloads.
The package keeps handlers simple by coupling decode-and-validate flows with reusable
response envelopes.

## Table of Contents
1. [Getting Started](#getting-started)
2. [Configuration](#configuration)
3. [Project Layout](#project-layout)
4. [Usage](#usage)
5. [Testing](#testing)
6. [Troubleshooting](#troubleshooting)
7. [Contributing](#contributing)
8. [License](#license)

## Getting Started

### Prerequisites
- Go 1.23.3 or later
- Clone this module: `git clone https://github.com/starme/go-zero/httpx`

### Installation

```shell
cd go-zero-httpx
go mod tidy
```

The module exposes reusable builders for request decoding, validation, success/error
responses, and authorized downloads, so it is intended to be consumed by Go services
that already depend on `go-zero/rest/httpx`.

## Configuration

- **Validation translations**: register additional locales via
  `validation.RegisterLocaleTranslation` and switch languages with `validation.NewTranslator`.
- **Download root**: call `response.SetDownloadRoot("/path/to/allowed/files")` before
  `response.Download` to prevent directory traversal. Requests outside the root return
  a wrapped `errors.DownloadError`.

## Project Layout

- `errors/` — HTTP error definitions such as `DownloadError` and shared `HttpError`.
- `request/` — Wrapper helpers (`Parse`, `ParseForm`, `ParseJsonBody`, `ParsePath`)
  that decode and validate incoming HTTP payloads in one step.
- `response/` — Unified `Body`, success/error helpers, JSON writers, and download logic.
- `validation/` — Validator wrapper exposing translation registration, custom
  validators, and the `ValidateError` aggregate.

## Usage

### Decode and validate requests

```go
type loginRequest struct {
    Account string `json:"account" validate:"required"`
    Password string `json:"password" validate:"required"`
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
    var req loginRequest
    if err := request.Parse(r, &req); err != nil {
        response.Error(w, err.(errors.HttpError))
        return
    }
    response.Success(w, map[string]any{"token": "secret"})
}
```

### Success and error responses

Use `response.Success(ctx, w, payload...)` for a 200-level envelope and
`response.Error(ctx, w, err)` for consistent error shaping. Customize status,
code, or body with `response.Response`.

### Controlled file downloads

```go
response.SetDownloadRoot("/var/downloads")
response.Download(w, "/reports/november.pdf", nil)
```

The helper sanitizes paths, rejects directories, sets content headers, and reports
`errors.DownloadError` instances for any filesystem issues.

## Testing

```bash
go test ./...
```

Unit tests cover response writing, download guards, and validation helpers.

## Troubleshooting

- If validation messages still appear in raw English, ensure a translator
  (`validation.NewTranslator`) has been registered for the desired locale.
- `response.Download` requires a configured root; otherwise, it returns a download error.

## Contributing

Feel free to open issues or pull requests against this repository. Follow Go
conventions, add tests for new behavior, and keep changes focused and reviewable.

## License

See the `LICENSE` file in the repository root for licensing details.

