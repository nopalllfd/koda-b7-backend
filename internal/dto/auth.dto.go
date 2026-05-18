package dto

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
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
