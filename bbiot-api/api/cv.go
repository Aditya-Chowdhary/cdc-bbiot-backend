// computer-vision.go
package api

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"
)

func (s *Server) ComputerVision(ctx *gin.Context) {
	obj, err := s.s3.GetObject(&s3.GetObjectInput{
		Bucket: aws.String("ubuntu20"),
		Key:    aws.String("detection_log.csv"),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("error while fetching json code(%s): %s", s3.ErrCodeNoSuchKey, aerr.Error()),
				})
				return
			case s3.ErrCodeInvalidObjectState:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("error while fetching json code(%s): %s", s3.ErrCodeInvalidObjectState, aerr.Error()),
				})
				return
			default:
				ctx.JSON(http.StatusInternalServerError, gin.H{
					"error": fmt.Sprintf("error fetching json file: %s", err.Error()),
				})
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error fetching json file: %s", err.Error()),
		})
		return
	}
	defer obj.Body.Close()

	csvBody := csv.NewReader(obj.Body)
	vals, err := csvBody.ReadAll()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprintf("error reading body: %s", err.Error()),
		})
		return
	}

	// fmt.Println(vals)
	for _, v := range vals {
		v[len(v)-1] = strings.ReplaceAll(v[len(v)-1], "'", "\"")
	}
	// fmt.Println(vals)

	var jsonMap map[string]interface{}
	json.Unmarshal([]byte(vals[len(vals)-1][len(vals[0])-1]), &jsonMap)
	// fmt.Println(jsonMap)

	prop, ok := jsonMap["prop"]
	if !ok {
		prop = 0.0
	}
	span, ok := jsonMap["span"]
	if !ok {
		span = 0.0
	}
	centering_sheet, ok := jsonMap["centering_sheet"]
	if !ok {
		centering_sheet = 0.0
	}
	fmt.Println(prop, span, centering_sheet)

	ctx.JSON(http.StatusOK, gin.H{
		"Camera1": map[string]int{
			"Centering_Sheet": int(centering_sheet.(float64)),
			"Prop":            int(prop.(float64)),
			"Span":            int(span.(float64)),
		},
		"Camera2": map[string]int{
			"Centering_Sheet": int(centering_sheet.(float64)),
			"Prop":            int(prop.(float64)),
			"Span":            int(span.(float64)),
		},
		"img": "https://ubuntu20.s3.eu-north-1.amazonaws.com/frame_1.jpg",
	})
}
