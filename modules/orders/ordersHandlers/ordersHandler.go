package ordersHandlers

import (
	"strings"
	"time"

	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/entities"
	"github.com/NatthawutSK/ri-shop/modules/orders"
	"github.com/NatthawutSK/ri-shop/modules/orders/ordersUsecases"
	"github.com/gofiber/fiber/v2"
)

type ordersHandlerErrCode string

const (
	findOneOrderErr ordersHandlerErrCode = "orders-001"
	findOrderErr ordersHandlerErrCode = "orders-002"
)

type IOrdersHandler interface{
	FindOneOrder(c *fiber.Ctx) error
	FindOrder(c *fiber.Ctx) error
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

func (h *ordersHandler) FindOrder(c *fiber.Ctx) error {
	req := &orders.OrderFilter{
		SortReq: &entities.SortReq{},
		PaginationReq: &entities.PaginationReq{},

	}

	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(findOrderErr),
			err.Error(),
		).Res()
	}

	// pagination
	if req.Page < 1 {
		req.Page = 1
	}

	if req.Limit < 3 {
		req.Limit = 3
	}

	// order by
	orderByMap := map[string]string{
		"id":         `"o"."id"`,
		"created_at": `"o"."created_at"`,
	}
	if orderByMap[req.OrderBy] == "" {
		req.OrderBy = orderByMap["id"]
	}

	// sort
	req.Sort = strings.ToUpper(req.Sort)
	sortMap := map[string]string{
		"DESC": "DESC",
		"ASC":  "ASC",
	}
	if sortMap[req.Sort] == "" {
		req.Sort = sortMap["DESC"]
	}

	// Date	YYYY-MM-DD
	if req.StartDate != "" {
		start, err := time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(findOrderErr),
				"start date is invalid",
			).Res()
		}
		req.StartDate = start.Format("2006-01-02")
	}
	if req.EndDate != "" {
		end, err := time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			return entities.NewResponse(c).Error(
				fiber.ErrBadRequest.Code,
				string(findOrderErr),
				"end date is invalid",
			).Res()
		}
		req.EndDate = end.Format("2006-01-02")
	}

	orders := h.orderUsecase.FindOrder(req)

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		orders,
	).Res()
}

