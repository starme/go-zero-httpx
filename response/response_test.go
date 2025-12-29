package response

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestSuccessWritesEmptyData(t *testing.T) {
	recorder := httptest.NewRecorder()

	Success(recorder)

	if recorder.Result().StatusCode != http.StatusOK {
		t.Fatalf("unexpected status: %d", recorder.Result().StatusCode)
	}

	var body Body
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}

	dataSlice, ok := body.Data.([]any)
	if !ok {
		t.Fatalf("data should be []any, got %T", body.Data)
	}
	if len(dataSlice) != 0 {
		t.Fatalf("expected empty data slice, got %v", dataSlice)
	}
}

func TestSuccessWithData(t *testing.T) {
	recorder := httptest.NewRecorder()
	payload := map[string]any{"key": "value"}

	Success(recorder, payload)

	var body Body
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}

	dataSlice, ok := body.Data.([]any)
	if !ok {
		t.Fatalf("data should be []any, got %T", body.Data)
	}
	if len(dataSlice) != 1 {
		t.Fatalf("expected one element, got %v", dataSlice)
	}
	if !reflect.DeepEqual(dataSlice[0], payload) {
		t.Fatalf("expected payload %v, got %v", payload, dataSlice[0])
	}
}

func TestResponseUsesProvidedData(t *testing.T) {
	recorder := httptest.NewRecorder()
	data := []any{1, 2, 3}

	Response(recorder, http.StatusBadRequest, 123, data, nil)

	if recorder.Result().StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status: %d", recorder.Result().StatusCode)
	}

	var body Body
	if err := json.NewDecoder(recorder.Body).Decode(&body); err != nil {
		t.Fatalf("decode body: %v", err)
	}

	dataSlice, ok := body.Data.([]any)
	if !ok {
		t.Fatalf("data should be []any, got %T", body.Data)
	}
	if len(dataSlice) != len(data) {
		t.Fatalf("expected %d elements, got %d", len(data), len(dataSlice))
	}
	for i, got := range dataSlice {
		want := float64(i + 1)
		if got != want {
			t.Fatalf("element %d mismatch: want %v, got %v", i, want, got)
		}
	}
	if body.Code != 123 {
		t.Fatalf("unexpected code: %d", body.Code)
	}
}
