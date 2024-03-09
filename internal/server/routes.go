package server

import (
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

	e.GET("/", echo.WrapHandler(templ.Handler(web.HomePage())))
	e.GET("/blog", echo.WrapHandler(templ.Handler(web.BlogHomePage(s.Posts))))
	e.GET("/blog/:slug", s.BlogPostHandler)
	e.POST("/blog/comment", s.CommitCommentHandler)

	return e
}

func (s *Server) BlogPostHandler(c echo.Context) error {
	slug := c.Param("slug")
	post := s.Octogo.GetPostBySlug(slug)
	comments := s.Octogo.GetCommentsBySlug(slug)
	return HTML(c, web.BlogPostPage(post, comments))
}
func (s *Server) CommitCommentHandler(c echo.Context) error {
	props := map[string]string{}
	props["author"] = c.FormValue("author")
	props["body"] = c.FormValue("body")
	sha := c.FormValue("sha")
	return HTML(c, web.Comment(*s.Octogo.CommitComment(props, sha)))
}

func HTML(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)
	return cmp.Render(c.Request().Context(), c.Response().Writer)
}
