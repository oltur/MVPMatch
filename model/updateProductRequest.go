package model

import "github.com/oltur/mvp-match/types"

type UpdateProductRequest struct {
	ID              types.Id `json:"id" example:"xxx"`
	ProductName     string   `json:"productName" example:"product_name"`
	AmountAvailable int      `json:"amountAvailable" example:"1"`
	Cost            int      `json:"cost" example:"5"`
}

func (a UpdateProductRequest) Validation() (err error) {
	if a.ID == "" {
		err = ErrInvalidID
		return
	}
	if a.ProductName == "" {
		err = ErrInvalidProductName
		return
	}
	if a.AmountAvailable < 0 {
		err = ErrInvalidAmount
		return
	}
	if a.Cost%5 != 0 || a.Cost <= 0 {
		err = ErrInvalidCost
		return
	}
	return
}
