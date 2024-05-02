package domain

import "time"

// BaseModel can be used to embed in other models, id must be string(UUID,ULID)
type BaseModel struct {
	ID        string     `gorm:"size:50" json:"id"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"index"`
}

// BaseAutoIncModel can be used to embed in other models, id is integer
type BaseAutoIncModel struct {
	ID        uint       `json:"id"`
	CreatedAt time.Time  `gorm:"not null" json:"created_at"`
	UpdatedAt time.Time  `gorm:"not null" json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}
