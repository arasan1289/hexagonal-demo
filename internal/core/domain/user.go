package domain

import "time"

type UserRole string

const (
	Admin    UserRole = "ADM"
	Rider    UserRole = "RID"
	Customer UserRole = "CUS"
)

type User struct {
	ID                    string    `json:"id"`
	FirstName             *string   `json:"first_name"`
	LastName              *string   `json:"last_name"`
	Email                 *string   `gorm:"-:all" json:"email"`
	EmailHash             *string   `json:"email_hash"`
	EmailEncrypted        *string   `json:"email_encrypted"`
	IsEmailVerified       bool      `json:"is_email_verified"`
	PhoneNumber           string    `gorm:"-:all" json:"phone_number"`
	PhoneNumberHash       string    `json:"phone_number_hash"`
	PhoneNumberEncrypted  string    `json:"phone_number_encrypted"`
	IsPhoneNumberVerified bool      `json:"is_phone_number_verified"`
	Role                  UserRole  `json:"user_role"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
