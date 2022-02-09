package model

type BuyResult struct {
	ProductName string  `json:"productName" example:"product_name"`
	Change      []*Coin `json:"change"`
	Total       int     `json:"total" example:"5"`
}
