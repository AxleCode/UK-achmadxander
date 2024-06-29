package models

import "github.com/jinzhu/gorm"

type SocialMedia struct {
    gorm.Model
    ID           uint   `gorm:"primaryKey"`
    Name         string `gorm:"not null" json:"name" validate:"required"`
    SocialMediaURL string `gorm:"not null" json:"social_media_url" validate:"required,url"`
    UserID       uint   `gorm:"not null" json:"user_id" gorm:"foreignKey:UserID"` 
}
