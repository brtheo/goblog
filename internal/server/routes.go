package server

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
	"github.com/brtheo/goblog/cmd/web"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// fileServer := http.FileServer(http.FS(web.Files))
	e.Static("/static", "cmd/web/assets")
	// e.GET("/js/*", echo.WrapHandler(fileServer))
	// e.GET("/css/*", echo.WrapHandler(fileServer))
	// e.GET("/css/**/*", echo.WrapHandler(fileServer))

	e.GET("/blog", echo.WrapHandler(templ.Handler(web.BlogHomePage(s.Posts))))
	e.GET("/blog/:slug", s.BlogPostHandler)
	// e.POST("/hello", echo.WrapHandler(http.HandlerFunc(web.HelloWebHandler)))

	e.GET("/", s.HelloWorldHandler)

	return e
}

func (s *Server) HelloWorldHandler(c echo.Context) error {
	fmt.Println(s.Posts[0].Title)
	resp := map[string]string{
		"message": s.Octogo.GetPostBySlug(s.Posts[3].Slug).Title,
	}

	return c.JSON(http.StatusOK, resp)
}
func (s *Server) BlogPostHandler(c echo.Context) error {
	return HTML(c, web.BlogPostPage(s.Octogo.GetPostBySlug(c.Param("slug"))))
}

func HTML(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
