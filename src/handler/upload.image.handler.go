package handler

import (
	"fmt"
	"net/http"
	"strings"

	// Import Cloudinary client
	"github.com/praction-networks/quantum-ISP365/webapp/src/cloudinary"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/response"
)

type ImageUploadHandler struct{}

// NewImageUploadHandler creates a new ImageUploadHandler
func NewImageUploadHandler() *ImageUploadHandler {
	return &ImageUploadHandler{}
}

// UploadImage handles image upload to Cloudinary
func (IU *ImageUploadHandler) UploadImage(w http.ResponseWriter, r *http.Request) {
	logger.Info("Initiating image upload process")

	cloudinaryClient, err := cloudinary.NewCloudinaryCLient()
	if err != nil {
		logger.Error("Failed to initialize Cloudinary client", "error", err.Error())
		response.SendInternalServerError(w, fmt.Sprintf("Error initializing Cloudinary client: %v", err))
		return
	}

	logger.Info("Cloudinary client successfully initialized")

	// Parse file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		logger.Error("Failed to parse file from request", "error", err.Error())
		response.SendBadRequestError(w, "Invalid file upload request")
		return
	}
	defer file.Close()

	// Step 2: Check if the file is empty
	if header.Size == 0 {
		logger.Error("Uploaded file is empty")
		response.SendBadRequestError(w, "Uploaded file is empty")
		return
	}

	// Step 3: Validate the file type (only allow .jpg, .jpeg, .png, .webp)
	if !ValidateFileType(header.Filename) {
		logger.Error("Invalid file type", "fileName", header.Filename)
		response.SendBadRequestError(w, "Invalid file type. Only .jpg, .jpeg, .png, and .webp are allowed.")
		return
	}

	// Log the file details for debugging
	logger.Info("Received file", "fileName", header.Filename, "fileSize", header.Size)

	// Check for duplicate images
	image, err := cloudinaryClient.UploadImage(file, header.Filename, "praction-blog")
	if err != nil {
		switch {
		case strings.Contains(err.Error(), "duplicate image detected"):
			logger.Warn("Duplicate image detected", "fileName", header.Filename)
			response.SendConflictError(w, "Duplicate image detected. The image already exists.")
			return
		case strings.Contains(err.Error(), "unsupported image format"):
			logger.Error("Unsupported image format during processing", "fileName", header.Filename)
			response.SendUnsupportedMediaTypeError(w, "Unsupported image format. Please upload a valid image.")
		case strings.Contains(err.Error(), "file is nil"):
			logger.Error("file is nil, no contenct", "fileName", header.Filename)
			response.SendBadRequestError(w, "nil image.")
		default:
			logger.Error("Failed to upload image to Cloudinary", "error", err.Error())
			response.SendInternalServerError(w, "Error uploading image")
			return
		}
	}
	logger.Info("Image uploaded successfully", "ImageURL", image.ImageURL)

	// Respond with success
	response.SendSuccess(w, map[string]string{
		"image_uuid": image.ImageUUID,
		"image_url":  image.ImageURL}, http.StatusOK)
	logger.Info("Image upload process completed successfully")
}

// ValidateFileType checks if the uploaded file is of a valid image type
func ValidateFileType(fileName string) bool {
	// List of allowed extensions
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}

	// Get the file extension
	fileExt := strings.ToLower(fileName[strings.LastIndex(fileName, "."):])

	// Check if the extension is in the allowed list
	for _, ext := range allowedExtensions {
		if ext == fileExt {
			return true
		}
	}
	return false
}
