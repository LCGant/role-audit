package http

import (
	"bytes"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"expvar"
	"io"
	"log/slog"
	nethttp "net/http"
	"strings"
	"time"

	"github.com/LCGant/role-audit/internal/config"
	"github.com/LCGant/role-audit/internal/store"
)

const maxBodyBytes = 1 << 20

type Server struct {
	cfg    config.Config
	logger *slog.Logger
	store  *store.FileStore
	mux    *nethttp.ServeMux
}

func New(cfg config.Config, logger *slog.Logger, fileStore *store.FileStore) nethttp.Handler {
	s := &Server{
		cfg:    cfg,
		logger: logger,
		store:  fileStore,
		mux:    nethttp.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("GET /healthz", func(w nethttp.ResponseWriter, r *nethttp.Request) {
		writeJSON(w, nethttp.StatusOK, map[string]string{"status": "ok"})
	})
	s.mux.Handle("/metrics", s.metricsGuard(expvar.Handler()))
	s.mux.Handle("POST /internal/events", s.internal(s.handleEvent))
}

func (s *Server) ServeHTTP(w nethttp.ResponseWriter, r *nethttp.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) internal(next nethttp.HandlerFunc) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if subtle.ConstantTimeCompare([]byte(strings.TrimSpace(r.Header.Get("X-Internal-Token"))), []byte(s.cfg.InternalToken)) != 1 {
			writeError(w, nethttp.StatusUnauthorized, "unauthorized")
			return
		}
		next(w, r)
	})
}

func (s *Server) metricsGuard(next nethttp.Handler) nethttp.Handler {
	return nethttp.HandlerFunc(func(w nethttp.ResponseWriter, r *nethttp.Request) {
		if s.cfg.MetricsToken == "" {
			writeError(w, nethttp.StatusForbidden, "forbidden")
			return
		}
		if subtle.ConstantTimeCompare([]byte(strings.TrimSpace(r.Header.Get("X-Metrics-Token"))), []byte(s.cfg.MetricsToken)) != 1 {
			writeError(w, nethttp.StatusForbidden, "forbidden")
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) handleEvent(w nethttp.ResponseWriter, r *nethttp.Request) {
	var event store.Event
	if err := decodeJSON(r, &event); err != nil {
		writeError(w, nethttp.StatusBadRequest, "bad_request")
		return
	}
	event.Source = strings.TrimSpace(event.Source)
	event.EventType = strings.TrimSpace(event.EventType)
	if event.Source == "" || event.EventType == "" {
		writeError(w, nethttp.StatusBadRequest, "bad_request")
		return
	}
	if event.CreatedAt.IsZero() {
		event.CreatedAt = time.Now().UTC()
	}
	event.ReceivedAt = time.Now().UTC()
	if err := s.store.Append(r.Context(), event); err != nil {
		s.logger.Error("append audit failed", "err", err)
		writeError(w, nethttp.StatusBadGateway, "store_failed")
		return
	}
	writeJSON(w, nethttp.StatusAccepted, map[string]string{"status": "accepted"})
}

func decodeJSON(r *nethttp.Request, dst any) error {
	payload, err := io.ReadAll(io.LimitReader(r.Body, maxBodyBytes+1))
	if err != nil {
		return err
	}
	if int64(len(payload)) > maxBodyBytes {
		return errors.New("too_large")
	}
	dec := json.NewDecoder(bytes.NewReader(payload))
	dec.DisallowUnknownFields()
	if err := dec.Decode(dst); err != nil {
		return err
	}
	if err := dec.Decode(new(struct{})); err != io.EOF {
		return errors.New("trailing_data")
	}
	return nil
}

func writeJSON(w nethttp.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

func writeError(w nethttp.ResponseWriter, status int, code string) {
	writeJSON(w, status, map[string]string{"error": code})
}
