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
	ID        int
	UserID    int
	Token     string
	IPAddress string
	Location  string
	CreatedAt time.Time
	IsValid   bool
}

type User struct {
	ID       int
	Username string
	Password string
	Email    string
}
