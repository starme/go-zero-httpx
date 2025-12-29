package response

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"

	"github.com/starme/go-zero/httpx/errors"
)

var (
	downloadRoot   string
	downloadRootMu sync.RWMutex
)

// SetDownloadRoot configures the absolute root directory that download paths must not escape.
func SetDownloadRoot(root string) error {
	if root == "" {
		return fmt.Errorf("download root cannot be empty")
	}

	absRoot, err := filepath.Abs(root)
	if err != nil {
		return fmt.Errorf("resolve download root %s: %w", root, err)
	}

	downloadRootMu.Lock()
	downloadRoot = absRoot
	downloadRootMu.Unlock()
	return nil
}

func resetDownloadRoot() {
	downloadRootMu.Lock()
	downloadRoot = ""
	downloadRootMu.Unlock()
}

// Download streams the file at path to the client, reporting failures via HttpError.
func Download(w http.ResponseWriter, path string, err errors.HttpError) {
	DownloadCtx(context.Background(), w, path, err)
}

// DownloadCtx streams a file using the supplied context for tracing and error reporting.
func DownloadCtx(ctx context.Context, w http.ResponseWriter, path string, err error) {
	if err != nil {
		ErrorCtx(ctx, w, wrapDownloadErr(path, err))
		return
	}

	resolved, resolveErr := resolveDownloadPath(path)
	if resolveErr != nil {
		ErrorCtx(ctx, w, wrapDownloadErr(path, resolveErr))
		return
	}

	stat, err := os.Stat(resolved)
	if err != nil {
		ErrorCtx(ctx, w, wrapDownloadErr(path, err))
		return
	}
	if !stat.Mode().IsRegular() {
		ErrorCtx(ctx, w, wrapDownloadErr(path, fmt.Errorf("not a regular file")))
		return
	}

	file, err := os.Open(resolved)
	if err != nil {
		ErrorCtx(ctx, w, wrapDownloadErr(path, err))
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", url.QueryEscape(stat.Name())))
	w.Header().Set("Access-Control-Expose-Headers", "Content-Disposition")
	w.Header().Set("Content-Length", strconv.FormatInt(stat.Size(), 10))

	if _, err = io.Copy(w, file); err != nil {
		ErrorCtx(ctx, w, wrapDownloadErr(path, err))
		return
	}
}

// wrapDownloadErr converts a filesystem error into a DownloadError with contextualized messaging.
func wrapDownloadErr(path string, err error) errors.HttpError {
	return errors.NewDownloadError(
		fmt.Errorf("download %s: %w", path, err))
}

// resolveDownloadPath ensures the requested path resolves inside the configured download root.
func resolveDownloadPath(path string) (string, error) {
	downloadRootMu.RLock()
	root := downloadRoot
	downloadRootMu.RUnlock()
	if root == "" {
		return "", fmt.Errorf("download root is not configured")
	}

	cleaned := filepath.Clean(path)
	if cleaned == "." {
		return "", fmt.Errorf("download path cannot be empty")
	}

	var joined string
	if filepath.IsAbs(cleaned) {
		joined = cleaned
	} else {
		joined = filepath.Join(root, cleaned)
	}

	absPath, err := filepath.Abs(joined)
	if err != nil {
		return "", fmt.Errorf("resolve absolute path: %w", err)
	}

	rel, err := filepath.Rel(root, absPath)
	if err != nil {
		return "", fmt.Errorf("relate path to root: %w", err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("path escapes download root")
	}

	return absPath, nil
}
