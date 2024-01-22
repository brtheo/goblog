package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/brtheo/goblog/internal/octogo"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	port   int
	Octogo *octogo.Octogo
	Posts  []*octogo.Post
}

func NewServer() *http.Server {
	octogo := octogo.NewOctogo(
		octogo.User("brtheo"),
		octogo.Repo("blog"),
	)

	posts := octogo.GetAllPosts(10)
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:   port,
		Octogo: octogo,
		Posts:  posts,
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
