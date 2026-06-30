package grpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/joanaeliseal/microservices-proto/golang/shipping"
	"github.com/joanaeliseal/microservices/shipping/config"
	"github.com/joanaeliseal/microservices/shipping/internal/application/core/domain"
	"github.com/joanaeliseal/microservices/shipping/internal/ports"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	api  ports.APIPort
	port int
	shipping.UnimplementedShippingServer
}

func NewAdapter(api ports.APIPort, port int) *Adapter {
	return &Adapter{api: api, port: port}
}

func (a Adapter) Ship(ctx context.Context, request *shipping.ShipOrderRequest) (*shipping.ShipOrderResponse, error) {
	var items []domain.ShippingItem
	for _, item := range request.OrderItems {
		items = append(items, domain.ShippingItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	result, err := a.api.Ship(request.OrderId, items)
	if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to ship: %v", err)).Err()
	}

	return &shipping.ShipOrderResponse{EstimatedDays: result.EstimatedDays}, nil
}

func (a Adapter) Run() {
	var err error
	listen, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Fatalf("failed to listen on port %d, error: %v", a.port, err)
	}

	grpcServer := grpc.NewServer()
	shipping.RegisterShippingServer(grpcServer, a)

	if config.GetEnv() == "development" {
		reflection.Register(grpcServer)
	}

	log.Printf("Shipping service running on port %d", a.port)
	if err := grpcServer.Serve(listen); err != nil {
		log.Fatalf("failed to serve grpc on port %d", a.port)
	}
}
