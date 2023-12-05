package orders

import (
	"github.com/NatthawutSK/ri-shop/modules/entities"
	"github.com/NatthawutSK/ri-shop/modules/products"
)

type Order struct {
	Id           string           `json:"id" db:"id"`
	UserId       string           `json:"user_id" db:"user_id"`
	TransferSlip *TransferSlip    `json:"transfer_slip" db:"transfer_slip"`
	Products     []*ProductsOrder `json:"products"`
	Address 	string           `json:"address" db:"address"`
	Contact 	string           `json:"contact" db:"contact"`
	Status 		string           `json:"status" db:"status"`
	TotalPaid 	float64          `json:"total_paid" db:"total_paid"`
	CreatedAt    string           `json:"created_at" db:"created_at"`
	UpdatedAt   string           `json:"updated_at" db:"updated_at"`
}

type TransferSlip struct {
	Id        string `json:"id"`
	FileName  string `json:"file_name"`
	Url       string `json:"url"`
	CreatedAt string `json:"created_at"`
}

type ProductsOrder struct {
	Id      string `json:"id" db:"id"`
	Qty     int    `json:"qty" db:"qty"`
	Product *products.Products `json:"product" db:"product"`
}


type OrderFilter struct {
	Search    string `query:"search"` // user_id, address, contact
	Status    string `query:"status"`
	StartDate string `query:"start_date"`
	EndDate   string `query:"end_date"`
	*entities.PaginationReq
	*entities.SortReq
}
