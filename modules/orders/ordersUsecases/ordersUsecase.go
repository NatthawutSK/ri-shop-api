package ordersUsecases

import (
	"github.com/NatthawutSK/ri-shop/modules/orders"
	"github.com/NatthawutSK/ri-shop/modules/orders/ordersRepositories"
	"github.com/NatthawutSK/ri-shop/modules/products/productsRepositories"
)

type IOrdersUsecase interface{
	FindOneOrder(orderId string) (*orders.Order, error) 
}

type ordersUsecase struct {
	ordersRepository ordersRepositories.IOrdersRepository
	productsRepository productsRepositories.IProductsRepository
}

func OrdersUsecase(ordersRepo ordersRepositories.IOrdersRepository, productsRepo productsRepositories.IProductsRepository) IOrdersUsecase {
	return &ordersUsecase{
		ordersRepository: ordersRepo,
		productsRepository: productsRepo,
	}
}

func (u *ordersUsecase) FindOneOrder(orderId string) (*orders.Order, error) {
	order, err := u.ordersRepository.FindOneOrder(orderId) 
	if err != nil {
		return nil, err
	}
	
	return order, nil
}