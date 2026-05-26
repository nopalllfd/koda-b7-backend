package dto

import "time"

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email" example:"naufal@test.com"`
	Password string `json:"password" binding:"required,min=8" example:"12345678"`
}

type RegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginResponse struct {
	DisplayName string `json:"display_name"`
	Photo       string `json:"photo"`
	IsPinExists bool   `json:"isPinExists"`
	Token       string `json:"token"`
}

type AddPinRequest struct {
	UserID int    `json:"user_id"`
	Pin    string `json:"pin" binding:"required,min=6"`
}

type UserPIN struct {
	Pin string `json:"pin"`
}

type ChangePasswordRequest struct {
	Id          int
	OldPassword string `json:"old_password" binding:"required,min=8"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type LoginSwaggerResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    LoginResponse `json:"data"`
}

type RegisterSwaggerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

type ErrorSwaggerResponse struct {
	Success bool   `json:"success" example:"false"`
	Message string `json:"message" example:"bad request"`
	Error   string `json:"error,omitempty" example:"validation error"`
}

type LogoutRequest struct {
	Token     string
	ExpiredAt time.Time
}

type LogoutSwaggerResponse struct {
	Success bool   `json:"success"`
	Message string `json:"logout success"`
}
