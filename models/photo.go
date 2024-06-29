package models

import (
    "time"
    "github.com/jinzhu/gorm"
)

type Photo struct {
    gorm.Model
    ID        uint       `gorm:"primaryKey"`
    Title     string     `gorm:"not null" json:"title"`
    Caption   string     `json:"caption"`
    URL       string     `gorm:"not null" json:"photo_url" validate:"required,url"`
    UserID    uint       `gorm:"not null" json:"user_id"`
    CreatedAt *time.Time `json:"created_at,omitempty"`
    UpdatedAt *time.Time `json:"updated_at,omitempty"`
    Comments  []Comment  `gorm:"foreignKey:PhotoID"`
    User      User       `gorm:"foreignKey:UserID"`
}

