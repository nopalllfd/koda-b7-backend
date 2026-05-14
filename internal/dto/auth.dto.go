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
	ID    int    `json:"id"`
	Email string `json:"email"`
	Token string `json:"token"`
}
