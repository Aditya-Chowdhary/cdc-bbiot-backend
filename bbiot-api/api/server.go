package api

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GDGVIT/bbiot-backend/internal/auth"
	"github.com/GDGVIT/bbiot-backend/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	_ "github.com/joho/godotenv/autoload"
)

type Server struct {
	Environment string
	port        int

	db       *pgxpool.Pool
	s3       *s3.S3
	uploader *s3manager.Uploader
}

func NewServer() *http.Server {
	auth.InitializeAuth()
	dbpool, err := database.InitialiseDB()
	if err != nil {
		log.Fatal("Error initialising db: ", err)
	}
	database.AutoMigrate()
	log.Printf("Database initalised")

	s3, uploader := InitialiseS3()

	News := Server{
		port:     8080,
		db:       dbpool,
		s3:       s3,
		uploader: uploader,
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", News.port),
		Handler:      News.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server

}

func InitialiseS3() (*s3.S3, *s3manager.Uploader) {
	sess, err := session.NewSession(&aws.Config{
		Region:      aws.String(os.Getenv("AWS_REGION")),
		Credentials: credentials.NewStaticCredentials(os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), ""),
	})
	if err != nil {
		return nil, nil
	}

	svc := s3.New(sess)
	uploader := s3manager.NewUploader(sess)
	return svc, uploader
}

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.Use(CORSMiddleware())

	r.GET("/", s.HelloWorldHandler)

	r.GET("/health", s.healthHandler)

	r.POST("/auth/login", s.login)
	r.POST("/auth/register", s.register)

	r.POST("/products/new", s.NewProduct)
	r.GET("/products", s.GetAllProducts)
	r.POST("/products/update", s.UpdateInventory)

	r.GET("/inventory/:location", s.TransferList)
	r.GET("/details/:location", s.TransferDetails)

	r.POST("/locations/new", s.AddLocation)
	r.GET("/locations", s.GetAllLocations)
	r.POST("/locations/user", s.MapUserToLocation)

	r.GET("/cv", s.ComputerVision)

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
