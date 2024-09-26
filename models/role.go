package models

import "time"

type Role struct {
	ID          int           `gorm:"primary_key" json:"id"`
	Name        string        `gorm:"index;size:100;not null" json:"name" binding:"required"`
	RoleModules []*RoleModule `gorm:"foreignKey:RoleId"`
	CreatedAt   time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
}

type NewRole struct {
	Name           string              `json:"name" binding:"required"`
	AllowedModules []*NewAllowedModule `json:"allowed_modules"`
}

type NewAllowedModule struct {
	ModuleID       int    `json:"moduleId"`
	AllowedActions string `json:"allowedActions"`
}