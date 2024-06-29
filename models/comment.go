package models

import (
    "time"
    "github.com/jinzhu/gorm"
)

type Comment struct {
    gorm.Model
    ID       uint       `gorm:"primaryKey"`
    Message  string     `gorm:"not null" json:"message"`
    UserID   uint       `gorm:"not null" json:"user_id"`
    PhotoID  uint       `gorm:"not null" json:"photo_id"`
    CreatedAt *time.Time `json:"created_at,omitempty"`
    UpdatedAt *time.Time `json:"updated_at,omitempty"`

    User  User  `gorm:"foreignKey:UserID"`
    Photo Photo `gorm:"foreignKey:PhotoID"`
}
