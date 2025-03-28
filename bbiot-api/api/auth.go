package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/GDGVIT/bbiot-backend/internal/auth"
	"github.com/GDGVIT/bbiot-backend/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

func (s *Server) register(ctx *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if len(request.Password) > 72 {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "password too long!",
		})
		return
	}

	pass_hash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	dbx := database.New(s.db)
	user, err := dbx.NewUser(ctx, database.NewUserParams{
		Username:     request.Username,
		Email:        &request.Email,
		PasswordHash: pass_hash,
	})
	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "user already exists",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("error during registration: %s", err.Error()),
			})
			return
		}
	}

	token, err := auth.CreateJWTToken(user.Username, user.Role, auth.JWTSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while creating JWT token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func (s *Server) login(ctx *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	dbx := database.New(s.db)
	user, err := dbx.GetUserByUsername(ctx, request.Username)
	if errors.Is(err, pgx.ErrNoRows) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "user not found",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("internal server error: %s", err.Error()),
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password)); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("incorrect password: %s", err.Error()),
		})
		return
	}

	token, err := auth.CreateJWTToken(user.Username, user.Role, auth.JWTSecret)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error while creating JWT token",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}
