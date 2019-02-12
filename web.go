package main

import (
	"html/template"
	"io"
	"net/http"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

// Web shiz yo

// Needs to get current server ip and query port from ENV on start
// Should display a error page in browser if those are not present, explaining
// how to create those and to restart the server.

// Look at using vue.js to make it all pretty

// TemplateRegistry struct
type TemplateRegistry struct {
	templates *template.Template
}

// Render Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func initWeb(p string) {
	e := echo.New()
	e.Static("/css", "frontend/css")
	e.Static("/js", "frontend/js")
	e.File("/favicon.ico", "favicon.ico")
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.HideBanner = true
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.PUT, echo.POST, echo.DELETE},
	}))

	e.Renderer = &TemplateRegistry{
		templates: template.Must(template.ParseGlob("frontend/*.html")),
	}

	e.GET("/", hello)
	e.GET("/live", liveHandler)

	e.Logger.Fatal(e.Start(":" + p))

}

// Handlers
func hello(c echo.Context) error {
	// return c.String(http.StatusOK, "Hello, World!")

	return c.Render(http.StatusOK, "index.html", map[string]interface{}{
		"name":   "Go Gaming Automated Server Manager",
		"msg":    "Hello, Yeeter",
		"author": "Bobblehead",
		"desc":   "Manager for Atlas Server Grid on Docker",
	})
}

func liveHandler(c echo.Context) error {
	return c.String(http.StatusOK, "Live Stuff Goes Here")
}
