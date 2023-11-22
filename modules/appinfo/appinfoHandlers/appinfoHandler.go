package appinfoHandlers

import (
	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/appinfo/appinfoUsecases"
	"github.com/NatthawutSK/ri-shop/modules/entities"
	riAuth "github.com/NatthawutSK/ri-shop/pkg/riauth"
	"github.com/gofiber/fiber/v2"
)


type appinfoHandlersErrCode string

const (
	generateApiKeyErr appinfoHandlersErrCode = "appinfo-001"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
}

type appinfoHandler struct {
	cfg config.IConfig
	appinfoUsecase appinfoUsecases.IAppinfoUsecase
}

func AppinfoHandler(appinfoUsecase appinfoUsecases.IAppinfoUsecase, cfg config.IConfig) IAppinfoHandler {
	return &appinfoHandler{
		appinfoUsecase: appinfoUsecase,
		cfg: cfg,
	}
}


func (h *appinfoHandler) GenerateApiKey(c *fiber.Ctx) error {
	apiKey, err := riAuth.NewRiAuth(
		riAuth.ApiKey,
		h.cfg.Jwt(),
		nil,
	)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(generateApiKeyErr),
			err.Error(),
		).Res()

	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		&struct {
			Key string `json:"key"`
		}{
			Key: apiKey.SignToken(),
		},
	).Res()
	
}