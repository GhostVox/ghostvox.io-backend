package handlers

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"path/filepath"

	"github.com/GhostVox/ghostvox.io-backend/internal/auth"
	"github.com/GhostVox/ghostvox.io-backend/internal/config"
	"github.com/GhostVox/ghostvox.io-backend/internal/database"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/google/uuid"
)

type AWSS3Handler struct {
	s3Client *s3.Client
	cfg      *config.APIConfig
}

func NewAWSS3Handler(cfg *config.APIConfig, s3Client *s3.Client) *AWSS3Handler {
	return &AWSS3Handler{
		s3Client: s3Client,
		cfg:      cfg,
	}
}

func (h *AWSS3Handler) UpdateUserAvatar(w http.ResponseWriter, r *http.Request) {
	// 1. AUTHENTICATION (from your original code)
	accessToken, err := r.Cookie("accessToken")
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Access token", "Missing access token", err)
		return
	}

	claims, err := auth.ValidateJWT(accessToken.Value, h.cfg.GhostvoxSecretKey)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Access token", "Invalid access token", err)
		return
	}

	userUUID, err := uuid.Parse(claims.Subject)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Malformed token", "Invalid user ID in token", err)
		return
	}

	// 2. FILE RETRIEVAL & VALIDATION
	file, header, err := r.FormFile("image")
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Image", "Image file is required", err)
		return
	}
	defer file.Close()

	if err := validateImageFile(header); err != nil {
		respondWithError(w, http.StatusBadRequest, "Image", "Invalid image", err)
		return
	}

	// 3. GENERATE UNIQUE S3 OBJECT KEY
	// Key format: avatars/{user_id}/{random_uuid}.{extension}
	ext := filepath.Ext(header.Filename)
	objectKey := fmt.Sprintf("avatars/%s%s", userUUID.String(), ext)

	// 4. UPLOAD TO S3
	_, err = h.s3Client.PutObject(r.Context(), &s3.PutObjectInput{
		Bucket:      aws.String(h.cfg.AwsS3Bucket), // Assumes bucket name is in config
		Key:         aws.String(objectKey),
		Body:        file,
		ACL:         types.ObjectCannedACLPublicRead, // Make the file publicly readable
		ContentType: aws.String(header.Header.Get("Content-Type")),
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "", "Failed to upload image", err)
		return
	}

	// 5. CONSTRUCT PUBLIC URL
	imageURL := fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", h.cfg.AwsS3Bucket, h.cfg.AwsRegion, objectKey)

	// 6. UPDATE DATABASE
	_, err = h.cfg.Queries.UpdateUserAvatar(r.Context(), database.UpdateUserAvatarParams{
		PictureUrl: NullStringHelper(imageURL),
		ID:         userUUID,
	})
	if err != nil {
		// In a real app, you might want to delete the uploaded S3 object if the DB update fails.
		respondWithError(w, http.StatusInternalServerError, "Server", "Failed to update user profile", err)
		return
	}

	// 7. RESPOND WITH SUCCESS
	respondWithJSON(w, http.StatusOK, map[string]string{"avatar_url": imageURL})
}

func validateImageFile(header *multipart.FileHeader) error {
	// Check file size (max 1MB)
	if header.Size > 1<<20 {
		return fmt.Errorf("file size exceeds 1MB")
	}

	// Check file type
	// Allowed types: jpeg, png, gif
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	if !allowedTypes[header.Header.Get("Content-Type")] {
		return fmt.Errorf("invalid file type")
	}

	return nil
}
