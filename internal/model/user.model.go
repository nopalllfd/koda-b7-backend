package model

type User struct {
	Id       int    `db:"id"`
	Email    string `db:"employee_name"`
	Password string `db:"department_id"`
	Pin      string `db:"salary"`
}
