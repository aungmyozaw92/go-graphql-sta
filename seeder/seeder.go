package seeder

import (
	"fmt"

	"gorm.io/gorm"
)

func SeedDatabase(tx *gorm.DB) {
	// Seed data
	seedUser(tx)

}

func ClearDatabase(tx *gorm.DB) error {
	// Truncate or delete data in order of dependencies
	if err := tx.Exec("DELETE FROM role_modules").Error; err != nil {
		return fmt.Errorf("error clearing role_modules: %w", err)
	}
	if err := tx.Exec("DELETE FROM modules").Error; err != nil {
		return fmt.Errorf("error clearing modules: %w", err)
	}
	if err := tx.Exec("DELETE FROM users").Error; err != nil {
		return fmt.Errorf("error clearing users: %w", err)
	}
	if err := tx.Exec("DELETE FROM roles").Error; err != nil {
		return fmt.Errorf("error clearing roles: %w", err)
	}
	
	fmt.Println("Database tables cleared")
	return nil
}
