package main

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/ioblakov86/pdf2doc-go/internal/config"
	"github.com/ioblakov86/pdf2doc-go/internal/httpapi"
	"github.com/ioblakov86/pdf2doc-go/internal/jobs"
)

var manager *jobs.Manager

func main() {
	cfg := config.Load()

	manager = jobs.NewManager(cfg.MaxConcurrent)

	r := chi.NewRouter()
	r.Post("/convert", httpapi.ConvertHandler(manager))
	r.Get("/jobs/{id}", httpapi.StatusHandler(manager))

	log.Printf("Listening on :%s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}
