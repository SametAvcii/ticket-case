package db

import (
	"database/sql"
	"fmt"
	"os"
	models "samet-avci/gowit/models/ticket"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect() *gorm.DB {
	dbHost := os.Getenv("DATABASE_HOST")
	dbPort := os.Getenv("DATABASE_PORT")
	dbUser := os.Getenv("POSTGRES_USER")
	dbPass := os.Getenv("POSTGRES_PASSWORD")
	dbName := os.Getenv("POSTGRES_DB")

	dbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPass)
	fmt.Println("db", dbUrl)
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})

	sqlDB, err := sql.Open("pgx", dbUrl)
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	DB = gormDB

	err = AutoMigrate()
	if err != nil {
		panic(err.Error())
	}
	return db

}

func AutoMigrate() error {
	migrate := DB.AutoMigrate(
		&models.Ticket{},
		&models.SoldTicket{},
	)
	return migrate
}
