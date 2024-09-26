package seeder

import (
	"gorm.io/gorm"
)

func SeedDatabase(tx *gorm.DB) {
	// Seed data
	seedUser(tx)

}
