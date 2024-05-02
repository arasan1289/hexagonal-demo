package domain

// UserRole defines user roles
type UserRole string

const (
	Admin    UserRole = "ADM"
	Rider    UserRole = "RID"
	Customer UserRole = "CUS"
)

// User represents a user in the system
type User struct {
	BaseModel
	FirstName             string   `gorm:"size:50;not null" json:"first_name"`
	LastName              string   `gorm:"size:50;not null" json:"last_name"`
	Email                 *string  `gorm:"-:all" json:"email"`
	EmailHash             *string  `gorm:"size:256;index" json:"email_hash"`
	EmailEncrypted        *string  `json:"email_encrypted"`
	IsEmailVerified       bool     `gorm:"default:false" json:"is_email_verified"`
	PhoneNumber           string   `gorm:"-:all" json:"phone_number"`
	PhoneNumberHash       string   `gorm:"size:256;index" json:"phone_number_hash"`
	PhoneNumberEncrypted  string   `json:"phone_number_encrypted"`
	IsPhoneNumberVerified bool     `gorm:"default:false" json:"is_phone_number_verified"`
	Role                  UserRole `gorm:"size:5;not null" json:"user_role"`
	IsActive              bool     `gorm:"default:false" json:"is_active"`
}
