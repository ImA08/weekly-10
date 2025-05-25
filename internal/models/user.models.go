package models

type UserStruct struct {
	ID       int    `json:"id" form:"-" DB:"id"`
	Email    string `json:"email" form:"email" DB:"email"`
	Role     string `json:"role" form:"-" DB:"role"`
	Password string `json:"-" form:"password" DB:"password"`
}

type AuthForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
