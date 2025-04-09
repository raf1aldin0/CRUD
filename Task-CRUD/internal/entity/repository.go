package entity

import (
	"time"
)

// Repository represents a code repository linked to a user.
type Repository struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`                                           // Primary key
	Name      string    `gorm:"type:varchar(100);not null" json:"name"`                                       // Repository name
	UserID    uint      `gorm:"not null" json:"user_id"`                                                      // Foreign key to User
	User      User      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user"` // Join with users
	URL       string    `gorm:"type:varchar(255);not null" json:"url"`                                        // Repository URL
	AIEnabled bool      `gorm:"default:false" json:"ai_enabled"`                                              // AI feature flag
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`                                             // Creation timestamp
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`                                             // Last update timestamp
}

// TableName explicitly sets the table name to "repositories"
func (Repository) TableName() string {
	return "repositories"
}
