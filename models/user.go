package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

	"go-grst-boilerplate/app/users-svc/entity"
)

type User struct {
	ID             uint64          `gorm:"primaryKey;column:ID;auto_increment"`
	UserLogin      string          `gorm:"unique;column:user_login;type:varchar(60);not null;default:''"`
	UserPass       string          `gorm:"column:user_pass;type:varchar(255);not null;default:''"`
	UserNicename   string          `gorm:"column:user_nicename;type:varchar(250);not null;default:''"`
	UserEmail      string          `gorm:"column:user_email;type:varchar(100);not null;default:''"`
	UserURL        string          `gorm:"column:user_url;type:varchar(100);not null;default:''"`
	UserRegistered time.Time       `gorm:"column:user_registered;not null;default:CURRENT_TIMESTAMP"`
	ActivationKey  string          `gorm:"column:user_activation_key;type:varchar(255);not null;default:''"`
	UserStatus     int             `gorm:"column:user_status;not null;default:0"`
	DisplayName    string          `gorm:"column:display_name;type:varchar(250);not null;default:''"`
	FCMToken       *string         `gorm:"column:fcm_token;type:varchar(255)"`
	LastLoginAt    *time.Time      `gorm:"column:last_login_at"`
	CreatedAt      *time.Time      `gorm:"column:created_at"`
	UpdatedAt      *time.Time      `gorm:"column:updated_at"`
	DeletedAt      *gorm.DeletedAt `gorm:"column:deleted_at"`
}

// TableName specifies the table name for the User model
func (User) TableName() string {
	return "sfwc_users"
}

// ToEntity converts the User model to a User entity
func (user *User) ToEntity() *entity.User {
	return &entity.User{
		ID:             user.ID,
		UserLogin:      user.UserLogin,
		UserPass:       user.UserPass,
		UserNicename:   user.UserNicename,
		UserEmail:      user.UserEmail,
		UserURL:        user.UserURL,
		UserRegistered: user.UserRegistered.Local(),
		ActivationKey:  user.ActivationKey,
		UserStatus:     user.UserStatus,
		DisplayName:    user.DisplayName,
		FCMToken:       user.FCMToken,
		LastLoginAt:    user.LastLoginAt,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

// FromEntity creates a User model from a User entity
func (User) FromEntity(user *entity.User) *User {
	return &User{
		ID:             user.ID,
		UserLogin:      user.UserLogin,
		UserPass:       user.UserPass,
		UserNicename:   user.UserNicename,
		UserEmail:      user.UserEmail,
		UserURL:        user.UserURL,
		UserRegistered: user.UserRegistered,
		ActivationKey:  user.ActivationKey,
		UserStatus:     user.UserStatus,
		DisplayName:    user.DisplayName,
		FCMToken:       user.FCMToken,
		LastLoginAt:    user.LastLoginAt,
		CreatedAt:      user.CreatedAt,
		UpdatedAt:      user.UpdatedAt,
	}
}

// HashPassword hashes the provided password using bcrypt
func (user *User) HashPassword(password string) error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.UserPass = string(bytes)
	return nil
}

// CheckPassword verifies if the provided password matches the hashed password
func (user *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(user.UserPass), []byte(providedPassword))
	if err != nil {
		return err
	}
	return nil
}
