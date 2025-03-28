package api

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"time"

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
