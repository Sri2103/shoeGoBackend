package db

import (
	"fmt"

	"github.com/sri2103/shoeMart/internal/app/utils"
	"github.com/sri2103/shoeMart/internal/db/postgresSchemas"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func ConnectToDB(config *utils.Config) (*gorm.DB, error) {
	connString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=5432 sslmode=require", config.DBHost,config.DBUser,config.DBPass,config.DBName)
	// dsn := "host=postgres user=postgresUser password=postgresPW dbname=postgresDB port=5432 sslmode=disable"
	// dsn := "host=localhost user=postgresUser password=postgresPW dbname=postgresDB port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(connString), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func GetDB() *gorm.DB {
	return db
}

// MigrateModels runs the auto-migration for the models.
func MigrateModels(db *gorm.DB) error {
	return db.AutoMigrate(
		&postgresModels.Product{},
		&postgresModels.User{},
		&postgresModels.Cart{},
		&postgresModels.CartItem{},


	)
}
