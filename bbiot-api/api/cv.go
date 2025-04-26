// computer-vision.go
package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (s *Server) ComputerVision(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"Camera1": map[string]int{
			"Centering_Sheet": 10,
			"Prop":            20,
			"Span":            30,
		},
		"Camera2": map[string]int{
			"Centering_Sheet": 16,
			"Prop":            22,
			"Span":            34,
		},
		"img": "https://bbiot-imgs.s3.eu-north-1.amazonaws.com/locations/1743542153-vellore.jpg",
	})
}
