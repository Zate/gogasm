package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Web shiz yo

func initWeb() {
	e := echo.New()
	// e.Static("/static", "static")
	// e.File("/favicon.ico", "favicon.ico")
	// e.File("/common.css", "common.css")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.GET("/", hello)

	e.Logger.Fatal(e.Start(":2086"))

}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
