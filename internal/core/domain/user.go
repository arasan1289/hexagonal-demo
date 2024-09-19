package domain

import (
	"database/sql/driver"
	"encoding/binary"
	"errors"
	"math"

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
	Password              *string  `gorm:"size:256"`
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

type Address struct {
	BaseModel
	UserID      string `gorm:"not null" json:"user_id"`
	Name        string `gorm:"size:100;not null" json:"name"`
	HouseNumber string `gorm:"size:100;not null" json:"house_number"`
	Floor       int    `gorm:"default:0" json:"floor"`
	Street      string `gorm:"size:100;not null" json:"street"`
	City        string `gorm:"size:100;not null" json:"city"`
	State       string `gorm:"size:100;not null" json:"state"`
	Pincode     string `gorm:"size:100;not null" json:"pincode"`
	Location    Point  `gorm:"type:geometry(Point,4326)" json:"location"`
}

type Point struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Scan implements the sql.Scanner interface for reading from the database.
func (p *Point) Scan(value interface{}) error {
	if value == nil {
		return errors.New("value is nil")
	}

	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("cannot convert to byte slice")
	}

	// Check for endianess
	var byteOrder binary.ByteOrder
	if bytes[0] == 1 { // little-endian marker
		byteOrder = binary.LittleEndian
	}
	byteOrder = binary.BigEndian

	if len(bytes) < 21 {
		return errors.New("invalid coordinates length")
	}

	// Extract longitude and latitude from bytes
	p.Longitude = math.Float64frombits(byteOrder.Uint64(bytes[5:13]))
	p.Latitude = math.Float64frombits(byteOrder.Uint64(bytes[13:21]))

	return nil
}

// Value implements the driver.Valuer interface for writing to the database.
func (p Point) Value() (driver.Value, error) {
	buf := make([]byte, 21)
	buf[0] = 1 // little-endian marker
	binary.LittleEndian.PutUint32(buf[1:5], 4326)
	binary.LittleEndian.PutUint64(buf[5:13], math.Float64bits(p.Longitude))
	binary.LittleEndian.PutUint64(buf[13:21], math.Float64bits(p.Latitude))

	return buf, nil
}
