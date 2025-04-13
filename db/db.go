package db

import (
	"striplex/config"

	"github.com/kokizzu/gotro/L"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global database connection instance
var DB *gorm.DB

// Connect initializes the database connection and sets the global DB variable
func Connect() {
	var err error
	DB, err = gorm.Open(postgres.Open(config.GetConfig().GetString("database.dsn")), &gorm.Config{})
	L.PanicIf(err, `gorm.Open`, err)
}
