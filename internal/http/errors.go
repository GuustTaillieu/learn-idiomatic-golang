package http

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

func respondWithError(w http.ResponseWriter, r *http.Request, err error) {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		http.Error(w, err.Error(), http.StatusNotFound)
	case errors.Is(err, domain.ErrConflict):
		http.Error(w, err.Error(), http.StatusConflict)
	case errors.Is(err, domain.ErrValidation):
		http.Error(w, err.Error(), http.StatusBadRequest)
	default:
		slog.ErrorContext(r.Context(), "unhandled server error", "error", err)
		http.Error(w, "internal Server Error", http.StatusInternalServerError)
	}
}
