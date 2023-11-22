package servers

import (
	"github.com/NatthawutSK/ri-shop/modules/appinfo/appinfoHandlers"
	"github.com/NatthawutSK/ri-shop/modules/appinfo/appinfoRepositories"
	"github.com/NatthawutSK/ri-shop/modules/appinfo/appinfoUsecases"
	"github.com/NatthawutSK/ri-shop/modules/middlewares/middlewaresHandlers"
	"github.com/NatthawutSK/ri-shop/modules/middlewares/middlewaresRepositories"
	"github.com/NatthawutSK/ri-shop/modules/middlewares/middlewaresUsecases"
	"github.com/NatthawutSK/ri-shop/modules/monitor/monitorHandlers"
	"github.com/NatthawutSK/ri-shop/modules/users/usersHandlers"
	"github.com/NatthawutSK/ri-shop/modules/users/usersRepositories"
	"github.com/NatthawutSK/ri-shop/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface{
	MonitorModule()
	UsersModule()
	AppinfoModule()
}


type moduleFactory struct {
	r fiber.Router
	s *server
	mid middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		r: r,
		s: s,
		mid: mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler{
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}


func (m *moduleFactory) MonitorModule() {
	handle := monitorHandlers.MonitorHandler(m.s.cfg)

	m.r.Get("/", handle.HealthCheck)
}


func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepositoryHandler(m.s.db)
	usecase := usersUsecases.UserUsecaseHandler(repository, m.s.cfg)
	handler := usersHandlers.UsersHandler(m.s.cfg, usecase)

	router := m.r.Group("/users")
	
	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), handler.SignUpAdmin)
	router.Get("/admin/secret", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)
	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
}


func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.s.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(usecase, m.s.cfg)

	router := m.r.Group("/appinfo")

	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)
	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)
	router.Post("/categories",  m.mid.JwtAuth(), m.mid.Authorize(2), handler.InsertCategory)
	router.Delete("/:categoryId/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteCategory)
}
