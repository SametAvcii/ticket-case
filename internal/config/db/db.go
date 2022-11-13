package db

import (
	"fmt"
	models "samet-avci/gowit/models/ticket"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func configSetup() {

	viper.SetConfigFile("../app/.env")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func Connect() *gorm.DB {
	/*configSetup()
	dbHost := viper.GetString("DB_HOST")
	dbPort := viper.GetString("DB_PORT")
	dbUser := viper.GetString("DB_USER")
	dbPass := viper.GetString("DB_PASSWORD")
	dbName := viper.GetString("DB_NAME")
	//dbtimeZone := viper.GetString("DB_TIMEZONE")

	fmt.Println("port:", dbPort)*/

	dbHost := "localhost"
	dbPort := "5432"
	dbUser := "root"
	dbPass := "secret"
	dbName := "gowit-case-db"

	dbUrl := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", dbHost, dbPort, dbUser, dbName, dbPass)
	db, err := gorm.Open(postgres.Open(dbUrl), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}
	DB = db

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
