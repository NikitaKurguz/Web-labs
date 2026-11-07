package model

type User struct {
	Id       int
	Login    string
	Password string
	FullName string
	Email    string
	Avatar   string
	Hash     string
	IsAdmin  bool
}
