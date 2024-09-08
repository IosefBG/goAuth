package database

import "time"

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

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

// User model
type User struct {
	ID            int       `json:"id"`
	Username      string    `json:"username"`
	Password      string    `json:"password"`
	Email         string    `json:"email"`
	IsBlocked     bool      `json:"is_blocked"`
	LoginAttempts int       `json:"login_attempts"`
	LastLogin     time.Time `json:"last_login"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Addresses     []Address `json:"addresses"`
	Products      []Product `json:"products"`
	Cart          Cart      `json:"cart"`
}

// Product model
type Product struct {
	ID          int       `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	CategoryID  int       `json:"category_id"`
	UserID      int       `json:"user_id"`
}

// Category model
type Category struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Address model
type Address struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Street    string    `json:"street"`
	City      string    `json:"city"`
	State     string    `json:"state"`
	ZipCode   string    `json:"zip_code"`
	Country   string    `json:"country"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Cart model
type Cart struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Items     []CartItem `json:"items"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

// CartItem model
type CartItem struct {
	ID        int       `json:"id"`
	CartID    int       `json:"cart_id"`
	ProductID int       `json:"product_id"`
	Quantity  int       `json:"quantity"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
