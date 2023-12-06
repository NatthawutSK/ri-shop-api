package entities

type PaginationReq struct {
	Page      int `json:"page" query:"page"`
	Limit     int `json:"limit" query:"limit"`
	TotalPage int `json:"total_page" query:"total_page"`
	TotalItem int `json:"total_item" query:"total_item"`
}

type SortReq struct {
	OrderBy string `json:"order_by" query:"order_by"`
	Sort    string `json:"sort" query:"sort"` // asc or desc
}
