package model

import "time"

type Profile struct {
	User_id    int       `db:"user_id"`
	FullName   string    `db:"full_name"`
	Photo      string    `db:"photo"`
	Phone      string    `db:"phone"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
}
