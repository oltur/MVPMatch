package model

import "github.com/oltur/mvp-match/types"

type UpdateUserRequest struct {
	ID       types.Id `json:"id" example:"xxx"`
	Password string   `json:"password"`
}

func (a UpdateUserRequest) Validation() (err error) {
	if a.ID == "" {
		err = ErrInvalidID
		return
	}
	if a.Password == "" {
		err = ErrInvalidPassword
		return
	}
	// TODO: Add separate API for changing other fields?

	return
}
