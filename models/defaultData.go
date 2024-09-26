package models

import (
	"gorm.io/gorm"
)


func GetDefaultModules() map[string]string {
	defaultModules := map[string]string{
		"User":    "create;update;delete;read;resetPassword",
	}
	return defaultModules
}

func CreateDefaultRole(tx *gorm.DB) (*Role, error){
	role := Role{
		Name: "Admin",
	}

	err := tx.Create(&role).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}
	return &role, err
}

func CreateDefaultModules(tx *gorm.DB) ([]Module, error) {

	defaultModules := GetDefaultModules()

	var modules []Module
	for k, v := range defaultModules {
		modules = append(modules, Module{
			Name:       k,
			Actions:    v,
		})
	}

	if err := tx.Create(&modules).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return modules, nil
}