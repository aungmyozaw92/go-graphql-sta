package models

import "time"

type RoleModule struct {
	RoleId         int       `gorm:"primary_key;autoIncrement:false;not null" json:"role_id" binding:"required"`
	ModuleId       int       `gorm:"primary_key;autoIncrement:false;not null" json:"module_id" binding:"required"`
	AllowedActions string    `gorm:"not null" json:"allowed_actions" binding:"required"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Role           Role      `json:"role"`
	Module         Module    `json:"module"`
}

type NewRoleModule struct {
	RoleId         int    `json:"role_id" binding:"required"`
	ModuleId       int    `json:"module_id" binding:"required"`
	AllowedActions string `json:"allowed_actions" binding:"required"`
}