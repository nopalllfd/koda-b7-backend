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
	Email string `json:"email"`
	Token string `json:"token"`
}

type AddPinRequest struct {
	UserID int    `json:"user_id"`
	Pin    string `json:"pin"`
}
