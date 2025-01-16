package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/GDGVIT/bbiot-backend/internal/auth"
	"github.com/GDGVIT/bbiot-backend/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Server struct {
	Environment string
	port        int

	db *pgxpool.Pool
}

func NewServer() *http.Server {
	auth.InitializeAuth()
	dbpool, err := database.InitialiseDB()
	if err != nil {
		log.Fatal("Error initialising db: ", err)
	}
	database.AutoMigrate()
	log.Printf("Database initalised")

	News := Server{
		port: 8080,
		db:   dbpool,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf("localhost:%d", News.port),
		Handler:      News.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server

}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.POST("/auth/login", s.login)

	r.POST("/products/new", s.NewProduct)
	r.GET("/products", s.GetAllProducts)

	r.POST("/locations/new", s.AddLocation)
	r.GET("/locations", s.GetAllLocations)
	r.POST("/locations/user", s.MapUserToLocation)

	return r
}

func (s *Server) HelloWorldHandler(c *gin.Context) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	c.JSON(http.StatusOK, resp)
}

func (s *Server) healthHandler(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	err := s.db.Ping(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("%s", fmt.Sprintf("db down: %v", err)) // Log the error and terminate the program
		c.JSON(http.StatusInternalServerError, stats)
	}

	// Database is up, add more statistics
	stats["status"] = "up"
	stats["message"] = "It's healthy"
	c.JSON(http.StatusOK, stats)
}
