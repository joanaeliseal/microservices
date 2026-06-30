package ports

import "github.com/joanaeliseal/microservices/order/internal/application/core/domain"

type ShippingPort interface {
	Ship(order *domain.Order) (int32, error) // retorna dias estimados
}
