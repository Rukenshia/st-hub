package main

import (
	"log"
	"sthub/lib"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	testController := lib.NewTestController()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Recover())

	testController.RegisterRoutes(e.Group("/iterations"))

	log.Printf("%V", e.Routes())

	// Start server
	e.Logger.Fatal(e.Start(":1323"))
}
