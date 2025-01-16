package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/GDGVIT/bbiot-backend/internal/auth"
	"github.com/GDGVIT/bbiot-backend/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) login(ctx *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	dbx := database.New(s.db)
	user, err := dbx.GetUserByUsername(ctx, request.Username)
	if errors.Is(err, pgx.ErrNoRows) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("internal server error: %s", err.Error()),
		})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("incorrect password: %s", err.Error()),
		})
	}

	token, err := auth.CreateJWTToken(user.Username, user.Role, auth.JWTSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while creating JWT token",
		})
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
