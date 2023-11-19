package servers

import (
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
	
	router.Post("/signup", handler.SignUpCustomer)
}
