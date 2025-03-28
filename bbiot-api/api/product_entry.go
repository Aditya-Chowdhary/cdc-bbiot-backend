package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/GDGVIT/bbiot-backend/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

func (s *Server) NewProduct(ctx *gin.Context) {
	var request struct {
		ProductName string `json:"product_name" binding:"required"`
		ProductCode string `json:"product_code" binding:"required"`
		Location    string `json:"location" binding:"required"`
		Stock       int    `json:"stock"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	defer tx.Rollback(ctx)
	qtx := database.New(s.db).WithTx(tx)

	product_type, err := qtx.NewProduct(ctx, database.NewProductParams{
		Name: request.ProductName,
		Code: request.ProductCode,
	})
	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "product code already exists",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": fmt.Sprintf("error adding new product: %s", err.Error()),
			})
			return
		}
	}

	locationID, err := qtx.GetLocationByName(ctx, request.Location)
	if errors.Is(err, pgx.ErrNoRows) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "location not found",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("internal server error: %s", err.Error()),
		})
		return
	}

	inventory, err := qtx.NewInventoryProduct(ctx, database.NewInventoryProductParams{
		ProductTypeID: product_type.ID,
		LocationID:    locationID,
		Stock:         int32(request.Stock),
	})
	if err != nil {
		var e *pgconn.PgError
		if errors.As(err, &e) && e.Code == pgerrcode.UniqueViolation {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": "product already exists for given location, update the stock instead",
			})
			return
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error while committing transaction: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"product_type": product_type,
		"inventory":    inventory,
	})
}

func (s *Server) GetAllProducts(ctx *gin.Context) {
	dbx := database.New(s.db)
	products, err := dbx.ListProducts(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "no products found",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error while retrieving products: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"products": products,
	})
}
