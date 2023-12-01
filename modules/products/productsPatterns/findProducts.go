package productsPatterns

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/NatthawutSK/ri-shop/modules/products"
	"github.com/NatthawutSK/ri-shop/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type IFindProductBuilder interface {
	openJsonQuery()
	initQuery()
	countQuery()
	whereQuery()
	sort()
	paginate()
	closeJsonQuery()
	resetQuery()
	Result() []*products.Products
	Count() int
	PrintQuery()
}

type findProductBuilder struct {
	db             *sqlx.DB
	req            *products.ProductFilter
	query          string
	lastStackIndex int
	values         []any
}

func FindProductBuilder(db *sqlx.DB, req *products.ProductFilter) IFindProductBuilder {
	return &findProductBuilder{
		db:  db,
		req: req,
	}
}

func (b *findProductBuilder) openJsonQuery() {
	b.query += `SELECT
		array_to_json(array_agg("t"))
	FROM (`
}
func (b *findProductBuilder) initQuery() {
	
	b.query += `
		SELECT
			"p"."id",
			"p"."title",
			"p"."description",
			"p"."price",
			(
				SELECT
					to_jsonb("ct")
				FROM (
					SELECT
						"c"."id",
						"c"."title"
					FROM "categories" "c"
						LEFT JOIN "products_categories" "pc" ON "pc"."category_id" = "c"."id"
					WHERE "pc"."product_id" = "p"."id"
				) AS "ct"
			) AS "category",
			"p"."created_at",
			"p"."updated_at",
			(
				SELECT
					COALESCE(array_to_json(array_agg("it")), '[]'::json)
				FROM (
					SELECT
						"i"."id",
						"i"."filename",
						"i"."url"
					FROM "images" "i"
					WHERE "i"."product_id" = "p"."id"
				) AS "it"
			) AS "images"
		FROM "products" "p"
		WHERE 1 = 1`
}
func (b *findProductBuilder) countQuery() {
	b.query += `
		SELECT
			COUNT(*) AS "count"
		FROM "products" "p"
		WHERE 1 = 1`
}
func (b *findProductBuilder) whereQuery() {
	var queryWhere string
	queryWhereStack := make([]string, 0)

	// Id check
	if b.req.Id != "" {
		b.values = append(b.values, b.req.Id)

		queryWhereStack = append(queryWhereStack, `
		AND "p"."id" = ?`)
	}

	// Search check
	if b.req.Search != "" {
		b.values = append(
			b.values,
			"%"+strings.ToLower(b.req.Search)+"%",
			"%"+strings.ToLower(b.req.Search)+"%",
		)

		queryWhereStack = append(queryWhereStack, `
		AND (LOWER("p"."title") LIKE ? OR LOWER("p"."description") LIKE ?)`)
	}

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			queryWhere += strings.Replace(queryWhereStack[i], "?", "$"+strconv.Itoa(i+1), 1)
		} else {
			queryWhere += strings.Replace(queryWhereStack[i], "?", "$"+strconv.Itoa(i+1), 1)
			queryWhere = strings.Replace(queryWhere, "?", "$"+strconv.Itoa(i+2), 1)
		}
	}
	// Last stack record
	b.lastStackIndex = len(b.values)

	// Summary query
	b.query += queryWhere
}
func (b *findProductBuilder) sort() {
    orderByMap := map[string]string{
        "id":    "\"p\".\"id\"",
        "title": "\"p\".\"title\"",
        "price": "\"p\".\"price\"",
    }
 
    if orderByMap[strings.ToLower(b.req.OrderBy)] == "" {
        b.req.OrderBy = orderByMap["title"]
    } else {
        b.req.OrderBy = orderByMap[strings.ToLower(b.req.OrderBy)]
    }
 
    sortMap := map[string]string{
        "DESC": "DESC",
        "ASC":  "ASC",
    }
    if sortMap[strings.ToUpper(b.req.Sort)] == "" {
        b.req.Sort = sortMap["ASC"]
    } else {
        b.req.Sort = sortMap[strings.ToUpper(b.req.Sort)]
    }
 
    // โค้ดที่มีปัญหา เมื่อใช้แล้ว ORDER BY จะไม่ทำงาน
    /* b.values = append(b.values, b.req.OrderBy)
    b.query += fmt.Sprintf(`
        ORDER BY $%d %s`, b.lastStackIndex+1, b.req.Sort)
    b.lastStackIndex = len(b.values) */
 
    b.query += fmt.Sprintf(`
        ORDER BY %s %s`, b.req.OrderBy, b.req.Sort)
}
func (b *findProductBuilder) paginate() {
	// offset (page - 1)*limit
	b.values = append(b.values, (b.req.Page-1)*b.req.Limit, b.req.Limit)

	b.query += fmt.Sprintf(`	OFFSET $%d LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
	b.lastStackIndex = len(b.values)
}
func (b *findProductBuilder) closeJsonQuery() {
	b.query += `
	) AS "t";`
}
func (b *findProductBuilder) resetQuery() {
	b.query = ""
	b.values = make([]any, 0)
	b.lastStackIndex = 0
}
func (b *findProductBuilder) Result() []*products.Products {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	bytes := make([]byte, 0)
	productsData := make([]*products.Products, 0)

	if err := b.db.Get(&bytes, b.query, b.values...); err != nil {
		log.Printf("find products failed: %v\n", err)
		return make([]*products.Products, 0)
	}

	if err := json.Unmarshal(bytes, &productsData); err != nil {
		log.Printf("unmarshal products failed: %v\n", err)
		return make([]*products.Products, 0)
	}
	b.resetQuery()
	return productsData
}
func (b *findProductBuilder) Count() int {
	_, cancel := context.WithTimeout(context.Background(), time.Second*15)
	defer cancel()

	var count int
	if err := b.db.Get(&count, b.query, b.values...); err != nil {
		log.Printf("count products failed: %v\n", err)
		return 0
	}
	b.resetQuery()
	return count
}
func (b *findProductBuilder) PrintQuery() {
	utils.Debug(b.values)
	fmt.Println(b.query)
}

type findProductEngineer struct {
	builder IFindProductBuilder
}

func FindProductEngineer(builder IFindProductBuilder) *findProductEngineer {
	return &findProductEngineer{builder: builder}
}

func (en *findProductEngineer) FindProduct() IFindProductBuilder {
	en.builder.openJsonQuery()
	en.builder.initQuery()
	en.builder.whereQuery()
	en.builder.sort()
	en.builder.paginate()
	en.builder.closeJsonQuery()
	return en.builder
}

func (en *findProductEngineer) CountProduct() IFindProductBuilder {
	en.builder.countQuery()
	en.builder.whereQuery()
	return en.builder
}




// type IFindProductBuilder interface{
// 	openJsonQuery()
// 	initQuery()
// 	countQuery()
// 	whereQuery()
// 	sort()
// 	paginate()
// 	closeJsonQuery()
// 	resetQuery()
// 	Result() []*products.Products
// 	Count() int
// 	PrintQuery()
// }

// type findProductBuilder struct {
// 	db *sqlx.DB
// 	req *products.ProductFiler
// 	query string
// 	lastStackIndex int
// 	values []any
// }

// func FindProductBuilder(db *sqlx.DB, req *products.ProductFiler) *findProductBuilder {
// 	return &findProductBuilder{
// 		db: db,
// 		req: req,
// 	}
// }

// type findProductEngineer struct {
// 	builder IFindProductBuilder
// }

// func FindProductEngineer(builder IFindProductBuilder) *findProductEngineer {
// 	return &findProductEngineer{
// 		builder: builder,
// 	}
// }

// func (en *findProductEngineer) FindProduct() IFindProductBuilder {
// 	en.builder.openJsonQuery()
// 	en.builder.initQuery()
// 	en.builder.whereQuery()
// 	en.builder.sort()
// 	en.builder.paginate()
// 	en.builder.closeJsonQuery()
// 	return en.builder
// }

// func (en *findProductEngineer) CountProduct() IFindProductBuilder {
// 	en.builder.countQuery()
// 	en.builder.whereQuery()
// 	return en.builder
// } 


// func (b *findProductBuilder) openJsonQuery() {
// 	b.query += `
// 	SELECT
// 		array_to_json(array_agg("t"))
// 	FROM (
// 	`
	
// }
// func (b *findProductBuilder) closeJsonQuery() {
// 	b.query += `
// 	) AS "t";
// 	`
// }

// func (b *findProductBuilder) initQuery() {
// 	b.query += `
// 		SELECT
// 			"p"."id",
// 			"p"."title",
// 			"p"."description",
// 			"p"."price",
// 			(
// 				SELECT
// 					to_jsonb("ct")
// 				FROM (
// 					SELECT
// 						"c"."id",
// 						"c"."title"
// 					FROM "categories" "c"
// 						LEFT JOIN "products_categories" "pc" ON "pc"."category_id" = "c"."id"
// 					WHERE "pc"."product_id" = "p"."id"
// 				) AS "ct"
// 			) AS "category",
// 			"p"."created_at",
// 			"p"."updated_at",
// 			(
// 				SELECT
// 					COALESCE(array_to_json(array_agg("it")), '[]'::json)
// 				FROM (
// 					SELECT
// 						"i"."id",
// 						"i"."filename",
// 						"i"."url"
// 					FROM "images" "i"
// 					WHERE "i"."product_id" = "p"."id"
// 				) AS "it"
// 			) AS "images"
// 		FROM "products" "p"
// 		WHERE  1 = 1
// 	`
// }

// func (b *findProductBuilder) countQuery() {
// 	b.query += `
// 		SELECT
// 			COUNT(*) AS "count"
// 		FROM "products" "p"
// 		WHERE 1 = 1`
// }

// func (b *findProductBuilder) whereQuery() {
// 	var queryWhere string
// 	queryWhereStack := make([]string, 0)

// 	// Id check
// 	if b.req.Id != "" {
// 		b.values = append(b.values, b.req.Id)

// 		queryWhereStack = append(queryWhereStack, `
// 		AND "p"."id" = ?`)
// 	}

// 	// Search check
// 	if b.req.Search != "" {
// 		b.values = append(
// 			b.values,
// 			"%"+strings.ToLower(b.req.Search)+"%",
// 			"%"+strings.ToLower(b.req.Search)+"%",
// 		)

// 		queryWhereStack = append(queryWhereStack, `
// 		AND (LOWER("p"."title") LIKE ? OR LOWER("p"."description") LIKE ?)`)
// 	}

// 	for i := range queryWhereStack {
// 		if i != len(queryWhereStack)-1 {
// 			queryWhere += strings.Replace(queryWhereStack[i], "?", "$"+strconv.Itoa(i+1), 1)
// 		} else {
// 			//condition for last stack record
// 			queryWhere += strings.Replace(queryWhereStack[i], "?", "$"+strconv.Itoa(i+1), 1)
// 			queryWhere = strings.Replace(queryWhere, "?", "$"+strconv.Itoa(i+2), 1)
// 		}
// 	}
// 	// Last stack record
// 	b.lastStackIndex = len(b.values)

// 	fmt.Printf("queryWhere: %s\n", queryWhere)

// 	// Summary query
// 	b.query += queryWhere

// }

// // func for make it perfect : can sort multiple like sort by title and price

// func (b *findProductBuilder) sort() {
// 	orderByMap := map[string]string{
// 		"id":    `"p"."id"`,
// 		"title": `"p"."title"`,
// 		"price": `"p"."price"`,
// 	}
// 	orderBy := strings.ToLower(b.req.OrderBy)
// 	fmt.Printf("b.req.OrderBy before: %s\n", orderBy)
// 	if orderByMap[orderBy] == "" {
// 		b.req.OrderBy = orderByMap["title"]
// 	} else {
// 		b.req.OrderBy = orderByMap[orderBy]
// 	}
// 		fmt.Printf("orderByMap : %s\n", orderByMap[orderBy])
// 		fmt.Printf("b.req.OrderBy after: %s\n", orderBy)
		
// 	sortMap := map[string]string{
// 		"DESC": "DESC",
// 		"ASC":  "ASC",
// 	}
// 	if sortMap[b.req.Sort] == "" {
// 		b.req.Sort = sortMap["ASC"]
// 	} else {
// 		b.req.Sort = sortMap[strings.ToUpper(b.req.Sort)]
// 	}

// 	b.values = append(b.values, orderBy)
// 	b.query += fmt.Sprintf(`
// 		ORDER BY $%d %s`, b.lastStackIndex+1, b.req.Sort)
// 	b.lastStackIndex = len(b.values)
// }

// func (b *findProductBuilder) paginate() {
// 	// offset = (page - 1) * limit
// 	b.values = append(b.values, (b.req.Page-1)*b.req.Limit, b.req.Limit)
// 	b.query += fmt.Sprintf(`
// 		OFFSET $%d
// 		LIMIT $%d`, b.lastStackIndex+1, b.lastStackIndex+2)
// 	b.lastStackIndex = len(b.values)
// }


// func (b *findProductBuilder) resetQuery() {
// 	b.query = ""
// 	b.lastStackIndex = 0
// 	b.values = make([]any, 0)
// }

// func (b *findProductBuilder) Result() []*products.Products {
// 	_, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 	defer cancel()

// 	bytes := make([]byte, 0)
// 	productsData := make([]*products.Products, 0)

// 	if err := b.db.Get(&bytes, b.query, b.values...); err != nil {
// 		log.Printf("find product failed: %v\n", err)
// 		return make([]*products.Products, 0)
// 	}

// 	if err := json.Unmarshal(bytes, &productsData); err != nil {
// 		log.Printf("unmarshal product failed: %v\n", err)
// 		return make([]*products.Products, 0)
// 	}

// 	b.resetQuery()
// 	return productsData
	
// }

// func (b *findProductBuilder) Count() int {
// 	_, cancel := context.WithTimeout(context.Background(), 15*time.Second)
// 	defer cancel()
// 	var count int
// 	if err := b.db.Get(&count, b.query, b.values...); err != nil {
// 		log.Printf("count product failed: %v\n", err)
// 		return 0
// 	}
// 	b.resetQuery()

// 	return count
// }

// func (b *findProductBuilder) PrintQuery() {
// 	utils.Debug(b.values)
// 	fmt.Println(b.query)
// }
	