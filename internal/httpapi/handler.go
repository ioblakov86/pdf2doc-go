package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ioblakov86/pdf2doc-go/internal/convert"
	"github.com/ioblakov86/pdf2doc-go/internal/jobs"
)

func ConvertHandler_(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file required", 400)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.wordprocessingml.document")

	err = convert.Convert(ctx, file, w)
	if err != nil {
		if errors.Is(err, context.Canceled) {
			return
		}
		http.Error(w, err.Error(), 500)
	}
}

func ConvertHandler(m *jobs.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := m.Submit()
		json.NewEncoder(w).Encode(map[string]string{"job_id": id})
	}
}

func StatusHandler(m *jobs.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		job, ok := m.Get(id)
		if !ok {
			http.NotFound(w, r)
			return
		}
		json.NewEncoder(w).Encode(job)
	}
}
