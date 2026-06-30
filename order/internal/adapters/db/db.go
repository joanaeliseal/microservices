package db

import (
	"fmt"

	"github.com/joanaeliseal/microservices/order/internal/application/core/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	CustomerID int64
	Status     string
	OrderItems []OrderItem
}

type OrderItem struct {
	gorm.Model
	ProductCode string
	UnitPrice   float32
	Quantity    int32
	OrderID     uint
}

type StockItem struct {
	gorm.Model
	ProductCode string `gorm:"size:100;uniqueIndex"`
	Name        string
	Quantity    int32
}

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dataSourceURL string) (*Adapter, error) {
	db, openErr := gorm.Open(mysql.Open(dataSourceURL), &gorm.Config{})
	if openErr != nil {
		return nil, fmt.Errorf("db connection error: %v", openErr)
	}

	err := db.AutoMigrate(&Order{}, &OrderItem{}, &StockItem{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}

	return &Adapter{db: db}, nil
}

func (a Adapter) Get(id string) (domain.Order, error) {
	var orderEntity Order
	res := a.db.First(&orderEntity, id)

	orderItems := make([]domain.OrderItem, 0, len(orderEntity.OrderItems))
	for _, orderItem := range orderEntity.OrderItems {
		orderItems = append(orderItems, domain.OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	order := domain.Order{
		ID:         int64(orderEntity.ID),
		CustomerID: orderEntity.CustomerID,
		Status:     orderEntity.Status,
		OrderItems: orderItems,
		CreatedAt:  orderEntity.CreatedAt.UnixNano(),
	}

	return order, res.Error
}

func (a Adapter) Save(order *domain.Order) error {
	orderItems := make([]OrderItem, 0, len(order.OrderItems))
	for _, orderItem := range order.OrderItems {
		orderItems = append(orderItems, OrderItem{
			ProductCode: orderItem.ProductCode,
			UnitPrice:   orderItem.UnitPrice,
			Quantity:    orderItem.Quantity,
		})
	}

	orderModel := Order{
		CustomerID: order.CustomerID,
		Status:     order.Status,
		OrderItems: orderItems,
	}

	res := a.db.Create(&orderModel)
	if res.Error == nil {
		order.ID = int64(orderModel.ID)
	}

	return res.Error
}

// ProductsExist verifica se os produtos existem no estoque
func (a Adapter) ProductsExist(productCodes []string) (bool, []string, error) {
	var found []StockItem
	result := a.db.Where("product_code IN ?", productCodes).Find(&found)
	if result.Error != nil {
		return false, nil, result.Error
	}

	foundCodes := make(map[string]bool)
	for _, item := range found {
		foundCodes[item.ProductCode] = true
	}

	var missing []string
	for _, code := range productCodes {
		if !foundCodes[code] {
			missing = append(missing, code)
		}
	}

	return len(missing) == 0, missing, nil
}
