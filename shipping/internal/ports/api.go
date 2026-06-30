package ports

import "github.com/joanaeliseal/microservices/shipping/internal/application/core/domain"

type APIPort interface {
	Ship(orderID int64, items []domain.ShippingItem) (domain.Shipment, error)
}
