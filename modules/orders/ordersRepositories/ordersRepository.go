package ordersRepositories

import (
	"encoding/json"
	"fmt"

	"github.com/NatthawutSK/ri-shop/modules/orders"
	"github.com/NatthawutSK/ri-shop/modules/orders/ordersPattern"
	"github.com/jmoiron/sqlx"
)

type IOrdersRepository interface{
	FindOneOrder(orderId string) (*orders.Order, error)
	FindOrder(req *orders.OrderFilter) ([]*orders.Order, int)
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

	if err := r.db.Get(&bytes, query, orderId) ; err != nil {
		return nil, fmt.Errorf("cannot get order: %w", err)
	}

	if err := json.Unmarshal(bytes, &order); err != nil {
		return nil, fmt.Errorf("unmarshal order failed: %v", err)
	}
	
	return order, nil
}

func (r *ordersRepository) FindOrder(req *orders.OrderFilter) ([]*orders.Order, int){
	builder := ordersPattern.FindOrderBuilder(r.db, req)
	engineer := ordersPattern.FindOrderEngineer(builder)

	return engineer.FindOrder(), engineer.CountOrder()
}
