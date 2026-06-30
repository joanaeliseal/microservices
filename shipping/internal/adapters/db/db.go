package db

import (
	"fmt"

	"github.com/joanaeliseal/microservices/shipping/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Shipment struct {
	gorm.Model
	OrderID       int64
	EstimatedDays int32
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dataSourceURL string) (*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(dataSourceURL), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}

	err := db.AutoMigrate(&Shipment{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}

	return &Adapter{db: db}, nil
}

func (a Adapter) Save(shipment *domain.Shipment) error {
	shipmentModel := Shipment{
		OrderID:       shipment.OrderID,
		EstimatedDays: shipment.EstimatedDays,
	}

	res := a.db.Create(&shipmentModel)
	if res.Error == nil {
		shipment.ID = int64(shipmentModel.ID)
	}

	return res.Error
}

func (a Adapter) Get(id string) (domain.Shipment, error) {
	var shipmentEntity Shipment
	res := a.db.First(&shipmentEntity, id)

	shipment := domain.Shipment{
		ID:            int64(shipmentEntity.ID),
		OrderID:       shipmentEntity.OrderID,
		EstimatedDays: shipmentEntity.EstimatedDays,
	}

	return shipment, res.Error
}
