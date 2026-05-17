package model

import "time"

type User struct {
	Id         int       `db:"id"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	Pin        string    `db:"pin"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at_at"`
}
