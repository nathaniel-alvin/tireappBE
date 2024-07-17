package s3

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"path/filepath"

	tireappconfig "github.com/nathaniel-alvin/tireappBE/config"
	tireapperror "github.com/nathaniel-alvin/tireappBE/error"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var (
	region          = tireappconfig.Envs.BucketRegion
	bucketName      = tireappconfig.Envs.BucketName
	accessKey       = tireappconfig.Envs.S3AccessKey
	secretAccessKey = tireappconfig.Envs.S3SecretAccessKey
)

// upload image to s3 and returns the URL
func UploadImageToS3(imageBytes []byte, filename string) (string, error) {
	objectKey, err := createRandomizedFileName(filename)
	if err != nil {
		return "", err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKey, secretAccessKey, "")),
	)
	if err != nil {
		return "", tireapperror.Errorf(tireapperror.EINTERNAL, "unable to load SDK config: %v", err)
	}

	client := s3.NewFromConfig(cfg)

	// Upload the image to S3
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   bytes.NewReader(imageBytes),
		ACL:    "public-read",
	})
	if err != nil {
		return "", tireapperror.Errorf(tireapperror.EINTERNAL, "unable to upload image to S3: %v", err)
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", bucketName, region, objectKey), nil
}

// generateRandomHexString generates a random hexadecimal string of the given byte length
func generateRandomHexString(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", tireapperror.Errorf(tireapperror.EINTERNAL, "failed to generate random bytes: %v", err)
	}
	return hex.EncodeToString(bytes), nil
}

// createRandomizedFileName generates a new file name by appending a random hex string to the original file name
func createRandomizedFileName(filename string) (string, error) {
	// Generate a random hex string of 16 bytes (32 characters)
	randomHex, err := generateRandomHexString(16)
	if err != nil {
		return "", err
	}

	// Extract the file extension
	ext := filepath.Ext(filename)
	if len(ext) == 0 {
		return "", tireapperror.Errorf(tireapperror.EINVALID, "filename must have an extension")
	}

	// Extract the file name without the extension
	baseName := filename[:len(filename)-len(ext)]
	if len(baseName) == 0 {
		return "", tireapperror.Errorf(tireapperror.EINVALID, "filename must have a base name before the extension")
	}
	// Create the new randomized file name
	randomizedFileName := fmt.Sprintf("%s_%s%s", baseName, randomHex, ext)

	return randomizedFileName, nil
}
