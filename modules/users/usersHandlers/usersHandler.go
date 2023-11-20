package usersHandlers

import (
	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/entities"
	"github.com/NatthawutSK/ri-shop/modules/users"
	"github.com/NatthawutSK/ri-shop/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type userHandlerErrCode = string

const (
	signUpCustomerErr userHandlerErrCode = "users-001"
	signInCustomerErr userHandlerErrCode = "users-002"
)

type IUsersHandler interface{
	SignUpCustomer(c *fiber.Ctx) error
	SignIn(c *fiber.Ctx) error

}

type usersHandler struct {
	cfg config.IConfig
	userUsecase usersUsecases.IUserUsecase

}

func UsersHandler(cfg config.IConfig, UserUsecase usersUsecases.IUserUsecase) IUsersHandler {
	return &usersHandler{
		cfg: cfg,
		userUsecase: UserUsecase,
	}
}

func (h *usersHandler) SignUpCustomer(c *fiber.Ctx) error {
	// Request Body parser
	req := new(users.UserRegisterReq)

	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			err.Error(),
		).Res()
	}
	// Email validation
	if !req.IsEmail() {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signUpCustomerErr),
			"email is invalid",
		).Res()
	}

	// Insert user
	result, err := h.userUsecase.InsertCustomer(req)
	if err != nil {
		switch err.Error() {
			case "username has been used":
				return entities.NewResponse(c).Error(
					fiber.ErrBadRequest.Code,
					string(signUpCustomerErr),
					err.Error(),
				).Res()
				case "email has been used":
					return entities.NewResponse(c).Error(
						fiber.ErrBadRequest.Code,
						string(signUpCustomerErr),
						err.Error(),
					).Res()

			default:
				return entities.NewResponse(c).Error(
					fiber.ErrInternalServerError.Code,
					string(signUpCustomerErr),
					err.Error(),
				).Res()
		}
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, result).Res()
}


func (h *usersHandler) SignIn(c *fiber.Ctx) error {
	req := new(users.UserCredential)
	if err := c.BodyParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInCustomerErr),
			err.Error(),
		).Res()
	}

	result, err := h.userUsecase.GetPassport(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(signInCustomerErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, result).Res()
}