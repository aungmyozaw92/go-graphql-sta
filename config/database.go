package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

const SearchLimit = 10

var (
	db *gorm.DB
)

func GetDB() *gorm.DB {
	return db
}

func init() {
	// Load env from .env
	godotenv.Load()
	connectDatabase()
}

func connectDatabase() {
	databaseConfig := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true&parseTime=true",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_NAME"),
	)

	var err error
	db, err = gorm.Open(mysql.Open(databaseConfig), initConfig())

	if err != nil {
		panic("Fail To Connect Database")
	}

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}
}

// InitConfig Initialize Config
func initConfig() *gorm.Config {
	return &gorm.Config{
		// Logger: WriteGormLog(),
		Logger:         initLog(),
		NamingStrategy: initNamingStrategy(),
	}
}

// InitLog Connection Log Configuration
func initLog() logger.Interface {

	// Open or create the gorm.log file
	logFile, err := os.OpenFile("gorm.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to create gorm.log: %v", err)
	}

	// Configure GORM logger to write only to gorm.log
	newLogger := logger.New(
		log.New(logFile, "\r\n", log.LstdFlags), // Output only to gorm.log
		logger.Config{
			SlowThreshold: time.Second,   // Log queries that take more than 1 second
			LogLevel:      logger.Info,   // Log level: Info (adjust if needed)
			Colorful:      false,         // Disable colorization (since it's logging to file)
		},
	)

	return newLogger

	// newLogger := logger.New(
	// 	log.New(os.Stdout, "\r\n", log.LstdFlags), // Output to standard output
	// 	logger.Config{
	// 		Colorful:      true,
	// 		LogLevel:      logger.Error, // Adjust log level as needed
	// 		SlowThreshold: time.Second,
	// 	},
	// )
	// return newLogger
}

// InitNamingStrategy Init NamingStrategy
func initNamingStrategy() *schema.NamingStrategy {
	return &schema.NamingStrategy{
		SingularTable: false,
		TablePrefix:   "",
	}
}

// func WriteGormLog() logger.Interface {
// 	logFile := os.Getenv("GORM_LOG")
// 	if logFile == "" {
// 		return initLog()
// 	}

// 	f, _ := os.Create(logFile)
// 	// f, _ := os.Create("gorm.log")
// 	newLogger := logger.New(log.New(io.MultiWriter(f), "\r\n", log.LstdFlags), logger.Config{
// 		Colorful:      true,
// 		LogLevel:      logger.Info,
// 		SlowThreshold: time.Second,
// 	})
// 	return newLogger
// }
