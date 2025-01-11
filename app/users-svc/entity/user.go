package entity

import "time"

type User struct {
	ID             uint64     `json:"id"`
	UserLogin      string     `json:"user_login"`
	UserPass       string     `json:"user_pass,omitempty"` // omitempty to exclude password from JSON by default
	UserNicename   string     `json:"user_nicename"`
	UserEmail      string     `json:"user_email"`
	UserURL        string     `json:"user_url"`
	UserRegistered time.Time  `json:"user_registered"`
	ActivationKey  string     `json:"activation_key"`
	UserStatus     int        `json:"user_status"`
	DisplayName    string     `json:"display_name"`
	FCMToken       *string    `json:"fcm_token,omitempty"`
	LastLoginAt    *time.Time `json:"last_login_at,omitempty"`
	CreatedAt      *time.Time `json:"created_at,omitempty"`
	UpdatedAt      *time.Time `json:"updated_at,omitempty"`
	DeletedAt      *time.Time `json:"deleted_at,omitempty"`
}

// Optional helper method for creating a new User
func NewUser() *User {
	return &User{
		UserRegistered: time.Now(),
		UserStatus:     0,
	}
}
