package ordersHandlers

import (
	"strings"

	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/entities"
	"github.com/NatthawutSK/ri-shop/modules/orders/ordersUsecases"
	"github.com/gofiber/fiber/v2"
)

type ordersHandlerErrCode string

const (
	findOneOrderErr ordersHandlerErrCode = "orders-001"
)

type IOrdersHandler interface{
	FindOneOrder(c *fiber.Ctx) error
}

type ordersHandler struct {
	orderUsecase ordersUsecases.IOrdersUsecase
	cfg 		config.IConfig
}

func OrdersHandler(orderUsecase ordersUsecases.IOrdersUsecase, cfg config.IConfig) IOrdersHandler {
	return &ordersHandler{
		orderUsecase: orderUsecase,
		cfg: cfg,
	}
}

func (h *ordersHandler) FindOneOrder(c *fiber.Ctx) error {

	orderId := strings.Trim(c.Params("order_id"), " ")
	order, err := h.orderUsecase.FindOneOrder(orderId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(findOneOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		order,
	).Res()

	
}

