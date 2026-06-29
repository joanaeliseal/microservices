package api

import (
	"github.com/joanaeliseal/microservices/order/internal/application/core/domain"
	"github.com/joanaeliseal/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const MaxTotalItems = 50

type Application struct {
	db      ports.DBPort
	payment ports.PaymentPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort) *Application {
	return &Application{
		db:      db,
		payment: payment,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	// Validar quantidade total de itens
	var totalQuantity int32
	for _, item := range order.OrderItems {
		totalQuantity += item.Quantity
	}

	if totalQuantity > MaxTotalItems {
		return domain.Order{}, status.Errorf(
			codes.InvalidArgument,
			"pedido não pode ser realizado: quantidade total de itens (%d) excede o limite máximo permitido de %d itens",
			totalQuantity,
			MaxTotalItems,
		)
	}

	err := a.db.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}

	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		return domain.Order{}, paymentErr
	}

	return order, nil
}
