package dto

import "time"

type Profiles struct {
	User_id   int        `json:"user_id"`
	FullName  string     `json:"fullname"`
	Photo     string     `json:"photo"`
	Phone     string     `json:"phone"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type ProfileUpdateRequest struct {
	FullName string `json:"fullname"`
	Photo    string `json:"photo"`
	Phone    string `json:"phone"`
}
