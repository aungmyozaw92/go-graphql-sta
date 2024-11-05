package utils

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/storage"
	"google.golang.org/api/option"
)


func getGoogleCredentials() (string, []byte, error) {
	// Read the environment variable
	credentials := os.Getenv("GOOGLE_CLOUD_CREDENTIALS")
	if credentials == "" {
		return "", nil, fmt.Errorf("Google credentials not set in environment")
	}

	// Convert the string to JSON bytes directly
	credBytes := []byte(credentials)

	return "meetup-b8d1b", credBytes, nil
}

// getGoogleClient initializes a Google Cloud Storage client
func getGoogleClient() (*storage.Client, error) {
	// _, credentials := getGoogleCredentials()
	_, credentials, err := getGoogleCredentials()
	if err != nil {
		return nil, err
	}

	client, err := storage.NewClient(context.Background(), option.WithCredentialsJSON(credentials))
	if err != nil {
		return nil, err
	}
	return client, nil
}

func SaveImageToGCS(objectName, imageData string) error {
	// Decode the base64 data
	decodedData, err := base64.StdEncoding.DecodeString(imageData)
	if err != nil {
		return err
	}
	bucketName := os.Getenv("GCS_BUCKET")

	// Get the Google Cloud Storage client
	client, err := getGoogleClient()
	if err != nil {
		return err
	}
	defer client.Close()

	// Upload the decoded image data to the specified object name in your GCS bucket
	contentType := "image/jpeg"
	ctx := context.Background()
	wc := client.Bucket(bucketName).Object(objectName).NewWriter(ctx)
	wc.ContentType = contentType
	wc.Metadata = map[string]string{
		"x-goog-acl": "public-read",

	}

	_, err = wc.Write(decodedData)
	if err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}

	return nil
}


// DeleteImageFromGCS deletes an image from Google Cloud Storage
func DeleteImageFromGCS(objectName string) error {
	// Get the Google Cloud Storage client
	client, err := getGoogleClient()
	if err != nil {
		return err
	}
	defer client.Close()

	bucketName := os.Getenv("GCS_BUCKET")

	// Remove the specified object from your Bucket
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = client.Bucket(bucketName).Object(objectName).Delete(ctx)
	if err != nil {
		// Check if the error is due to the object not existing
		if err == storage.ErrObjectNotExist {
			fmt.Println("Object does not exist:", objectName)
			return nil
		}
		return err
	}

	fmt.Println("Object deleted successfully:", objectName)
	return nil
}

// CheckImageExistInCloud checks if an image exists on the internet
func CheckImageExistInGCS(imageURL string) error {
	resp, err := http.Head(imageURL)
	if err != nil {
		return errors.New("invalid image url")
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	return errors.New("image does not exist")
}

// ObjectExists checks if an object exists in Google Cloud Storage
func ObjectExistsInGCS(objectName string) (bool, error) {
	// Get the Google Cloud Storage client
	client, err := getGoogleClient()
	if err != nil {
		return false, err
	}
	defer client.Close()

	bucketName := os.Getenv("GCS_BUCKET")

	// Attrs is used to check the existence of an object without downloading its content
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	_, err = client.Bucket(bucketName).Object(objectName).Attrs(ctx)
	if err != nil {
		if err == storage.ErrObjectNotExist {
			return false, nil // Object does not exist
		}
		return false, err // Other error
	}

	return true, nil // Object exists
}
