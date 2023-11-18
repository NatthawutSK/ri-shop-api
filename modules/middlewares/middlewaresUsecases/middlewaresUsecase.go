package middlewaresUsecases

import "github.com/NatthawutSK/ri-shop/modules/middlewares/middlewaresRepositories"

type IMiddlewaresUsecase interface {
}

type middlewaresUsecase struct {
	middlewareRepository middlewaresRepositories.IMiddlewaresRepository
}

func MiddlewaresUsecase(middlewareRepository middlewaresRepositories.IMiddlewaresRepository) IMiddlewaresUsecase {
	return &middlewaresUsecase{
		middlewareRepository: middlewareRepository,
	}
}