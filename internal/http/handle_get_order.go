package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"uuid"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

func (h *Handler) handleGetorder(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	parsedId, err := uuid.Parse(id)
	if err != nil {
		respondWithError(w, r, domain.ErrValidation)
		return
	}
	order, err := h.orderStore.Get(r.Context(), domain.OrderID(parsedId))
	if err != nil {
		respondWithError(w, r, err)
		return
	}
	slog.Info("order retrieved", "order_id", order.ID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(order)
}
