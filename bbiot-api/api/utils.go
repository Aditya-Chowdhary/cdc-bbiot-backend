package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func (s *Server) uploadImg(imgData, folder, name string) (string, error) {
	bucketName := "bbiot-imgs"

	data, err := base64.StdEncoding.DecodeString(imgData)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64: %w", err)
	}

	key := fmt.Sprintf("%s/%d-%s.jpg", folder, time.Now().Unix(), name)
	_, err = s.uploader.Upload(&s3manager.UploadInput{
		Bucket: &bucketName,
		Key:    &key,
		Body:   bytes.NewReader(data),
	})

	// _, err := s3Client.PutObject(&s3.PutObjectInput{
	// 	Bucket: aws.String(bucketName),
	// 	Key:    aws.String(key),
	// 	Body:   reader,
	// })
	if err != nil {
		return "", fmt.Errorf("failed to upload to S3: %w", err)
	}

	s3URL := fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, key)
	return s3URL, nil
}

type jsonObject struct {
	Timestamp         time.Time `json:"timestamp"`
	EventType         string    `json:"eventType"`
	TagInventoryEvent *struct {
		EpcHex      string `json:"epcHex"`
		AntennaPort int    `json:"antennaPort"`
		Frequency   int    `json:"frequency"`
	} `json:"tagInventoryEvent"`
}

func getJsonFile(s3Client *s3.S3, bucket, key string) ([]jsonObject, error) {
	obj, err := s3Client.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case s3.ErrCodeNoSuchKey:
				return nil, fmt.Errorf("error while fetching json code(%s): %w", s3.ErrCodeNoSuchKey, aerr)
			case s3.ErrCodeInvalidObjectState:
				return nil, fmt.Errorf("error while fetching json code(%s): %w", s3.ErrCodeInvalidObjectState, aerr)
			default:
				return nil, fmt.Errorf("error fetching json file: %w", err)
			}
		}
		return nil, fmt.Errorf("error fetching json file: %w", err)
	}
	defer obj.Body.Close()

	decoder := json.NewDecoder(obj.Body)
	_, err = decoder.Token()
	if err != nil {
		return nil, fmt.Errorf("JSON Decoding Error: %w", err)
	}

	var objects []jsonObject

	for decoder.More() {
		var obj jsonObject
		err := decoder.Decode(&obj) // Try to decode into the struct
		if obj.TagInventoryEvent != nil {
			objects = append(objects, obj) // Only add valid objects
		}
		if err != nil {
			return nil, fmt.Errorf("error unrelated to invalid struct: %w", err)
		}
	}

	return objects, nil
}
