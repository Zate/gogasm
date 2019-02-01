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
	e.Debug = true
	// e.Static("/static", "static")
	e.File("/favicon.ico", "favicon.ico")
	// e.File("/common.css", "common.css")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HideBanner = true
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	// a1url, err := url.Parse("http://atlas.zate.systems")
	// if err != nil {
	// 	e.Logger.Fatal(err)
	// }
	// a1target := []*middleware.ProxyTarget{
	// 	{
	// 		Name: "A1",
	// 		URL:  a1url,
	// 	},
	// }

	// rrb := middleware.NewRoundRobinBalancer(a1target)

	// g := e.Group("/a1")
	// g.Use(middleware.ProxyWithConfig(rrb))

	a1url, err := url.Parse("http://atlas.zate.systems:80")
	if err != nil {
		e.Logger.Fatal(err)
	}

	targets := []*middleware.ProxyTarget{
		{
			URL: a1url,
		},
	}
	e.GET("/", hello)
	e.GET("/live", liveHandler)

	g := e.Group("/a1")

	g.Use(middleware.Proxy(middleware.NewRoundRobinBalancer(targets)))

	e.Logger.Fatal(e.Start(":" + p))

}

// Handlers
func hello(c echo.Context) error {
	return c.String(http.StatusOK, "Hello, World!")
}

func liveHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Live Stuff Goes Here")
}
