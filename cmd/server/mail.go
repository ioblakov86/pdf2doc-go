package main

import (
	"log"
	"net/http"

	httpapi "github.com/ioblakov86/pdf2doc-go/internal/http"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/convert", httpapi.ConvertHandler)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
