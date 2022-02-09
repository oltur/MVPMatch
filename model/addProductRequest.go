package model

type AddProductReq struct {
	ProductName     string `json:"productName" example:"product_name"`
	AmountAvailable int    `json:"amountAvailable" example:"1"`
	Cost            int    `json:"cost" example:"5"`
}

func (a AddProductReq) Validation() (err error) {
	if a.Cost%5 != 0 {
		err = ErrInvalidCost
		return
	}
	return
}
