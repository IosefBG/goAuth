package entities

import "time"

type Session struct {
	ID              int       `json:"id"`
	UserID          int       `json:"user_id"`
	IPAddress       string    `json:"ip_address"`
	Location        string    `json:"location"`
	DeviceConnected string    `json:"device_connected"`
	BrowserUsed     string    `json:"browser_used"`
	IsActive        bool      `json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
