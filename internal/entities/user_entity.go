package entities

import "time"

type User struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	Email         string    `json:"email"`
	IsBlocked     bool      `json:"is_blocked"`
	IsActive      bool      `json:"is_active"`
	LoginAttempts int       `json:"login_attempts"`
	LastLogin     time.Time `json:"last_login"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Addresses     []Address `json:"addresses"`
	Products      []Product `json:"products"`
	Cart          Cart      `json:"cart"`
}
