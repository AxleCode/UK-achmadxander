package models

import (
    "time"
    "github.com/jinzhu/gorm"
)

type User struct {
    gorm.Model
    ID          uint        `gorm:"primaryKey" json:"userID"`
    Username    string      `gorm:"not null;uniqueIndex" json:"username"`
    Email       string      `gorm:"not null;uniqueIndex" json:"email" form:"email" valid:"required~Your email is required,email~Invalid email format"`
    Password    string      `gorm:"not null" json:"password" form:"password" valid:"required~Your password is required,minstringlength(6)~Password has to have a minimum length of 6 characters"`
    Age         uint        `gorm:"not null" json:"age" validate:"min=9"`
    CreatedAt	*time.Time	`json:"created_at,omitempty"`
	UpdatedAt	*time.Time	`json:"updated_at,omitempty"`
    Photos      []Photo     `gorm:"foreignKey:UserID"`
    Comments    []Comment
    SocialMedias []SocialMedia
}
