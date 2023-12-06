package ordersUsecases

import (
	"fmt"
	"math"

	"github.com/NatthawutSK/ri-shop/modules/entities"
	"github.com/NatthawutSK/ri-shop/modules/orders"
	"github.com/NatthawutSK/ri-shop/modules/orders/ordersRepositories"
	"github.com/NatthawutSK/ri-shop/modules/products/productsRepositories"
)

type IOrdersUsecase interface{
	FindOneOrder(orderId string) (*orders.Order, error) 
	FindOrder(req *orders.OrderFilter) *entities.PaginateRes
	InsertOrder(req *orders.Order) (*orders.Order, error)
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

func (u *ordersUsecase) FindOrder(req *orders.OrderFilter) *entities.PaginateRes {
	orders, count := u.ordersRepository.FindOrder(req)

	return &entities.PaginateRes{
		Data: orders,
		Page: req.Page,
		Limit: req.Limit,
		TotalPage: int(math.Ceil(float64(count)/float64(req.Limit))),
		TotalItem: count,
	}
}

func (u *ordersUsecase) InsertOrder(req *orders.Order) (*orders.Order, error) {
	// Check product is exist and correct price
	for i := range req.Products {
		if req.Products[i].Product == nil {
			return nil, fmt.Errorf("product is required")
		
		}

		prod, err := u.productsRepository.FindOneProduct(req.Products[i].Product.Id)
		if err != nil {
			return nil, fmt.Errorf("find one product failed : %v", err)
		}

		// set price from product
		req.TotalPaid += req.Products[i].Product.Price * float64(req.Products[i].Qty)
		req.Products[i].Product = prod
	}

	orderId, err := u.ordersRepository.InsertOrder(req)
	if err != nil {
		return nil, err
	}

	order, err := u.ordersRepository.FindOneOrder(orderId)
	if err != nil {
		return nil, err
	}

	return order, nil
}