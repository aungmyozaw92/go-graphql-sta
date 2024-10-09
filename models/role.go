package models

import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

	"github.com/aungmyozaw92/go-graphql/config"
	"github.com/aungmyozaw92/go-graphql/utils"
)

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

func extractModuleActions(s string) []string {
	return strings.Split(strings.ToLower(s), ";")
}

// retrieve allowed query paths for role
func GetQueryPathsFromRole(ctx context.Context, roleId int) (map[string]bool, error) {
	db := config.GetDB()
	var role Role
	if err := db.WithContext(ctx).
			Preload("RoleModules").
			Preload("RoleModules.Module").
			Where("id = ?", roleId).
			First(&role).Error; err != nil {
		return nil, errors.New("role not found")
	}

	allowedPaths := make(map[string]bool, 0)
	for _, permission := range role.RoleModules {
		validActions := extractModuleActions(permission.Module.Actions)
		allowedActions := extractModuleActions(permission.AllowedActions)
		module := permission.Module.Name

		for _, action := range allowedActions {
			// check if the action is valid

			if slices.Contains(validActions, action) {
				// changing case of action & module for older module name convention
				module = utils.UppercaseFirst(module)
				switch action {
				case "read":
					allowedPaths["get"+module] = true
					allowedPaths["get"+module+"s"] = true
					allowedPaths["paginate"+module] = true
				case "update":
					allowedPaths["update"+module] = true
					allowedPaths["toggleActive"+module] = true
				default:
					action = utils.LowercaseFirst(action)
					allowedPaths[action+module] = true
				}
			}
		}
	}
	return allowedPaths, nil
}