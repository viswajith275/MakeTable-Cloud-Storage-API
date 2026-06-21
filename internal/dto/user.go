package dto

import "time"

type UserRegisterRequest struct {
	Email    string `json:"email" binding:"required,email"`
	UserName string `json:"username" binding:"required,min=3,max=10"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type UserLoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8,max=20"`
}

type UserResponse struct {
	Email     string    `json:"email"`
	Username  string    `json:"username"`
	Disabled  bool      `json:"disabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
