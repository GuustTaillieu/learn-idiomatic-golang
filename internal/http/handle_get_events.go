package http

import (
	"fmt"
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
				http.Error(w, "Failed to marshal event", http.StatusInternalServerError)
				return
			}
			fmt.Fprintf(w, "data: %s\n\n", jsonData)
			flusher.Flush()
		}
	}
}
