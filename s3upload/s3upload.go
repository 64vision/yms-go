package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
)

var presignClient *s3.PresignClient

const bucketName = "zerasuite"

func InitS3() error {
	cfg, err := config.LoadDefaultConfig(
		context.Background(),
		config.WithRegion("ap-southeast-1"),
	)
	if err != nil {
		return err
	}

	client := s3.NewFromConfig(cfg)
	presignClient = s3.NewPresignClient(client)
	fmt.Println("NO error")
	return nil
}

type PresignRequest struct {
	FolderName  string `json:"folder_name"`
	FileName    string `json:"file_name"`
	ContentType string `json:"content_type"`
}

type PresignResponse struct {
	UploadURL string `json:"upload_url"`
	Key       string `json:"key"`
}

func GetPresignedUploadURL(w http.ResponseWriter, r *http.Request) {
	var input PresignRequest

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	ext := filepath.Ext(input.FileName)

	key := input.FolderName + "/" + uuid.New().String() + ext
	fmt.Println(input.ContentType)
	result, err := presignClient.PresignPutObject(
		r.Context(),
		&s3.PutObjectInput{
			Bucket: aws.String(bucketName),
			Key:    aws.String(key),
			//ContentType: aws.String(input.ContentType),
		},
		func(options *s3.PresignOptions) {
			options.Expires = 15 * time.Minute
		},
	)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	response := PresignResponse{
		UploadURL: result.URL,
		Key:       key,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
