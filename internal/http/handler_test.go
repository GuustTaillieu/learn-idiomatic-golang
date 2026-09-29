package http_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	httpapi "github.com/GuustTaillieu/idiomatic-go/internal/http"
	"github.com/GuustTaillieu/idiomatic-go/internal/memory"
)

func TestHandler_GetOrderRoute(t *testing.T) {
	s := memory.NewOrderStore()
	hc := &mockHealthChecker{}
	handler := httpapi.NewHandler(s, hc)

	// Pre-populate the store with a order
	item := domain.NewItem("payload")
	original := domain.NewOrder(item.ID, 5)
	s.Save(nil, original)

	req := httptest.NewRequest("GET", fmt.Sprintf("/orders/%s", original.ID), nil)
	rec := httptest.NewRecorder()

	handler.Routes().ServeHTTP(rec, req)

	// Assert
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
	// Arrange
	s := memory.NewOrderStore()
	hc := &mockHealthChecker{}
	handler := httpapi.NewHandler(s, hc)

	fakeItem := domain.NewItem("payload")
	reqBody := fmt.Sprintf(`{"item_id": "%s"}`, fakeItem.ID)
	body := strings.NewReader(reqBody)
	req := httptest.NewRequest("POST", "/orders", body)
	rec := httptest.NewRecorder()

	// Act
	handler.Routes().ServeHTTP(rec, req)

	// Assert
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
	if result.Status != domain.OrderStatusPending {
		t.Errorf("Expected order status %s, got %s", domain.OrderStatusPending, result.Status)
	}
}

func TestHandler_Healthz(t *testing.T) {
	// Arrange
	s := memory.NewOrderStore()
	hc := &mockHealthChecker{}
	handler := httpapi.NewHandler(s, hc)

	req := httptest.NewRequest("GET", "/healthz", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.Routes().ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rec.Code)
	}

	var result map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	if result["status"] != "alive" {
		t.Errorf("Expected status 'alive', got '%s'", result["status"])
	}
}

func TestHandler_Readyz(t *testing.T) {
	// Arrange
	s := memory.NewOrderStore()
	hc := &mockHealthChecker{}
	handler := httpapi.NewHandler(s, hc)

	req := httptest.NewRequest("GET", "/readyz", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.Routes().ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", rec.Code)
	}

	var result map[string]string
	if err := json.NewDecoder(rec.Body).Decode(&result); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	if result["status"] != "ready" {
		t.Errorf("Expected status 'ready', got '%s'", result["status"])
	}
}

func TestHandler_Readyz_Unhealthy(t *testing.T) {
	// Arrange
	s := memory.NewOrderStore()
	hc := &mockUnhealthyHealthChecker{}
	handler := httpapi.NewHandler(s, hc)

	req := httptest.NewRequest("GET", "/readyz", nil)
	rec := httptest.NewRecorder()

	// Act
	handler.Routes().ServeHTTP(rec, req)

	// Assert
	if rec.Code != http.StatusServiceUnavailable {
		t.Errorf("Expected status code 503, got %d", rec.Code)
	}
}

type mockUnhealthyHealthChecker struct{}

func (m *mockUnhealthyHealthChecker) Ping(ctx context.Context) error {
	return fmt.Errorf("service is unhealthy")
}

type mockHealthChecker struct{}

func (m *mockHealthChecker) Ping(ctx context.Context) error {
	return nil
}
