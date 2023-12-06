package ordersRepositories

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/NatthawutSK/ri-shop/modules/orders"
	"github.com/NatthawutSK/ri-shop/modules/orders/ordersPattern"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepository interface {
	FindOneOrder(orderId string) (*orders.Order, error)
	FindOrder(req *orders.OrderFilter) ([]*orders.Order, int)
	InsertOrder(req *orders.Order) (string, error)
	UpdateOrder(req *orders.OrderUpdate) error
}

type ordersRepository struct {
	db *sqlx.DB
}

func OrdersRepository(db *sqlx.DB) IOrdersRepository {
	return &ordersRepository{
		db: db,
	}
}

func (r *ordersRepository) FindOneOrder(orderId string) (*orders.Order, error) {
	query := `
	SELECT
		to_jsonb("t")
	FROM (
		SELECT
			"o"."id",
			"o"."user_id",
			"o"."transfer_slip",
			"o"."status",
			(
				SELECT
					array_to_json(array_agg("pt"))
				FROM (
					SELECT
						"spo"."id",
						"spo"."qty",
						"spo"."product"
					FROM "products_orders" "spo"
					WHERE "spo"."order_id" = "o"."id"
				) AS "pt"
			) AS "products",
			"o"."address",
			"o"."contact",
			(
				SELECT
					SUM(COALESCE(("po"."product"->>'price')::FLOAT*("po"."qty")::FLOAT, 0))
				FROM "products_orders" "po"
				WHERE "po"."order_id" = "o"."id"
			) AS "total_paid",
			"o"."created_at",
			"o"."updated_at"
		FROM "orders" "o"
		WHERE "o"."id" = $1
	) AS "t";`

	//SUM(COALESCE(("po"."product"->>'price')::FLOAT*("po"."qty")::FLOAT, 0))
	//  "po"."product"->>'price' คือการเข้าถึง value ของ key "price" ใน jsonb ของ column "product" เพราะ product เก็บเป็น json
	// ::FLOAT คือ casting ค่าให้เป็น float
	// COALESCE คือ ถ้าเป็น null ให้ return เป็น 0 แทน

	bytes := make([]byte, 0)
	order := &orders.Order{
		Products: make([]*orders.ProductsOrder, 0),
	}

	if err := r.db.Get(&bytes, query, orderId); err != nil {
		return nil, fmt.Errorf("cannot get order: %w", err)
	}

	if err := json.Unmarshal(bytes, &order); err != nil {
		return nil, fmt.Errorf("unmarshal order failed: %v", err)
	}

	return order, nil
}

func (r *ordersRepository) FindOrder(req *orders.OrderFilter) ([]*orders.Order, int) {
	builder := ordersPattern.FindOrderBuilder(r.db, req)
	engineer := ordersPattern.FindOrderEngineer(builder)

	return engineer.FindOrder(), engineer.CountOrder()
}

func (r *ordersRepository) InsertOrder(req *orders.Order) (string, error) {
	builder := ordersPattern.InsertOrderBuilder(req, r.db)
	orderId, err := ordersPattern.InsertOrderEngineer(builder).InsertOrder()
	if err != nil {
		return "", err
	}

	return orderId, nil
}

// func (r *ordersRepository) UpdateOrder(req *orders.OrderUpdate) error {
// 	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
// 	defer cancel()

// 	query := `
// 	UPDATE "orders" SET`

// 	values := make([]any, 0)
// 	lastIndex := 1
// 	whereStackQuery := make([]string, 0)

// 	if req.Status != "" {
// 		values = append(values, req.Status)
// 		whereStackQuery = append(whereStackQuery, fmt.Sprintf(`
// 		"status" = $%d?`, lastIndex)) // ? มีไว้เพื่อ replace เป็น , ในกรณีที่ไม่ใช่ตัวสุดท้าย
// 		lastIndex++
// 	}

// 	if req.TransferSlip != nil {
// 		values = append(values, req.TransferSlip)
// 		whereStackQuery = append(whereStackQuery, fmt.Sprintf(`
// 		"transfer_slip" = $%d?`, lastIndex))
// 		lastIndex++
// 	}

// 	values = append(values, req.Id)

// 	for i := range whereStackQuery {
// 		if i != len(whereStackQuery)-1 {
// 			query += strings.Replace(whereStackQuery[i], "?", ",", 1)
// 		} else {
// 			query += strings.Replace(whereStackQuery[i], "?", "", 1)
// 		}
// 	}

// 	queryClose := fmt.Sprintf(`WHERE "id" = $%d;`, lastIndex)

// 	query += queryClose

// 	if _, err := r.db.ExecContext(ctx, query, values...); err != nil {
// 		return fmt.Errorf("cannot update order: %w", err)
// 	}

// 	return nil
// }

func (r *ordersRepository) UpdateOrder(req *orders.OrderUpdate) error {
	query := `
	UPDATE "orders" SET`

	queryWhereStack := make([]string, 0)
	values := make([]any, 0)
	lastIndex := 1

	if req.Status != "" {
		values = append(values, req.Status)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"status" = $%d?`, lastIndex))

		lastIndex++
	}

	if req.TransferSlip != nil {
		values = append(values, req.TransferSlip)

		queryWhereStack = append(queryWhereStack, fmt.Sprintf(`
		"transfer_slip" = $%d?`, lastIndex))

		lastIndex++
	}

	values = append(values, req.Id)

	queryClose := fmt.Sprintf(`
	WHERE "id" = $%d;`, lastIndex)

	for i := range queryWhereStack {
		if i != len(queryWhereStack)-1 {
			query += strings.Replace(queryWhereStack[i], "?", ",", 1)
		} else {
			query += strings.Replace(queryWhereStack[i], "?", "", 1)
		}
	}
	query += queryClose

	if _, err := r.db.ExecContext(context.Background(), query, values...); err != nil {
		return fmt.Errorf("update order failed: %v", err)
	}
	return nil
}
