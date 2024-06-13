package domain

import (
	"errors"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/arasan1289/hexagonal-demo/internal/core/util"
	"gorm.io/gorm"
)

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

func (u *User) AfterFind(tx *gorm.DB) (err error) {
	// Fetch configuration from GORM instance
	confValue, _ := tx.InstanceGet("config")
	if confValue == nil {
		// Handle nil value
		return errors.New("config not found in GORM instance")
	}

	conf, ok := confValue.(*config.Container)
	if !ok {
		// Handle incorrect type
		return errors.New("config has unexpected type")
	}
	if u.PhoneNumberEncrypted != "" {
		p, err := util.DecryptString(u.PhoneNumberEncrypted, conf.App.SecretKey)
		if err != nil {
			return err
		}
		u.PhoneNumber = p
	}
	if u.EmailEncrypted != nil {
		e, err := util.DecryptString(*u.EmailEncrypted, conf.App.SecretKey)
		if err != nil {
			return err
		}
		u.Email = &e
	}
	return nil
}
