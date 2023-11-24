package productsUsecases

import (
	"github.com/NatthawutSK/ri-shop/modules/products"
	"github.com/NatthawutSK/ri-shop/modules/products/productsRepositories"
)

type IProductsUsecase interface{
	FindOneProduct(productId string) (*products.Products, error)
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