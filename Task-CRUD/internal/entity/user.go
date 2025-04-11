package entity

import "time"

// User represents the user entity stored in the database.
type User struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`             // Primary key, auto increment
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`         // User's full name
	Email     string    `gorm:"type:varchar(100);unique;not null" json:"email"` // Unique email address
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`               // Created timestamp
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`               // Updated timestamp
}
