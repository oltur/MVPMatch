package model

import (
	"errors"
	"github.com/oltur/mvp-match/tools"
	"github.com/oltur/mvp-match/types"
	"github.com/rs/xid"
	"time"
)

const (
	UserRoleSeller = "seller"
	UserRoleBuyer  = "buyer"
	UserRoleAdmin  = "admin"
)

var allowedUserRoles = map[string]*struct{}{UserRoleSeller: nil, UserRoleBuyer: nil, UserRoleAdmin: nil}

type User struct {
	ID           types.Id `json:"id" example:"xxx"`
	UserName     string   `json:"userName" example:"user_name"`
	PasswordHash string   `json:"passwordHash"`
	Deposit      int      `json:"deposit" example:"5"`
	Role         string   `json:"role"`
	Token        string   `json:"token"`
	TokenExpires int64    `json:"tokenExpires"`
}

func UsersAll(q string) (res []*User, err error) {
	allUsers := GetMapValuesForUsers(usersByIds)
	if q == "" {
		res = allUsers
		return
	}
	res = []*User{}
	for k, v := range allUsers {
		if q == v.UserName {
			res = append(res, allUsers[k])
		}
	}
	return
}

func UserOne(id types.Id) (res *User, err error) {
	for k := range usersByIds {
		if id == k {
			res = usersByIds[k]
			return
		}
	}
	return nil, ErrNotFound
}

func UserInsert(req *User) (res *User, err error) {
	req.ID = types.Id(xid.New().String())

	_, err = UserOne(req.ID)
	if err != nil {
		if !errors.Is(err, ErrNotFound) {
			return
		}
	} else {
		err = ErrUserIdExists
		return
	}

	isFree, err := IsUserNameFree(req.UserName)
	if err != nil {
		return
	}
	if !isFree {
		err = ErrUserNameExists
		return
	}

	usersByIds[req.ID] = req
	res = req
	return
}

func UserResetDeposit(id types.Id) (err error) {
	user := usersByIds[id]

	user.Deposit = 0

	err = UserSave(user)
	if err != nil {
		return
	}
	return
}

// UserUpdate part of CRUD
func UserUpdate(req *UpdateUserRequest) (err error) {
	user := usersByIds[req.ID]

	user.PasswordHash = tools.Hash(req.Password)

	err = UserSave(user)
	if err != nil {
		return
	}
	return
}

// UserSave Internal use only
func UserSave(req *User) (err error) {
	usersByIds[req.ID] = req
	return
}

func UserDelete(id types.Id) (err error) {
	if _, ok := usersByIds[id]; ok {
		err = ErrNotFound
		return
	}
	delete(usersByIds, id)
	return
}

var usersByIds map[types.Id]*User

func GetMapValuesForUsers(m map[types.Id]*User) (res []*User) {
	res = make([]*User, len(m))
	i := 0
	for _, v := range m {
		res[i] = v
		i++
	}
	return res
}

func IsUserNameFree(userName string) (res bool, err error) {
	for k := range usersByIds {
		if usersByIds[k].UserName == userName {
			res = false
			return
		}
	}
	res = true
	return
}

func GetUserByCredentials(userName string, password string) (res *User, err error) {
	passwordHash := tools.Hash(password)
	for k := range usersByIds {
		if usersByIds[k].UserName == userName && usersByIds[k].PasswordHash == passwordHash {
			res = usersByIds[k]
			return
		}
	}
	err = ErrNotFound
	return
}

func UserLogout(id types.Id) (err error) {
	for k := range usersByIds {
		if usersByIds[k].ID == id {
			usersByIds[k].Token = ""
			usersByIds[k].TokenExpires = 0
			return
		}
	}
	err = ErrNotFound
	return
}

func VerifyToken(userId string, token string, expires int64) (res bool, err error) {
	if expires < time.Now().UnixMilli() {
		res = false
		return
	}
	for k := range usersByIds {
		if usersByIds[k].ID == types.Id(userId) && usersByIds[k].Token == token && usersByIds[k].TokenExpires == expires {
			res = true
			return
		}
	}
	err = ErrNotFound
	return
}
