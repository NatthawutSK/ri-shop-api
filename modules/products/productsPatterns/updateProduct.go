package productsPatterns

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/NatthawutSK/ri-shop/config"
	"github.com/NatthawutSK/ri-shop/modules/entities"
	"github.com/NatthawutSK/ri-shop/modules/files"
	"github.com/NatthawutSK/ri-shop/modules/files/filesUsecases"
	"github.com/NatthawutSK/ri-shop/modules/products"
	"github.com/jmoiron/sqlx"
)

type IUpdateProductBuilder interface {
	initTransaction() error
	initQuery()
	updateTitleQuery()
	updateDescriptionQuery()
	updatePriceQuery()
	updateCategory() error
	insertImages() error
	getOldImages() []*entities.Image
	deleteOldImages() error
	closeQuery()
	updateProduct() error
	getQueryFields() []string
	getValues() []any
	getQuery() string
	setQuery(query string)
	getImagesLen() int
	commit() error
}

type updateProductBuilder struct{
	db             *sqlx.DB
	tx             *sqlx.Tx
	req            *products.Products
	filesUsecases  filesUsecases.IFilesUsecase
	query          string
	queryFields    []string
	lastStackIndex int
	values         []any
	cfg 		   config.IConfig
}

func UpdateProductBuilder(db *sqlx.DB, req *products.Products, fileUsecase filesUsecases.IFilesUsecase, cfg config.IConfig) IUpdateProductBuilder {
	return &updateProductBuilder{
		db:             db,
		req:            req,
		filesUsecases:  fileUsecase,
		queryFields:    make([]string, 0),
		values:         make([]any, 0),
		lastStackIndex: 0,
		cfg: cfg,
	}
}

type updateProductEngineer struct {
	builder IUpdateProductBuilder
}

func UpdateProductEngineer(builder IUpdateProductBuilder) *updateProductEngineer {
	return &updateProductEngineer{
		builder: builder,
	}
}


func (b *updateProductBuilder) initTransaction() error {

	tx, err := b.db.BeginTxx(context.Background(), nil)
	if err != nil {
		return err
	}
	b.tx = tx
	return nil
}

func (b *updateProductBuilder) initQuery() {
	b.query += `
	UPDATE "products" SET`
}

func (b *updateProductBuilder) updateTitleQuery() {
	if b.req.Title != "" {
		b.values = append(b.values, b.req.Title)
		b.lastStackIndex = len(b.values)

		b.queryFields = append(b.queryFields, fmt.Sprintf(`
		"title" = $%d`, b.lastStackIndex))
	}
}

func (b *updateProductBuilder) updateDescriptionQuery() {
	if b.req.Description != "" {
		b.values = append(b.values, b.req.Description)
		b.lastStackIndex = len(b.values)

		b.queryFields = append(b.queryFields, fmt.Sprintf(`
		"description" = $%d`, b.lastStackIndex))
	}
}

func (b *updateProductBuilder) updatePriceQuery() {
	if b.req.Price != 0 {
		b.values = append(b.values, b.req.Price)
		b.lastStackIndex = len(b.values)

		b.queryFields = append(b.queryFields, fmt.Sprintf(`
		"price" = $%d`, b.lastStackIndex))
	}
}

func (b *updateProductBuilder) updateCategory() error {

	if b.req.Category == nil {
		return nil
	}
	if b.req.Category.Id == 0 {
		return nil
	}

	query := `
	UPDATE "products_categories" SET
		"category_id" = $1
	WHERE "product_id" = $2;`

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		b.req.Category.Id,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update products_categories failed: %v", err)
	}
	return nil
}

func (b *updateProductBuilder) insertImages() error {
	query := `
	INSERT INTO "images" (
		"filename",
		"url",
		"product_id"
	)
	VALUES`

	valueStack := make([]any, 0)
	var index int
	for i := range b.req.Images {
		valueStack = append(valueStack,
			b.req.Images[i].FileName,
			b.req.Images[i].Url,
			b.req.Id,
		)

		if i != len(b.req.Images)-1 {
			query += fmt.Sprintf(`
			($%d, $%d, $%d),`, index+1, index+2, index+3)
		} else {
			query += fmt.Sprintf(`
			($%d, $%d, $%d);`, index+1, index+2, index+3)
		}
		index += 3
	}

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		valueStack...,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("insert images failed: %v", err)
	}
	return nil
}

