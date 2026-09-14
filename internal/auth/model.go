package auth

import "time"

type User struct {
	ID         int64
	Email      string
	Password   string
	CreatedAt  time.Time
	ModifiedAt time.Time
}

type SignupRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type AuthResponse struct {
	Token string `json:"token"`
}
