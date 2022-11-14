package main

import (
	"samet-avci/gowit/internal/config/db"

	"samet-avci/gowit/router"

	_ "samet-avci/gowit/docs"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func init() {
	db.Connect()
}

// @title Gowit API Documentation
// @version 1.0
// @description This is a server api docs for
// @termsOfService http://swagger.io/terms/

// @contact.name Samet Avcı
// @contact.url https://www.linkedin.com/in/samet-avci/
// @contact.email sametavc05@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
func main() {

	db := db.Connect()
	e := echo.New()
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	router.Init(e, db)
	e.Logger.Fatal(e.Start(":8080"))

}
