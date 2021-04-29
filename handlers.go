package main

import (
	"context"
	"io"
	"net/http"
	"net/url"
	"os"

	"cloud.google.com/go/storage"
	"github.com/gin-gonic/gin"
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

var (
	storageClient *storage.Client
	secretClient *secretmanager.Client
)


// Uploads file to bucket
func uploadFileHandler(c *gin.Context) {
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

func viewSecretsHandler(c *gin.Context) {
	// Grab Secret name from request
	// secret := c.Params.ByName("name")
	secret := "projects/312654790392/secrets/super-secure-secret/versions/1"
	// Create the client.
	ctx := context.Background()

	secretClient, err := secretmanager.NewClient(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// Build the request.
	req := &secretmanagerpb.AccessSecretVersionRequest{
			Name: secret,
	}

	// Call the API.
	result, err := secretClient.AccessSecretVersion(ctx, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
			"error":   true,
		})
		return
	}

	// WARNING: Do not print the secret in a production environment - this snippet
	// is showing how to access the secret material.
	c.JSON(http.StatusOK, gin.H{
		// "message":  "file uploaded successfully",
		"message":  "secret retrieved successfully",
		"secret value": result.Payload.Data,
	})
}