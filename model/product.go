package model

import (
	"errors"
	"github.com/oltur/mvp-match/types"
	"github.com/rs/xid"
)

type Product struct {
	ID              types.Id `json:"id" example:"xxx"`
	ProductName     string   `json:"productName" example:"product_name"`
	SellerId        types.Id `json:"sellerId"`
	AmountAvailable int      `json:"amountAvailable" example:"1"`
	Cost            int      `json:"cost" example:"5"`
}

func ProductsAll(q string) (res []*Product, err error) {
	allProducts := GetMapValuesForProducts(productsByIds)
	if q == "" {
		res = allProducts
		return
	}
	res = []*Product{}
	for k, v := range allProducts {
		if q == v.ProductName {
			res = append(res, allProducts[k])
		}
	}
	return
}

func ProductOne(id types.Id) (res *Product, err error) {
	for k := range productsByIds {
		if id == k {
			res = productsByIds[k]
			return
		}
	}
	return nil, ErrNotFound
}

func ProductUpdate(req *UpdateProductRequest) (err error) {
	product := productsByIds[req.ID]

	product.ProductName = req.ProductName
	product.Cost = req.Cost
	product.AmountAvailable = req.AmountAvailable

	err = ProductSave(product)
	if err != nil {
		return
	}
	return
}

// ProductSave Internal use only
func ProductSave(req *Product) (err error) {
	productsByIds[req.ID] = req
	return
}

func ProductDelete(id types.Id) (err error) {
	if _, ok := productsByIds[id]; ok {
		err = ErrNotFound
		return
	}
	delete(productsByIds, id)
	return
}

var productsByIds map[types.Id]*Product

//func GetMapKeysForProducts(m map[types.Id]*Product) (res []types.Id) {
//	res = make([]types.Id, len(m))
//	i := 0
//	for k := range m {
//		res[i] = k
//		i++
//	}
//	return
//}

func GetMapValuesForProducts(m map[types.Id]*Product) (res []*Product) {
	res = make([]*Product, len(m))
	i := 0
	for _, v := range m {
		res[i] = v
		i++
	}
	return res
}

func ProductInsert(req *Product) (res *Product, err error) {
	req.ID = types.Id(xid.New().String())

	_, err = ProductOne(req.ID)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return
		} else {
			err = nil
		}
	} else {
		err = ErrProductIdExists
		return
	}

	productsByIds[req.ID] = req
	res = req
	return
}
