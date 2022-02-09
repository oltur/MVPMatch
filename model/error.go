package model

import "errors"

var (
	ErrNotFound                = errors.New("not found")
	ErrInvalidCost             = errors.New("cost should be positive and in multiples of 5")
	ErrInvalidAmount           = errors.New("amount should be non-negative")
	ErrInvalidSeller           = errors.New("user is not a seller")
	ErrInvalidBuyer            = errors.New("user is not a buyer")
	ErrNotEnoughAmount         = errors.New("not enough items to buy")
	ErrProductIdExists         = errors.New("given product id already exists")
	ErrWrongSeller             = errors.New("the current user does not own this product")
	ErrInvalidAdmin            = errors.New("user is not an admin")
	ErrInvalidProductName      = errors.New("invalid product name")
	ErrInvalidID               = errors.New("invalid id")
	ErrInvalidUserName         = errors.New("invalid user name")
	ErrInvalidPassword         = errors.New("invalid password")
	ErrInvalidUserRole         = errors.New("unsupported user role")
	ErrCannotCreateAdmin       = errors.New("cannot create admin user")
	ErrUserNotFoundInContext   = errors.New("user not found in context")
	ErrNotEnoughDeposit        = errors.New("not enough deposit")
	ErrCannotGenerateUserToken = errors.New("cannot generate user token")
	ErrCannotValidateUserToken = errors.New("cannot validate user token")
	ErrUnauthorized            = errors.New("not authorized")
	ErrUserIdExists            = errors.New("user with given ID already exists")
	ErrUserNameExists          = errors.New("user with given name already exists")
	ErrActiveSessionExists     = errors.New("there is already an active session using your account")
)
