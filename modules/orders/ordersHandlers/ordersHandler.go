package ordersHandlers

import (
	"strings"
	"time"

	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/entities"
	"github.com/NatthawutSK/ri-shop/modules/orders"
	"github.com/NatthawutSK/ri-shop/modules/orders/ordersUsecases"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type ordersHandlerErrCode string

const (
	findOneOrderErr ordersHandlerErrCode = "orders-001"
	findOrderErr    ordersHandlerErrCode = "orders-002"
	insertOrderErr  ordersHandlerErrCode = "orders-003"
	updateOrderErr  ordersHandlerErrCode = "orders-004"
)

type IOrdersHandler interface {
	FindOneOrder(c *fiber.Ctx) error
	FindOrder(c *fiber.Ctx) error
	InsertOrder(c *fiber.Ctx) error
	UpdateOrder(c *fiber.Ctx) error
}

type ordersHandler struct {
	orderUsecase ordersUsecases.IOrdersUsecase
	cfg          config.IConfig
}

func OrdersHandler(orderUsecase ordersUsecases.IOrdersUsecase, cfg config.IConfig) IOrdersHandler {
	return &ordersHandler{
		orderUsecase: orderUsecase,
		cfg:          cfg,
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
		SortReq:       &entities.SortReq{},
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

func (h *ordersHandler) InsertOrder(c *fiber.Ctx) error {
	userId := c.Locals("userId").(string)

	req := &orders.Order{
		Products: make([]*orders.ProductsOrder, 0),
	}

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertOrderErr),
			err.Error(),
		).Res()
	}

	if len(req.Products) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(insertOrderErr),
			"products are empty",
		).Res()
	}

	// set user_id ให้เป็นของตัวเองเสมอ ยกเว้นเป็น admin
	if c.Locals("userRoleId").(int) != 2 {
		req.UserId = userId
	}

	req.Status = "waiting"
	req.TotalPaid = 0

	order, err := h.orderUsecase.InsertOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(insertOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		order,
	).Res()
}

func (h *ordersHandler) UpdateOrder(c *fiber.Ctx) error {
	orderId := strings.Trim(c.Params("order_id"), " ")

	req := new(orders.OrderUpdate)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(updateOrderErr),
			err.Error(),
		).Res()
	}

	req.Id = orderId

	statusMap := map[string]string{
		"waiting":   "waiting",
		"shipping":  "shipping",
		"completed": "completed",
		"canceled":  "canceled",
	}

	// ถ้าเป็น admin จะสามารถเปลี่ยนสถานะได้ทั้งหมด แต่ถ้าเป็น user จะสามารถเปลี่ยนเป็น canceled ได้เท่านั้น
	if c.Locals("userRoleId").(int) == 2 {
		req.Status = statusMap[strings.ToLower(req.Status)]
	} else if strings.ToLower(req.Status) != statusMap["canceled"] {
		req.Status = statusMap["canceled"]
	}

	if req.TransferSlip != nil {
		if req.TransferSlip.Id == "" {
			req.TransferSlip.Id = uuid.NewString()
		}

		if req.TransferSlip.CreatedAt == "" {
			loc, err := time.LoadLocation("Asia/Bangkok")
			if err != nil {
				return entities.NewResponse(c).Error(
					fiber.ErrInternalServerError.Code,
					string(updateOrderErr),
					err.Error(),
				).Res()
			}
			now := time.Now().In(loc)

			// YYYY-MM-DD HH:MM:SS
			// 2006-01-02 15:04:05
			req.TransferSlip.CreatedAt = now.Format("2006-01-02 15:04:05")
		}
	}

	order, err := h.orderUsecase.UpdateOrder(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(updateOrderErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusCreated,
		order,
	).Res()
}