func (b *updateProductBuilder) getOldImages() []*entities.Image {
	query := `
	SELECT
		"id",
		"filename",
		"url"
	FROM "images"
	WHERE "product_id" = $1;`

	images := make([]*entities.Image, 0)
	if err := b.db.Select(
		&images,
		query,
		b.req.Id,
	); err != nil {
		return make([]*entities.Image, 0)
	}
	return images
}

func (b *updateProductBuilder) deleteOldImages() error {
	query := `
	DELETE FROM "images"
	WHERE "product_id" = $1;`

	images := b.getOldImages()
	if len(images) > 0 {
		deleteFileReq := make([]*files.DeleteFileReq, 0)
		for _,img := range images {
			parsedURL, err := url.Parse(img.Url)
			if err != nil {
				fmt.Println("Error parsing URL:", err)
			}	

			// Get the path from the parsed URL
			path := parsedURL.Path

			// Remove the leading '/' character from the path
			path = strings.TrimPrefix(path, fmt.Sprintf("/%s/", b.cfg.App().GCPBucket()))
			deleteFileReq = append(deleteFileReq, &files.DeleteFileReq{
				Destination: fmt.Sprint(path),
			})
		}
		 
		if err := b.filesUsecases.DeleteFileOnGCP(deleteFileReq) ; err != nil {
			return fmt.Errorf("delete old images failed: %v", err)
		}
			
	}

	if _, err := b.tx.ExecContext(
		context.Background(),
		query,
		b.req.Id,
	); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("delete images failed: %v", err)
	}
	return nil
}

func (b *updateProductBuilder) closeQuery() {
	b.values = append(b.values, b.req.Id)
	b.lastStackIndex = len(b.values)

	b.query += fmt.Sprintf(`
	WHERE "id" = $%d`, b.lastStackIndex)
}

func (b *updateProductBuilder) updateProduct() error {
	if _, err := b.tx.ExecContext(context.Background(), b.query, b.values...); err != nil {
		b.tx.Rollback()
		return fmt.Errorf("update product failed: %v", err)
	}
	return nil
}

func (b *updateProductBuilder) getQueryFields() []string {
	return b.queryFields
}

func (b *updateProductBuilder) getValues() []any {
	return b.values
}

func (b *updateProductBuilder) getQuery() string {
	return b.query
}

func (b *updateProductBuilder) setQuery(query string) {
	b.query = query
}

func (b *updateProductBuilder) getImagesLen() int {
	return len(b.req.Images)
}

func (b *updateProductBuilder) commit() error {
	if err := b.tx.Commit(); err != nil {
		return fmt.Errorf("commit failed: %v", err)
	}
	return nil
}

func (en *updateProductEngineer) UpdateProduct() error {

	en.builder.initTransaction();
	en.builder.initQuery()
	en.sumQueryFields()
	en.builder.closeQuery()

	// update category
	if err := en.builder.updateCategory(); err != nil {
		return fmt.Errorf("update category failed: %v", err)
	}

	// update product
	if err := en.builder.updateProduct(); err != nil {
		return fmt.Errorf("update product failed: %v", err)
	}

	fmt.Print("len image", en.builder.getImagesLen())
	if en.builder.getImagesLen() > 0 {
		// delete old images
		if err := en.builder.deleteOldImages(); err != nil {
			return fmt.Errorf("delete old images failed: %v", err)
		}

		// insert new images
		if err := en.builder.insertImages(); err != nil {
			return fmt.Errorf("insert new images failed: %v", err)
		}
	}

	// commit transaction
	if err := en.builder.commit(); err != nil {
		return fmt.Errorf("commit failed: %v", err)
	}
	
	return nil
}

func (en *updateProductEngineer) sumQueryFields() {
	en.builder.updateTitleQuery()
	en.builder.updateDescriptionQuery()
	en.builder.updatePriceQuery()

	fields := en.builder.getQueryFields()

	for i := range fields {
		query := en.builder.getQuery()
		if i != len(fields) - 1 {
			en.builder.setQuery(query + fields[i] + ",")
		} else {
			en.builder.setQuery(query + fields[i])
		}
	}

}