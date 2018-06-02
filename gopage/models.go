package gopage

type User struct {
	Id       int
	Username string
	Hash     string
	Token    Token
	Validity string
}
