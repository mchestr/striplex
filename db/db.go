package db

import (
	"striplex/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection instance
var DB *gorm.DB

// Connect initializes the database connection and sets the global DB variable
func Connect() error {
	var err error
	DB, err = gorm.Open(postgres.Open(config.Config.GetString("database.dsn")), &gorm.Config{})
	return err
}
