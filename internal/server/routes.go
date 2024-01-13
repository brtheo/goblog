package server

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/brtheo/goblog/cmd/web"
	"github.com/brtheo/goblog/internal/octogo"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	fileServer := http.FileServer(http.FS(web.Files))
	e.GET("/js/*", echo.WrapHandler(fileServer))

	e.GET("/web", echo.WrapHandler(templ.Handler(web.HelloForm())))
	e.POST("/hello", echo.WrapHandler(http.HandlerFunc(web.HelloWebHandler)))

	e.GET("/", s.HelloWorldHandler)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	octogo := octogo.NewOctogo(
		octogo.
			NewOctoConf().
			Repo("blog").
			User("brtheo"),
	)
	fmt.Println(octogo.GetAllPosts(7)[0].Title)
	resp := map[string]string{
		"message": "Hello World",
	}

	return c.JSON(http.StatusOK, resp)
}
