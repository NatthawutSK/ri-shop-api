package productsUsecases

import (
	"math"

	"github.com/NatthawutSK/ri-shop/modules/entities"
	"github.com/NatthawutSK/ri-shop/modules/products"
	"github.com/NatthawutSK/ri-shop/modules/products/productsRepositories"
)

type IProductsUsecase interface{
	FindOneProduct(productId string) (*products.Products, error)
	FindProduct(req *products.ProductFilter) *entities.PaginateRes
	AddProduct(req *products.Products) (*products.Products, error)
	UpdateProduct(req *products.Products) (*products.Products, error)
}

type productsUsecase struct {
	productsRepository productsRepositories.IProductsRepository
}

func ProductsUsecase(productsRepository productsRepositories.IProductsRepository) IProductsUsecase {
	return &productsUsecase{
		productsRepository: productsRepository,
	}
}

func (u *productsUsecase) FindOneProduct(productId string) (*products.Products, error) {
	product, err := u.productsRepository.FindOneProduct(productId)
	if err != nil {
		return nil, err
	}
	return product, nil
}


func (u *productsUsecase) FindProduct(req *products.ProductFilter) *entities.PaginateRes {
	products, count := u.productsRepository.FindProduct(req)
	return &entities.PaginateRes{
		Data: products,
		TotalItem: count,
		Page: req.Page,
		Limit: req.Limit,
		TotalPage: int(math.Ceil(float64(count) / float64(req.Limit))),

	}
	
}

func (u *productsUsecase) AddProduct(req *products.Products) (*products.Products, error) {
	product, err := u.productsRepository.InsertProduct(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (u *productsUsecase) UpdateProduct(req *products.Products) (*products.Products, error) {
	product, err := u.productsRepository.UpdateProduct(req)
	if err != nil {
		return nil, err
	}
	return product, nil
}