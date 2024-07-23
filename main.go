package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/tade3910/blog_aggregator/internal/database"
	feedfollows "github.com/tade3910/blog_aggregator/internal/routes/FeedFollows"
	feeds "github.com/tade3910/blog_aggregator/internal/routes/Feeds"
	users "github.com/tade3910/blog_aggregator/internal/routes/Users"
	v1 "github.com/tade3910/blog_aggregator/internal/routes/v1"
	"github.com/tade3910/blog_aggregator/middleware"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DATABASE_URL")
	if port == "" || dbURL == "" {
		log.Fatal("Could not read port from .env file")
	}
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Could not open postgress server")
	}
	dbQueries := database.New(db)
	router := http.NewServeMux()
	middleware := middleware.NewMiddleWare()
	router.Handle("/v1/", v1.V1(0))
	router.Handle("/users", middleware.EnsureAuthenticated(users.NewUsersHandler(dbQueries)))
	router.Handle("/feeds", middleware.EnsureAuthenticated(feeds.NewFeedsHandler(dbQueries)))
	feedFollowHandler := feedfollows.NewFeedsHandler(dbQueries)
	router.Handle("/feed_follows", middleware.EnsureAuthenticated(feedFollowHandler))
	router.Handle("/feed_follows/", feedFollowHandler)
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	fmt.Printf("Server listening on port %s\n", port)
	server.ListenAndServe()
}
