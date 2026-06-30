package domain

type ShippingItem struct {
	ProductCode string
	Quantity    int32
}

type Shipment struct {
	ID            int64
	OrderID       int64
	EstimatedDays int32
}

func NewShipment(orderID int64, estimatedDays int32) Shipment {
	return Shipment{
		OrderID:       orderID,
		EstimatedDays: estimatedDays,
	}
}
