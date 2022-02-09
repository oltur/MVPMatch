package model

type LoginResult struct {
	Token        string `json:"token"`
	TokenExpires int64  `json:"tokenExpires"`
}
