package db

import (
	"errors"
	"fmt"

	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func configSetup() {
	viper.SetConfigFile("../config.yaml")
	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %w", err))
	}
}

func Connect() (*gorm.DB, error) {
	dbHost := viper.GetString("database.host")
	dbPort := viper.GetString("database.port")
	dbUserName := viper.GetString("database.username")
	dbPass := viper.GetString("database.password")
	dbName := viper.GetString("database.name")
	dbtimeZone := viper.GetString("database.timezone")

	dsn := "host=" + dbHost + "user=" + dbUserName + " password=" + dbPass + " dbname=" + dbName + " port=" + dbPort + " sslmode=disable TimeZone=" + dbtimeZone
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, errors.New("cannot connect to database")
	}
	return db, nil
}
