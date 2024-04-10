package main

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/segmentio/encoding/json"
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

type ResponseData []struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
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
		response, err := http.Get("https://api.thecatapi.com/v1/images/search?api_key=" + apiKey)
		if err != nil {
			return c.Render(500, "error.html", nil)
		}
		defer response.Body.Close()

		body, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return c.Render(500, "error.html", nil)
		}

		var responseData ResponseData
		err := json.Unmarshal(body, &responseData)
		if err != nil {
			return c.Render(500, "error.html", nil)
		}
		fmt.Println(responseData)
		return c.Render(200, "index.html", responseData)
	})

	e.Logger.Fatal(e.Start(":8080"))
}
