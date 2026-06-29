package payment_adapter

import (
	"context"

	"github.com/joanaeliseal/microservices-proto/golang/payment"
	"github.com/joanaeliseal/microservices/order/internal/application/core/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Adapter struct {
	payment payment.PaymentClient // codigo gerado a partir do arquivo protobuf
}

func NewAdapter(paymentServiceUrl string) (*Adapter, error) {
	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(paymentServiceUrl, opts...)
	if err != nil {
		return nil, err
	}
	client := payment.NewPaymentClient(conn) // inicializa o stub
	return &Adapter{payment: client}, nil
}

func (a *Adapter) Charge(order *domain.Order) error {
	_, err := a.payment.CreatePayment(context.Background(), &payment.CreatePaymentRequest{
		UserId:     order.UserID,
		OrderId:    order.ID,
		TotalPrice: order.TotalPrice,
	})
	return err
}

func (a Adapter) Create(ctx context.Context, request *payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	log.WithContext(ctx).Info("Creating payment...")

	newPayment := domain.NewPayment(request.UserId, request.OrderId, request.TotalPrice)
	result, err := a.api.Charge(ctx, newPayment)
	code := status.Code(err)
	 if code == codes.InvalidArgument {
		return nil, err
	} else if err != nil {
		return nil, status.New(codes.Internal, fmt.Sprintf("failed to charge: %v", err)).Err()
	}
	return &payment.CreatePaymentResponse{PaymentId: result.ID}, nil