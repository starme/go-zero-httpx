package response

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

func TestResolveDownloadPathRootNotConfigured(t *testing.T) {
	resetDownloadRoot()
	t.Cleanup(resetDownloadRoot)

	if _, err := resolveDownloadPath("file.txt"); err == nil {
		t.Fatalf("expected error without download root")
	}
}

func TestResolveDownloadPathEscapesRoot(t *testing.T) {
	root := t.TempDir()
	if err := SetDownloadRoot(root); err != nil {
		t.Fatalf("failed to set download root: %v", err)
	}
	t.Cleanup(resetDownloadRoot)

	if _, err := resolveDownloadPath("../etc/passwd"); err == nil {
		t.Fatalf("expected path escape to fail")
	}
}

func TestResolveDownloadPathInsideRoot(t *testing.T) {
	root := t.TempDir()
	if err := SetDownloadRoot(root); err != nil {
		t.Fatalf("failed to set download root: %v", err)
	}
	t.Cleanup(resetDownloadRoot)

	allowed := filepath.Join(root, "reports", "report.txt")
	if err := os.MkdirAll(filepath.Dir(allowed), 0o755); err != nil {
		t.Fatalf("create directories: %v", err)
	}
	if err := os.WriteFile(allowed, []byte("ok"), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	resolved, err := resolveDownloadPath(filepath.ToSlash(filepath.Join("reports", "report.txt")))
	if err != nil {
		t.Fatalf("resolve path: %v", err)
	}
	if resolved != allowed {
		t.Fatalf("expected %s, got %s", allowed, resolved)
	}
}

func TestDownloadCtxWritesFile(t *testing.T) {
	root := t.TempDir()
	if err := SetDownloadRoot(root); err != nil {
		t.Fatalf("failed to set download root: %v", err)
	}
	t.Cleanup(resetDownloadRoot)

	content := []byte("download content")
	filename := "report.txt"
	if err := os.WriteFile(filepath.Join(root, filename), content, 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}

	recorder := httptest.NewRecorder()
	DownloadCtx(context.Background(), recorder, filename, nil)

	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", recorder.Result().StatusCode)
	}
	expectedDisposition := fmt.Sprintf("attachment; filename=%s", url.QueryEscape(filename))
	if got := recorder.Header().Get("Content-Disposition"); got != expectedDisposition {
		t.Fatalf("unexpected Content-Disposition: %s", got)
	}
	if got := recorder.Header().Get("Content-Length"); got != fmt.Sprintf("%d", len(content)) {
		t.Fatalf("unexpected Content-Length: %s", got)
	}
	if !bytes.Equal(recorder.Body.Bytes(), content) {
		t.Fatalf("unexpected body: %s", recorder.Body.String())
	}
}

func TestDownloadCtxWithoutRootFails(t *testing.T) {
	resetDownloadRoot()
	t.Cleanup(resetDownloadRoot)

	recorder := httptest.NewRecorder()
	DownloadCtx(context.Background(), recorder, "file.txt", nil)
	if recorder.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("expected bad request when root is missing, got %d", recorder.Result().StatusCode)
	}
}
