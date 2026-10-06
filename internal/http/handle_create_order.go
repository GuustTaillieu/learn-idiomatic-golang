package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/telemetry"
)

type CreateorderRequest struct {
	ItemID domain.ItemID `json:"item_id"`
	Amount int           `json:"amount"`
}

func (h *Handler) handleCreateorder(w http.ResponseWriter, r *http.Request) {
	ctx, span := telemetry.StartSpan(r.Context(), "http.handleCreateorder")
	defer span.End()

	var req CreateorderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, r, domain.ErrValidation)
		return
	}
	order := domain.NewOrder(req.ItemID, req.Amount, domain.WithTraceParent(telemetry.InjectTraceParent(ctx)))
	if err := order.Validate(); err != nil {
		respondWithError(w, r, err)
		return
	}
	if err := h.orderStore.Save(ctx, order); err != nil {
		respondWithError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)

	slog.Info("order submitted", "order_id", order.ID)
}
