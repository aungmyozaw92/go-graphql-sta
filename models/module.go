package models

import "time"

type Module struct {
	ID         int       `gorm:"primary_key" json:"id"`
	Name       string    `gorm:"index;size:100;not null" json:"name" binding:"required"`
	Actions    string    `gorm:"not null" json:"action" binding:"required"`
	CreatedAt  time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewModule struct {
	Name       string `json:"name" binding:"required"`
	Actions    string `json:"action" binding:"required"`
}
