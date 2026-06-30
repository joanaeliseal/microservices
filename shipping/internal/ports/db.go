package ports

import "github.com/joanaeliseal/microservices/shipping/internal/application/core/domain"

type DBPort interface {
	Save(shipment *domain.Shipment) error
	Get(id string) (domain.Shipment, error)
}
