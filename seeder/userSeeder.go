package seeder

import (
	"fmt"

	"github.com/aungmyozaw92/go-graphql/models"
	"github.com/aungmyozaw92/go-graphql/utils"
	"gorm.io/gorm"
)

func seedUser(tx *gorm.DB) {

	// // Seed Roles

	role, err := models.CreateDefaultRole(tx)
	if err != nil {
		tx.Rollback()
		fmt.Println("Err CreateDefaultRole: " + err.Error())
		return
	}

	// create modules for 
	modules, err := models.CreateDefaultModules(tx)
	if err != nil {
		tx.Rollback()
		fmt.Println("Error CreateDefaultModules: " + err.Error())
		return
	}

	// gives permission to owner
	for _, module := range modules {
		roleModule := models.RoleModule{
			RoleId:         role.ID,
			ModuleId:       module.ID,
			AllowedActions: module.Actions,
		}
		if err := tx.Create(&roleModule).Error; err != nil {
			tx.Rollback()
			fmt.Println("Error seeding role modules: " + err.Error())
			return
		}
	}


	// Seed Users
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		tx.Rollback()
		fmt.Println("Error hashing password: " + err.Error())
		return
	}

	users := []models.User{
		{
			Username: "super_admin",
			Name:     "SuperAdmin",
			Email:    "superadmin@example.com",
			RoleId:  role.ID,
			IsActive: utils.NewTrue(),
			Password: string(hashedPassword),
		},
	}

	err = tx.Create(&users).Error
	if err != nil {
		tx.Rollback()
		fmt.Println("Error seeding users: " + err.Error())
		return
	}

}