package models

type User struct {
	FirstName string `json:"firstname"`
	SurName   string `json:"surname"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Age       int    `json:"age"`
}
