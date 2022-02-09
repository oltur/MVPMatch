package model

type BuyResponse struct {
	ProductName string  `json:"productName" example:"product_name"`
	Change      []*Coin `json:"change"`
	Total       int     `json:"total" example:"5"`
}
