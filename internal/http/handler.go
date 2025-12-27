package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/ioblakov86/pdf2doc-go/internal/convert"
)

func ConvertHandler(w http.ResponseWriter, r *http.Request) {
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
