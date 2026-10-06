package http_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"connectrpc.com/connect"
	httpapi "github.com/GuustTaillieu/idiomatic-go/internal/http"
	"github.com/GuustTaillieu/idiomatic-go/internal/memory"
	v1 "github.com/GuustTaillieu/idiomatic-go/proto/order/v1"
	"github.com/GuustTaillieu/idiomatic-go/proto/order/v1/orderv1connect"
)

func TestOrderRPCServer_CreateAndGetOrder(t *testing.T) {
	store := memory.NewOrderStore()
	mux := http.NewServeMux()
	path, handler := orderv1connect.NewOrderServiceHandler(httpapi.NewOrderRPCServer(store))
	mux.Handle(path, handler)

	server := httptest.NewServer(mux)
	defer server.Close()

	// Create type-safe client pointing to test server
	client := orderv1connect.NewOrderServiceClient(server.Client(), server.URL)

	// First, create an order to retrieve later
	createReq := connect.NewRequest(&v1.CreateOrderRequest{
		ItemId: "00000000-0000-0000-0000-000000000002",
		Amount: 3,
	})
	createRes, err := client.CreateOrder(context.Background(), createReq)
	if err != nil {
		t.Fatalf("CreateOrder failed: %v", err)
	}

	// Now, retrieve the created order
	getReq := connect.NewRequest(&v1.GetOrderRequest{
		OrderId: createRes.Msg.OrderId,
	})
	getRes, err := client.GetOrder(context.Background(), getReq)
	if err != nil {
		t.Fatalf("GetOrder failed: %v", err)
	}

	if getRes.Msg.OrderId != createRes.Msg.OrderId {
		t.Fatalf("Expected OrderId %s, got %s", createRes.Msg.OrderId, getRes.Msg.OrderId)
	}
	if getRes.Msg.ItemId != "00000000-0000-0000-0000-000000000002" {
		t.Fatalf("Expected ItemId 00000000-0000-0000-0000-000000000002, got %s", getRes.Msg.ItemId)
	}
	if getRes.Msg.Amount != 3 {
		t.Fatalf("Expected Amount 3, got %d", getRes.Msg.Amount)
	}
}

func TestOrderRPCServer_GetOrderWithInvalidUUID(t *testing.T) {
	store := memory.NewOrderStore()
	mux := http.NewServeMux()
	path, handler := orderv1connect.NewOrderServiceHandler(httpapi.NewOrderRPCServer(store))
	mux.Handle(path, handler)

	server := httptest.NewServer(mux)
	defer server.Close()

	// Create type-safe client pointing to test server
	client := orderv1connect.NewOrderServiceClient(server.Client(), server.URL)

	// Attempt to retrieve an order with an invalid UUID
	getReq := connect.NewRequest(&v1.GetOrderRequest{
		OrderId: "invalid-uuid",
	})
	_, err := client.GetOrder(context.Background(), getReq)
	if err == nil {
		t.Fatalf("Expected error for invalid UUID, got nil")
	}

	if connect.CodeOf(err) != connect.CodeInvalidArgument {
		t.Fatalf("Expected InvalidArgument error code, got %v", connect.CodeOf(err))
	}
}
