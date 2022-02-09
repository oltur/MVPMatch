package model

import (
	"errors"
)

var (
	ErrInvalidCoinValue = errors.New("invalid coin value")
)

var allowedCoinValues = map[int]*struct{}{5: nil, 10: nil, 20: nil, 50: nil, 100: nil}

type Coin struct {
	Value int `json:"value" example:"5"`
}

// Validation example
func (a Coin) Validation() (err error) {
	if _, ok := allowedCoinValues[a.Value]; !ok {
		return ErrInvalidCoinValue
	}
	return
}
