package http

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/LCGant/role-audit/internal/config"
	"github.com/LCGant/role-audit/internal/store"
)

func TestInternalEventsRejectTrailingJSONData(t *testing.T) {
	cfg := config.Config{
		InternalToken: "secret",
		MetricsToken:  "metrics",
		LogFile:       t.TempDir() + "/audit.jsonl",
	}
	fileStore, err := store.NewFileStore(cfg.LogFile)
	if err != nil {
		t.Fatalf("new file store: %v", err)
	}
	h := New(cfg, slog.New(slog.NewTextHandler(io.Discard, nil)), fileStore)

	req := httptest.NewRequest(http.MethodPost, "/internal/events", bytes.NewBufferString(`{"source":"auth","event_type":"login","success":true}junk`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Internal-Token", "secret")
	rr := httptest.NewRecorder()

	h.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rr.Code)
	}
}
