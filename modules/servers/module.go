package servers

import (
	"github.com/NatthawutSK/ri-shop/modules/appinfo/appinfoHandlers"
	"github.com/NatthawutSK/ri-shop/modules/appinfo/appinfoRepositories"
	"github.com/NatthawutSK/ri-shop/modules/appinfo/appinfoUsecases"
	"github.com/NatthawutSK/ri-shop/modules/files/filesHandlers"
	"github.com/NatthawutSK/ri-shop/modules/files/filesUsecases"
	"github.com/NatthawutSK/ri-shop/modules/middlewares/middlewaresHandlers"
	"github.com/NatthawutSK/ri-shop/modules/middlewares/middlewaresRepositories"
	"github.com/NatthawutSK/ri-shop/modules/middlewares/middlewaresUsecases"
	"github.com/NatthawutSK/ri-shop/modules/monitor/monitorHandlers"
	"github.com/NatthawutSK/ri-shop/modules/orders/ordersHandlers"
	"github.com/NatthawutSK/ri-shop/modules/orders/ordersRepositories"
	"github.com/NatthawutSK/ri-shop/modules/orders/ordersUsecases"
	"github.com/NatthawutSK/ri-shop/modules/products/productsHandlers"
	"github.com/NatthawutSK/ri-shop/modules/products/productsRepositories"
	"github.com/NatthawutSK/ri-shop/modules/products/productsUsecases"
	"github.com/NatthawutSK/ri-shop/modules/users/usersHandlers"
	"github.com/NatthawutSK/ri-shop/modules/users/usersRepositories"
	"github.com/NatthawutSK/ri-shop/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface{
	MonitorModule()
	UsersModule()
	AppinfoModule()
	FilesModule()
	ProductsModule()
	OrdersModule()
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


func (m *moduleFactory) FilesModule(){
	usecase := filesUsecases.FilesUsecase(m.s.cfg)
	handler := filesHandlers.FileHandler(m.s.cfg, usecase)

	router := m.r.Group("/files")

	router.Post("/upload", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UploadFiles)
	router.Patch("/delete", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteFile)

}

func (m *moduleFactory) ProductsModule(){
	fileUsecase := filesUsecases.FilesUsecase(m.s.cfg)
	repository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, fileUsecase)
	usecase := productsUsecases.ProductsUsecase(repository)
	handler := productsHandlers.ProductsHandler(usecase, m.s.cfg, fileUsecase)

	router := m.r.Group("/products")

	router.Post("/", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddProduct)
	router.Patch("/:productId", m.mid.JwtAuth(), m.mid.Authorize(2), handler.UpdateProduct)
	router.Get("/", m.mid.ApiKeyAuth(), handler.FindProduct)
	router.Get("/:productId", m.mid.ApiKeyAuth(), handler.FindOneProduct)
	router.Delete("/:productId", m.mid.JwtAuth(), m.mid.Authorize(2), handler.DeleteProduct)


}

func (m *moduleFactory) OrdersModule(){
	fileUsecase := filesUsecases.FilesUsecase(m.s.cfg)
	productRepository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, fileUsecase)

	ordersRepository := ordersRepositories.OrdersRepository(m.s.db)
	ordersUsecase := ordersUsecases.OrdersUsecase(ordersRepository, productRepository)
	ordersHandler := ordersHandlers.OrdersHandler(ordersUsecase, m.s.cfg)

	router := m.r.Group("/orders")

	router.Post("/", m.mid.JwtAuth(), ordersHandler.InsertOrder)
	router.Get("/", m.mid.JwtAuth(), m.mid.Authorize(2), ordersHandler.FindOrder)
	router.Get("/:user_id/:order_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), ordersHandler.FindOneOrder)
	
}