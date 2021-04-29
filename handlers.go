package main

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	sm "cloud.google.com/go/secretmanager/apiv1"
	pb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var (
	storageClient *storage.Client
)

// Uploads file to bucket
func UploadFileHandler(c *gin.Context) {
	bucket := os.Getenv("BUCKET_NAME")

	var err error

	ctx := context.Background()

	storageClient, err = storage.NewClient(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	f, uploadedFile, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	defer f.Close()

	sw := storageClient.Bucket(bucket).Object(uploadedFile.Filename).NewWriter(ctx)

	if _, err := io.Copy(sw, f); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	if err := sw.Close(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	u, err := url.Parse("/" + bucket + "/" + sw.Attrs().Name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"Error":   true,
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "file uploaded successfully",
		"pathname": u.EscapedPath(),
	})
}

// Get Secret
func ViewSecretsHandler(c *gin.Context) {
	projectId := os.Getenv("PROJECT_ID")
	secretName := os.Getenv("SECRET_NAME")

	ctx := context.Background()

	secretClient, err := sm.NewClient(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Build the request.
	req := pb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%v/secrets/%v/versions/latest", projectId, secretName),
	}

	// Call the API.
	result, err := secretClient.AccessSecretVersion(ctx, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Get Secret value and return
	c.JSON(http.StatusOK, gin.H{
		"message":  result.Payload.Data,
	})
}