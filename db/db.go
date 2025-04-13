package db

import (
	"striplex/config"

	"github.com/kokizzu/gotro/L"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection instance
var db *gorm.DB

// Connect initializes the database connection and sets the global DB variable
func Connect() *gorm.DB {
	var err error
	db, err = gorm.Open(postgres.Open(config.GetConfig().GetString("postgres.dsn")), &gorm.Config{})
	L.PanicIf(err, `gorm.Open`, err)
	return db
}

// GetDB returns the global database connection
func GetDB() *gorm.DB {
	return db
}
