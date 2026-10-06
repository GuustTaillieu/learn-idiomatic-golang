package http

import (
	"context"
	"uuid"

	"connectrpc.com/connect"

	"github.com/GuustTaillieu/idiomatic-go/internal/domain"
	"github.com/GuustTaillieu/idiomatic-go/internal/telemetry"
	v1 "github.com/GuustTaillieu/idiomatic-go/proto/order/v1"
)

type OrderRPCServer struct {
	orderStore domain.OrderStorer
}

func NewOrderRPCServer(orderStore domain.OrderStorer) *OrderRPCServer {
	return &OrderRPCServer{
		orderStore: orderStore,
	}
}

func (s *OrderRPCServer) CreateOrder(ctx context.Context, req *connect.Request[v1.CreateOrderRequest]) (*connect.Response[v1.CreateOrderResponse], error) {
	ctx, span := telemetry.StartSpan(ctx, "rpc.CreateOrder")
	defer span.End()

	id, err := uuid.Parse(req.Msg.ItemId)
	if err != nil {
		span.RecordError(err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	itemID := domain.ItemID(id)

	order := domain.NewOrder(itemID, int(req.Msg.Amount), domain.WithTraceParent(telemetry.InjectTraceParent(ctx)))
	if err := s.orderStore.Save(ctx, order); err != nil {
		span.RecordError(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&v1.CreateOrderResponse{
		OrderId: order.ID.String(),
		Status:  string(order.Status),
	}), nil
}

func (s *OrderRPCServer) GetOrder(ctx context.Context, req *connect.Request[v1.GetOrderRequest]) (*connect.Response[v1.GetOrderResponse], error) {
	ctx, span := telemetry.StartSpan(ctx, "rpc.GetOrder")
	defer span.End()

	id, err := uuid.Parse(req.Msg.OrderId)
	if err != nil {
		span.RecordError(err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}
	orderID := domain.OrderID(id)

	order, err := s.orderStore.Get(ctx, orderID)
	if err != nil {
		span.RecordError(err)
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	return connect.NewResponse(&v1.GetOrderResponse{
		OrderId: order.ID.String(),
		ItemId:  order.ItemID.String(),
		Amount:  int32(order.Amount),
		Status:  string(order.Status),
	}), nil
}
