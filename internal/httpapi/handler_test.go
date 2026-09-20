package httpapi_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/httpapi"
	"github.com/GuustTaillieu/idiomatic-go/internal/processor"
	"github.com/GuustTaillieu/idiomatic-go/internal/store"
)

func TestHandler_GetOrderRoute(t *testing.T) {
	p := processor.NewPaymentProcessor()
	s := store.NewOrderMemoryStore()
	q := domain.NewQueue(p, s)
	handler := httpapi.NewHandler(q, s)

	// Pre-populate the store with a order
	item := domain.NewItem("payload")
	original := domain.NewOrder(item.ID, 5)
	s.Save(nil, original)

	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/%s", original.ID), nil)
	rec := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rec.Code)
	}

	var result domain.Order
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	if result.ID != original.ID {
		t.Errorf("Expected order ID %s, got %s", original.ID, result.ID)
	}
	if result.ItemID != original.ItemID {
		t.Errorf("Expected order payload %s, got %s", original.ItemID, result.ItemID)
	}
}

func TestHandler_CreateOrderRoute(t *testing.T) {
	p := processor.NewPaymentProcessor()
	s := store.NewOrderMemoryStore()
	q := domain.NewQueue(p, s)
	handler := httpapi.NewHandler(q, s)

	fakeItem := domain.NewItem("payload")
	reqBody := fmt.Sprintf(`{"item_id": "%s"}`, fakeItem.ID)
	body := strings.NewReader(reqBody)
	req := httptest.NewRequest("POST", "/orders", body)
	rec := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rec, req)

	if rec.Code != http.StatusCreated {
		t.Errorf("Expected status code 201, got %d", rec.Code)
	}

	var result domain.Order
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	if result.ItemID != fakeItem.ID {
		t.Errorf("Expected order item ID %s, got %s", fakeItem.ID, result.ItemID)
	}
}
