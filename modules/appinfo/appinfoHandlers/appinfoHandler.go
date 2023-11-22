package appinfoHandlers

import (
	"strconv"
	"strings"

	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/appinfo"
	"github.com/NatthawutSK/ri-shop/modules/appinfo/appinfoUsecases"
	"github.com/NatthawutSK/ri-shop/modules/entities"
	riAuth "github.com/NatthawutSK/ri-shop/pkg/riauth"
	"github.com/gofiber/fiber/v2"
)


type appinfoHandlersErrCode string

const (
	generateApiKeyErr appinfoHandlersErrCode = "appinfo-001"
	FindCategoryErr appinfoHandlersErrCode = "appinfo-002"
	InsertCategoryErr appinfoHandlersErrCode = "appinfo-003"
	DeleteCategoryErr appinfoHandlersErrCode = "appinfo-004"
)

type IAppinfoHandler interface {
	GenerateApiKey(c *fiber.Ctx) error
	FindCategory(c *fiber.Ctx) error
	InsertCategory(c *fiber.Ctx) error
	DeleteCategory(c *fiber.Ctx) error
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


func (h *appinfoHandler) FindCategory(c *fiber.Ctx) error {
	req := new(appinfo.CategoryFilter)
	//if only one parameter
	//use c.Query("title")
	if err := c.QueryParser(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(FindCategoryErr),
			err.Error(),
		).Res()
	}

	category, err := h.appinfoUsecase.FindCategory(req)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(FindCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(
		fiber.StatusOK,
		category,
	).Res()
	
	
}

func (h *appinfoHandler) InsertCategory(c *fiber.Ctx) error {
	req := make([]*appinfo.Category, 0)
	if err := c.BodyParser(&req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(InsertCategoryErr),
			err.Error(),
		).Res()
	}

	if len(req) == 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(InsertCategoryErr),
			"category request body is empty",
		).Res()
	}

	if err := h.appinfoUsecase.InsertCategory(req); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(InsertCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusCreated, req).Res()
}

func (h *appinfoHandler) DeleteCategory(c *fiber.Ctx) error {
	categoryId := strings.Trim(c.Params("categoryId"), " ")
	categoryIdInt, err := strconv.Atoi(categoryId)
	if err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(DeleteCategoryErr),
			"category id type is invalid",
		).Res()
	} 
	if categoryIdInt <= 0 {
		return entities.NewResponse(c).Error(
			fiber.ErrBadRequest.Code,
			string(DeleteCategoryErr),
			"category id must more than 0",
		).Res()
	}

	if err := h.appinfoUsecase.DeleteCategory(categoryIdInt); err != nil {
		return entities.NewResponse(c).Error(
			fiber.ErrInternalServerError.Code,
			string(DeleteCategoryErr),
			err.Error(),
		).Res()
	}

	return entities.NewResponse(c).Success(fiber.StatusOK, 
		&struct{
			CategoryId int `json:"category_id"`
		}{
			CategoryId: categoryIdInt,
		},
	).Res()
}