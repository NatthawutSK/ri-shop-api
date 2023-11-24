package productsHandlers

import (
	"strings"

	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/entities"
	"github.com/NatthawutSK/ri-shop/modules/files/filesUsecases"
	"github.com/NatthawutSK/ri-shop/modules/products/productsUsecases"
	"github.com/gofiber/fiber/v2"
)


type productsHandlerErrCode string

const (
	findOneProductErr productsHandlerErrCode = "products-001"
)

type IProductsHandler interface{
	FindOneProduct(c *fiber.Ctx) error
}

type productsHandler struct {
	productsUsecase productsUsecases.IProductsUsecase
	cfg config.IConfig
	fileUsecase filesUsecases.IFilesUsecase
}

func ProductsHandler(productsUsecase productsUsecases.IProductsUsecase, cfg config.IConfig, fileUsecase filesUsecases.IFilesUsecase) IProductsHandler {
	return &productsHandler{
		productsUsecase: productsUsecase,
		cfg: cfg,
		fileUsecase: fileUsecase,
	}
}

func (h *productsHandler) FindOneProduct(c *fiber.Ctx) error {
	productId := strings.Trim(c.Params("productId"), " ")
	product, err := h.productsUsecase.FindOneProduct(productId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneProductErr),
			err.Error(),

		).Res()
	}
	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		product,
	).Res()
}