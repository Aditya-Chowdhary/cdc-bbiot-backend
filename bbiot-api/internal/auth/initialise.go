package auth

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

// JWTSecret is the secret key used to sign the JWT
var JWTSecret string

// InitializeAuth initializes the auth package
func InitializeAuth() {
	JWTSecret = os.Getenv("JWT_SECRET")
}
