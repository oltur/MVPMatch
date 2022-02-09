package model

import (
	"github.com/oltur/mvp-match/types"
)

func init() {
	var id types.Id

	usersByIds = make(map[types.Id]*User)

	id = "1" // types.Id(xid.New().String())
	user1 := &User{
		ID:           id,
		UserName:     "User #1, Seller",
		PasswordHash: "6b86b273ff34fce19d6b804eff5a3f5747ada4eaa22f1d49c01e52ddb7875b4b", // 1
		Role:         UserRoleSeller,
	}
	usersByIds[id] = user1
	id = "2" // types.Id(xid.New().String())
	user2 := &User{
		ID:           id,
		UserName:     "User #2, Seller",
		PasswordHash: "d4735e3a265e16eee03f59718b9b5d03019c07d8b6c51f90da3a666eec13ab35", // 2
		Role:         UserRoleSeller,
	}
	usersByIds[id] = user2
	id = "3" // types.Id(xid.New().String())
	user3 := &User{
		ID:           id,
		UserName:     "User #3, Buyer",
		PasswordHash: "4e07408562bedb8b60ce05c1decfe3ad16b72230967de01f640b7e4729b49fce", // 3
		Role:         UserRoleBuyer,
	}
	usersByIds[id] = user3
	id = "4" // types.Id(xid.New().String())
	user4 := &User{
		ID:           id,
		UserName:     "User #4, Admin",
		PasswordHash: "4b227777d4dd1fc61c6f884f48641d02b4d121d3fd328cb08b5531fcacdabf8a", // 4
		Role:         UserRoleAdmin,
	}
	usersByIds[id] = user4

	productsByIds = make(map[types.Id]*Product)

	id = "1" // types.Id(xid.New().String())
	product1 := &Product{
		ID:              id,
		ProductName:     "Product #1",
		SellerId:        user1.ID,
		AmountAvailable: 1000,
		Cost:            20,
	}
	productsByIds[id] = product1

	id = "2" // types.Id(xid.New().String())
	product2 := &Product{
		ID:              id,
		ProductName:     "Product #2",
		SellerId:        user1.ID,
		AmountAvailable: 1,
		Cost:            30,
	}
	productsByIds[id] = product2

	id = "3" // types.Id(xid.New().String())
	product3 := &Product{
		ID:              id,
		ProductName:     "Product #3",
		SellerId:        user2.ID,
		AmountAvailable: 3000,
		Cost:            40,
	}
	productsByIds[id] = product3
}
