package cmd

import (
	"fmt"
	"os"

	"github.com/aungmyozaw92/go-graphql/config"
	"github.com/aungmyozaw92/go-graphql/seeder"
	"github.com/spf13/cobra"
)

// Root command for the CLI
var rootCmd = &cobra.Command{
	Use:   "MyApp",
	Short: "A Golang CLI tool for database operations",
	Long:  `This is a CLI tool for managing the application, including database seeding.`,
}



var seedingCommand = &cobra.Command{
	Use:   "db:seed",
	Short: "Seed the database with initial data",
	Long:  `This command will seed the database with default mock data.`,
	Run: func(cmd *cobra.Command, args []string) {

		db := config.GetDB()
		tx := db.Begin()

		seeder.SeedDatabase(tx)
		fmt.Println("Database seeded successfully")

		if err := tx.Commit().Error; err != nil {
			fmt.Println(err)
		}
	},
}

// var freshCommand = &cobra.Command{
// 	Use:   "fresh:seed",
// 	Short: "Fresh database tables and seed",
// 	Run: func(cmd *cobra.Command, args []string) {

// 		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_PORT"), os.Getenv("DB_NAME"))

// 		DB, err := gorm.Open(mysql.Open(DBURL), &gorm.Config{})
// 		if err != nil {
// 			fmt.Println("Failed to connect to the database")
// 			return
// 		}

// 		// DB.Migrator().DropTable(&models.Township{}, &models.City{})
// 		// DB.Migrator().AutoMigrate(&models.City{}, &models.Township{})
// 		fmt.Println("tables freshed successfully")
// 		seeder.SeedDatabase(DB)
// 		fmt.Println("Database seeded successfully")
// 	},
// }

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(seedingCommand)
}