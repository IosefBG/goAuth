package models

import "time"

type RegistrationRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
}

type LoginRequest struct {
	Identifier string `json:"identifier" binding:"required"`
	Password   string `json:"password" binding:"required"`
}

type UserData struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
}

type AuthResponse struct {
	Token string   `json:"token"`
	User  UserData `json:"user"`
}

type SessionResponse struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	Token           string    `json:"session_token"`
	IPAddress       string    `json:"ip_address"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_date_at"`
	UpdatedAt       time.Time `json:"updated_date_at"`
	Location        string    `json:"location"`
	DeviceConnected string    `json:"device_connected"`
	BrowserUsed     string    `json:"browser_used"`
}
