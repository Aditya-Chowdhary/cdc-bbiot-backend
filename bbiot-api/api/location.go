package api

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/GDGVIT/bbiot-backend/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func (s *Server) AddLocation(ctx *gin.Context) {
	var request struct {
		Location  string `json:"location" binding:"required"`
		Address   string `json:"address" binding:"required"`
		Site_Desc string `json:"site_desc" binding:"required"`
		Notes     string `json:"notes"`
		Img       string `json:"img"`
	}

	if err := ctx.ShouldBindJSON(&request); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": fmt.Sprintf("bad input: %s", err.Error()),
		})
		return
	}

	fmt.Println("Uploading to s3")
	url, err := s.uploadImg(request.Img, "locations", request.Location)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error uploading to s3: %s", err.Error()),
		})
		return
	}
	fmt.Printf("Done uploading to s3: %s\n", url)

	dbx := database.New(s.db)
	location, err := dbx.AddLocation(ctx, database.AddLocationParams{
		Location:        request.Location,
		Address:         request.Address,
		SiteDesc:        request.Site_Desc,
		AdditionalNotes: &request.Notes,
		ImgUrl:          &url,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error while adding location: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"location": location,
	})
}

func (s *Server) GetAllLocations(ctx *gin.Context) {
	dbx := database.New(s.db)
	locations, err := dbx.ListLocations(ctx)
	if errors.Is(err, pgx.ErrNoRows) {
		ctx.JSON(http.StatusNotFound, gin.H{
			"error": "no locations found",
		})
		return
	} else if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error while retrieving locations: %s", err.Error()),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"locations": locations,
	})
}

func (s *Server) MapUserToLocation(ctx *gin.Context) {
	var request struct {
		Username string `json:"username" binding:"required"`
		Location string `json:"location" binding:"required"`
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

	user, err := qtx.GetUserByUsername(ctx, request.Username)
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

	location, err := qtx.GetLocationByName(ctx, request.Location)
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

	qtx.MapUserIDToLocationID(ctx, database.MapUserIDToLocationIDParams{
		LocationID: &location,
		ID:         user.ID,
	})

}
