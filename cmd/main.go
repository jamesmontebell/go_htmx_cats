package main

import (
	"html/template"
	"io"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
	template *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.template.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
	return &Templates{
		template: template.Must(template.ParseGlob("views/*.html")),
	}
}

type Count struct {
	Count int
}

func main() {
	apiKey := os.Getenv("API_KEY")

	e := echo.New()
	e.Use(middleware.Logger())

	count := Count{Count: 0}
	e.Renderer = newTemplate()

	e.GET("/", func(c echo.Context) error {
		count.Count++
		return c.Render(200, "index.html", count)
	})

	e.GET("/cat", func(c echo.Context) error {
		count.Count++
		return c.Render(200, "index.html", count)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
