package ports

import "github.com/joanaeliseal/microservices/order/internal/application/core/domain"

type DBPort interface {
	Get(id string) (domain.Order, error)
	Save(*domain.Order) error
	ProductsExist(productCodes []string) (bool, []string, error) // retorna (existe, códigos ausentes, erro)
}
