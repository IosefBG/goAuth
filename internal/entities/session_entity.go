package entities

import "time"

type Session struct {
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
