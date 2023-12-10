package servers

import (
	"github.com/NatthawutSK/ri-shop/modules/products/productsHandlers"
	"github.com/NatthawutSK/ri-shop/modules/products/productsRepositories"
	"github.com/NatthawutSK/ri-shop/modules/products/productsUsecases"
)

type IProductModule interface {
	Init()
	Repository() productsRepositories.IProductsRepository
	Usecase() productsUsecases.IProductsUsecase
	Handler() productsHandlers.IProductsHandler
}

type ProductsModule struct {
	*moduleFactory
	repository productsRepositories.IProductsRepository
	usecase    productsUsecases.IProductsUsecase
	handler    productsHandlers.IProductsHandler
}

func (m *moduleFactory) ProductsModule() IProductModule {
	repository := productsRepositories.ProductsRepository(m.s.db, m.s.cfg, m.FilesModule().Usecase())
	usecase := productsUsecases.ProductsUsecase(repository)
	handler := productsHandlers.ProductsHandler(usecase, m.s.cfg, m.FilesModule().Usecase())

	return &ProductsModule{
		moduleFactory: m,
		repository:    repository,
		usecase:       usecase,
		handler:       handler,
	}
}

func (p *ProductsModule) Init() {
	router := p.r.Group("/products")

	router.Post("/", p.mid.JwtAuth(), p.mid.Authorize(2), p.handler.AddProduct)
	router.Patch("/:productId", p.mid.JwtAuth(), p.mid.Authorize(2), p.handler.UpdateProduct)
	router.Get("/", p.mid.ApiKeyAuth(), p.handler.FindProduct)
	router.Get("/:productId", p.mid.ApiKeyAuth(), p.handler.FindOneProduct)
	router.Delete("/:productId", p.mid.JwtAuth(), p.mid.Authorize(2), p.handler.DeleteProduct)
}

func (p *ProductsModule) Repository() productsRepositories.IProductsRepository { return p.repository }
func (p *ProductsModule) Usecase() productsUsecases.IProductsUsecase           { return p.usecase }
func (p *ProductsModule) Handler() productsHandlers.IProductsHandler           { return p.handler }
