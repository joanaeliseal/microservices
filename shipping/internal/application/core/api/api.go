package api

import (
	"github.com/joanaeliseal/microservices/shipping/internal/application/core/domain"
	"github.com/joanaeliseal/microservices/shipping/internal/ports"
)

type Application struct {
	db ports.DBPort
}

func NewApplication(db ports.DBPort) *Application {
	return &Application{
		db: db,
	}
}

// calculateDeliveryDays calcula o prazo de entrega
// Mínimo: 1 dia, +1 dia a cada 5 unidades
func calculateDeliveryDays(items []domain.ShippingItem) int32 {
	var totalQuantity int32
	for _, item := range items {
		totalQuantity += item.Quantity
	}
	days := int32(1) + (totalQuantity / 5)
	return days
}

func (a Application) Ship(orderID int64, items []domain.ShippingItem) (domain.Shipment, error) {
	days := calculateDeliveryDays(items)
	shipment := domain.NewShipment(orderID, days)

	err := a.db.Save(&shipment)
	if err != nil {
		return domain.Shipment{}, err
	}

	return shipment, nil
}
