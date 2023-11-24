package products

import (
	"github.com/NatthawutSK/ri-shop/modules/appinfo"
	"github.com/NatthawutSK/ri-shop/modules/entities"
)

type Products struct {
	Id          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    *appinfo.Category `json:"category"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
	Price 	 float64 `json:"price"`
	Images []*entities.Image `json:"images"`
}