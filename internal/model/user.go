package model

type User struct {
	GUID        string `json:"guid"`
	RefreshHash string `json:"refresh_token"`
	Email       string `json:"email"`
}
