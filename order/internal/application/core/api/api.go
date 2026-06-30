package api

import (
	"fmt"
	"log"

	"github.com/joanaeliseal/microservices/order/internal/application/core/domain"
	"github.com/joanaeliseal/microservices/order/internal/ports"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const MaxTotalItems = 50

type Application struct {
	db       ports.DBPort
	payment  ports.PaymentPort
	shipping ports.ShippingPort
}

func NewApplication(db ports.DBPort, payment ports.PaymentPort, shipping ports.ShippingPort) *Application {
	return &Application{
		db:       db,
		payment:  payment,
		shipping: shipping,
	}
}

func (a Application) PlaceOrder(order domain.Order) (domain.Order, error) {
	// Validar quantidade total de itens
	var totalQuantity int32
	var productCodes []string
	for _, item := range order.OrderItems {
		totalQuantity += item.Quantity
		productCodes = append(productCodes, item.ProductCode)
	}

	if totalQuantity > MaxTotalItems {
		return domain.Order{}, status.Errorf(
			codes.InvalidArgument,
			"pedido não pode ser realizado: quantidade total de itens (%d) excede o limite máximo permitido de %d itens",
			totalQuantity,
			MaxTotalItems,
		)
	}

	// Verificar existência dos produtos no estoque
	exists, missing, err := a.db.ProductsExist(productCodes)
	if err != nil {
		return domain.Order{}, status.Errorf(codes.Internal, "falha ao verificar estoque: %v", err)
	}
	if !exists {
		return domain.Order{}, status.Errorf(
			codes.NotFound,
			"produtos não encontrados no estoque: %v",
			fmt.Sprintf("%v", missing),
		)
	}

	// Salvar pedido no banco
	err = a.db.Save(&order)
	if err != nil {
		return domain.Order{}, err
	}

	// Processar pagamento
	paymentErr := a.payment.Charge(&order)
	if paymentErr != nil {
		order.Status = "Canceled"
		return domain.Order{}, paymentErr
	}

	// Só aciona Shipping se o pagamento foi OK
	estimatedDays, shippingErr := a.shipping.Ship(&order)
	if shippingErr != nil {
		// Pagamento já feito; logamos o erro mas não cancelamos o pedido
		log.Printf("[Order] Shipping failed for order %d: %v", order.ID, shippingErr)
	} else {
		log.Printf("[Order] Shipping scheduled for order %d: %d days", order.ID, estimatedDays)
	}

	order.Status = "Paid"
	return order, nil
}
