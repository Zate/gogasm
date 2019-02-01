package main

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Web shiz yo

// Needs to get current server ip and query port from ENV on start
// Should display a error page in browser if those are not present, explaining
// how to create those and to restart the server.

// Look at using vue.js to make it all pretty

func initWeb(p string) {
	e := echo.New()
	e.Static("/static", "static")
	e.File("/favicon.ico", "favicon.ico")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HideBanner = true
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.GET("/", hello)
	e.GET("/live", liveHandler)

	e.Logger.Fatal(e.Start(":" + p))

}

// Handlers
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func liveHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Live Stuff Goes Here")
}
