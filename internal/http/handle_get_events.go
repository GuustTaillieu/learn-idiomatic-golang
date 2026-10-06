package http

import (
	"fmt"
	"log/slog"
	"net/http"
)

func (h *Handler) handleGetEvents(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	subscriber := h.orderHub.Subscribe()
	defer h.orderHub.Unsubscribe(subscriber)

	for {
		select {
		case <-r.Context().Done(): // Client closed the tab/connection
			return
		case event := <-subscriber:
			// Write to io.Writer
			jsonData, err := event.MarshalJSON()
			if err != nil {
				slog.Error("Failed to marshal event", "error", err)
				http.Error(w, "Something went wrong with sending the event", http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()
		}
	}
}
