package shipping

import (
	"context"
	"log"
	"time"

	grpc_retry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	"github.com/joanaeliseal/microservices-proto/golang/shipping"
	"github.com/joanaeliseal/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type Adapter struct {
	shipping shipping.ShippingClient
}

func NewAdapter(shippingServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption

	// Interceptor de retry com backoff linear (mesmo padrão do Payment)
	opts = append(opts,
		grpc.WithUnaryInterceptor(grpc_retry.UnaryClientInterceptor(
			grpc_retry.WithCodes(codes.Unavailable, codes.ResourceExhausted),
			grpc_retry.WithMax(5),
			grpc_retry.WithBackoff(grpc_retry.BackoffLinear(time.Second)),
		)),
	)

	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))

	conn, err := grpc.NewClient(shippingServiceUrl, opts...)
	if err != nil {
		return nil, err
	}

	client := shipping.NewShippingClient(conn)
	return &Adapter{shipping: client}, nil
}

func (a *Adapter) Ship(order *domain.Order) (int32, error) {
	// Deadline individual de 2 segundos para esta chamada
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var items []*shipping.OrderItem
	for _, item := range order.OrderItems {
		items = append(items, &shipping.OrderItem{
			ProductCode: item.ProductCode,
			Quantity:    item.Quantity,
		})
	}

	resp, err := a.shipping.Ship(ctx, &shipping.ShipOrderRequest{
		OrderId:    order.ID,
		OrderItems: items,
	})

	if err != nil {
		if status.Code(err) == codes.DeadlineExceeded {
			log.Printf("[Shipping] Timeout excedido para pedido %d: %v", order.ID, err)
		}
		return 0, err
	}

	return resp.EstimatedDays, nil
}
