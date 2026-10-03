package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
)

type CreateorderRequest struct {
	ItemID domain.ItemID `json:"item_id"`
	Amount int           `json:"amount"`
}

func (h *Handler) handleCreateorder(w http.ResponseWriter, r *http.Request) {
	var req CreateorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}
	order := domain.NewOrder(req.ItemID, req.Amount)
	if err := h.orderStore.Save(r.Context(), order); err != nil {
		http.Error(w, "Failed to save order", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)

	slog.Info("order submitted", "order_id", order.ID)
}
