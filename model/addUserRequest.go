package model

type AddUserReq struct {
	UserName string `json:"userName" example:"user_name"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

func (a AddUserReq) Validation() (err error) {
	if a.UserName == "" {
		err = ErrInvalidUserName
		return
	}
	if a.Password == "" {
		err = ErrInvalidPassword
		return
	}
	if _, ok := allowedUserRoles[a.Role]; !ok {
		err = ErrInvalidUserRole
		return
	}
	// TODO: Add separate API?
	if a.Role == UserRoleAdmin {
		err = ErrCannotCreateAdmin
		return
	}

	return
}
