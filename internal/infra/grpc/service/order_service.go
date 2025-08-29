package service

import (
	"context"

	"github.com/Berchon/Clean-Architecture/internal/infra/grpc/pb"
	"github.com/Berchon/Clean-Architecture/internal/usecase"
)

type OrderService struct {
	pb.UnimplementedOrderServiceServer
	OrderUseCase usecase.OrderUseCase
}

func NewOrderService(OrderUseCase usecase.OrderUseCase) *OrderService {
	return &OrderService{
		OrderUseCase: OrderUseCase,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, in *pb.CreateOrderRequest) (*pb.OrderResponse, error) {
	dto := usecase.OrderInputDTO{
		ID:    in.Id,
		Price: float64(in.Price),
		Tax:   float64(in.Tax),
	}
	output, err := s.OrderUseCase.Execute(dto)
	if err != nil {
		return nil, err
	}
	return &pb.OrderResponse{
		Id:         output.ID,
		Price:      float32(output.Price),
		Tax:        float32(output.Tax),
		FinalPrice: float32(output.FinalPrice),
	}, nil
}

func (s *OrderService) ListOrders(ctx context.Context, in *pb.Blank) (*pb.OrdersListResponse, error) {
	orders, err := s.OrderUseCase.GetOrders()
	if err != nil {
		return nil, err
	}

	var ordersResponse []*pb.OrderResponse
	for _, order := range orders {
		orderResponse := &pb.OrderResponse{
			Id:         order.ID,
			Price:      float32(order.Price),
			Tax:        float32(order.Tax),
			FinalPrice: float32(order.FinalPrice),
		}

		ordersResponse = append(ordersResponse, orderResponse)
	}

	return &pb.OrdersListResponse{Orders: ordersResponse}, nil
}
