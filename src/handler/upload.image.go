package handler

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/google/uuid"
	googleDrive "github.com/praction-networks/quantum-ISP365/webapp/src/googledrive"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

type ImageUploadHandler struct{}

// UploadImage handles image upload to Google Drive
func (IU *ImageUploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	logger.Info("Initiating image upload process")

	ctx := context.Background()
	// Create Google Drive client
	client, err := googleDrive.NewGoogleDriveClient(ctx, "src/config/praction-web-12376f05facd.json")
	if err != nil {
		logger.Error("Failed to create Google Drive client", "error", err.Error())
		response.SendInternalServerError(w, "Failed to create Google Drive client")
		return
	}
	logger.Info("Google Drive client created successfully")

	// Access or create the target folder
	folderID, err := client.GetOrCreateFolder("web-blog")
	if err != nil {
		logger.Error("Failed to access or create the Google Drive folder", "error", err.Error())
		response.SendInternalServerError(w, "Error accessing Google Drive folder")
		return
	}
	logger.Info("Google Drive folder verified or created successfully", "FolderID", folderID)

	// Parse file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		logger.Error("Failed to parse file from request", "error", err.Error())
		response.SendBadRequestError(w, "Invalid file upload request")
		return
	}
	defer file.Close()
	logger.Info("File received successfully", "OriginalFilename", header.Filename)

	// Generate a unique filename
	originalFilename := header.Filename
	fileExt := filepath.Ext(originalFilename)
	if fileExt == "" {
		logger.Error("File has no extension", "Filename", originalFilename)
		response.SendBadRequestError(w, "File has no valid extension")
		return
	}

	uniqueFilename := fmt.Sprintf("%s%s", uuid.New().String(), fileExt)
	logger.Info("Generated unique filename", "UniqueFilename", uniqueFilename)

	// Upload the image to Google Drive
	imageMeta, err := client.UploadImageToDrive(file, uniqueFilename, folderID)
	if err != nil {
		logger.Error("Failed to upload image to Google Drive", "error", err.Error())
		imageError := fmt.Sprintf("Error uploading image: %v", err)
		response.SendInternalServerError(w, imageError)
		return
	}
	logger.Info("Image uploaded successfully", "ImageURL", imageMeta.ImageURL)

	// Respond with success
	response.SendSuccess(w, imageMeta, http.StatusOK)
	logger.Info("Image upload process completed successfully")
}
