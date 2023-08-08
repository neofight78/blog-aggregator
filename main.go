package main

import (
	"com.github/neofight78/blog-aggregator/internal/database"
	"database/sql"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type apiConfig struct {
	DB      *sql.DB
	Queries *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("unable to load configuration: %v", err)
	}

	db, err := sql.Open("mysql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("unable to connect to database: %v", err)
	}

	apiConfig := apiConfig{
		DB:      db,
		Queries: database.New(db),
	}

	go fetchFeeds(apiConfig)

	port := os.Getenv("PORT")

	v1 := chi.NewRouter()
	v1.Get("/readiness", readiness)
	v1.Get("/err", errHandler)
	v1.Post("/users", apiConfig.createUser)
	v1.Get("/users", apiConfig.middlewareAuth(currentUser))
	v1.Post("/feeds", apiConfig.middlewareAuth(apiConfig.createFeed))
	v1.Get("/feeds", apiConfig.listFeeds)
	v1.Post("/feed_follows", apiConfig.middlewareAuth(apiConfig.createFeedFollow))
	v1.Delete("/feed_follows/{feedFollowID}", apiConfig.middlewareAuth(apiConfig.deleteFeedFollow))
	v1.Get("/feed_follows", apiConfig.middlewareAuth(apiConfig.listFeedFollows))
	v1.Get("/posts", apiConfig.middlewareAuth(apiConfig.listPosts))

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
	}))

	router.Mount("/v1", v1)

	err = http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Fatalf("unable to start server: %v", err)
	}
}
