package servers

import (
	"github.com/NatthawutSK/ri-shop/modules/files/filesHandlers"
	"github.com/NatthawutSK/ri-shop/modules/files/filesUsecases"
)

type IFilesModule interface {
	Init()
	Usecase() filesUsecases.IFilesUsecase
	Handler() filesHandlers.IFileHandler
}

type filesModule struct {
	*moduleFactory
	usecase filesUsecases.IFilesUsecase
	handler filesHandlers.IFileHandler
}

func (m *moduleFactory) FilesModule() IFilesModule {
	usecase := filesUsecases.FilesUsecase(m.s.cfg)
	handler := filesHandlers.FileHandler(m.s.cfg, usecase)

	return &filesModule{
		moduleFactory: m,
		usecase:       usecase,
		handler:       handler,
	}
}

func (f *filesModule) Init() {
	router := f.r.Group("/files")

	router.Post("/upload", f.mid.JwtAuth(), f.mid.Authorize(2), f.handler.UploadFiles)
	router.Patch("/delete", f.mid.JwtAuth(), f.mid.Authorize(2), f.handler.DeleteFile)
}

func (f *filesModule) Usecase() filesUsecases.IFilesUsecase { return f.usecase }
func (f *filesModule) Handler() filesHandlers.IFileHandler  { return f.handler }
