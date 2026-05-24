package model

import "time"

type User struct {
	Id        int       `db:"id"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	Pin       string    `db:"pin"`
	FullName  string    `db:"full_name"`
	Photo     string    `db:"photo"`
	Phone     string    `db:"phone"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserPIN struct {
	Pin *string `db:"pin"`
}
