package model

type User struct {
	ID           int    `json:"id"`
	User_id      string `json:"user_id"`
	Access_token string `json:"token"`
}
