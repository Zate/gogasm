package main

import (
	"net/http"
	"net/url"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Web shiz yo

// Needs to get current server ip and query port from ENV on start
// Should display a error page in browser if those are not present, explaining
// how to create those and to restart the server.

// Ability to crawl Live servers and create a view on that
// Live container stats via glances.  Need to reverse proxy to each container
// for that.

func initWeb(p string) {
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

	a1url, err := url.Parse("http://atlas-s:61208")
	if err != nil {
		e.Logger.Fatal(err)
	}
	a1target := []*middleware.ProxyTarget{
		{
			URL: a1url,
		},
	}

	rrb := middleware.NewRandomBalancer(a1target)

	e.Group("/a1*")
	e.Use(middleware.Proxy(rrb))

	e.GET("/", hello)

	e.Logger.Fatal(e.Start(":" + p))

}

// Handler
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}
