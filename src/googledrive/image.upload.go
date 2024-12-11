package googleDrive

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/models"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
	"google.golang.org/api/drive/v3"
)

// UploadImageToDrive uploads an image to Google Drive, checks for duplicates, and returns the image URL.
func (client *GoogleDriveClient) UploadImageToDrive(file multipart.File, fileName string, folderID string) (models.Image, error) {
	logger.Info("Starting image upload process", "fileName", fileName, "folderID", folderID)

	// Validate file extension
	if !isValidImageExtension(fileName) {
		errMsg := "invalid file type"
		logger.Error("Image validation failed", "fileName", fileName, "error", errMsg)
		return models.Image{}, errors.New(errMsg)
	}
	logger.Info("Image extension validated successfully", "fileName", fileName)

	// Convert image to WebP format
	logger.Info("Converting image to WebP format")

	webpBuffer, err := convertToWebP(file, fileName)

	if err != nil {
		logger.Error("Failed to convert image to WebP", "error", err.Error())
		return models.Image{}, fmt.Errorf("failed to convert image to WebP: %v", err)
	}

	fileName = strings.TrimSuffix(fileName, filepath.Ext(fileName)) + ".webp"

	logger.Info("Image converted to WebP format successfully", "fileName", fileName)

	// Compute file hash for duplicate check
	fileHash := hashBuffer(webpBuffer)
	logger.Info("File hash computed successfully", "fileHash", fileHash)

	// Check for duplicates using hash
	isDuplicate, err := isDuplicateImage(client, folderID, fileHash)
	if err != nil {
		logger.Error("Error checking for duplicate images", "error", err.Error())
		return models.Image{}, fmt.Errorf("error checking for duplicate files: %v", err)
	}
	if isDuplicate {
		errMsg := "image already uploaded"
		logger.Warn("Duplicate image detected", "fileHash", fileHash)
		return models.Image{}, errors.New(errMsg)
	}
	logger.Info("No duplicate image found", "fileHash", fileHash)

	// Reset file reader position before uploading
	_, err = file.Seek(0, io.SeekStart)
	if err != nil {
		logger.Error("Failed to reset file pointer before upload", "error", err.Error())
		return models.Image{}, fmt.Errorf("failed to reset file pointer: %v", err)
	}
	logger.Info("File pointer reset successfully for upload")

	// Upload file to Google Drive
	uploadMetadata := &drive.File{
		Name:    fileName,
		Parents: []string{folderID},
		Properties: map[string]string{
			"fileHash": fileHash,
		},
	}

	logger.Info("Initiating file upload to Google Drive", "fileName", fileName, "folderID", folderID)

	uploadRequest := client.service.Files.Create(uploadMetadata).Media(webpBuffer)
	uploadedFile, err := uploadRequest.Do()
	if err != nil {
		logger.Error("Failed to upload image to Google Drive", "fileName", fileName, "error", err.Error())
		return models.Image{}, fmt.Errorf("unable to upload image to Drive: %v", err)
	}

	_, err = client.service.Permissions.Create(uploadedFile.Id, &drive.Permission{
		Type: "anyone",
		Role: "reader",
	}).Do()

	if err != nil {
		errMsg := fmt.Sprintf("Failed to update Permission of file with file ID %s", uploadedFile.Id)
		logger.Error(errMsg, "fileID", uploadedFile.Id)
		return models.Image{}, errors.New(errMsg)
	}

	fileData, err := client.service.Files.Get(uploadedFile.Id).Fields("webContentLink").Do()

	if err != nil {
		log.Fatalf("Unable to retrieve file: %v", err)
	}

	logger.Info("File uploaded successfully", "fileName", fileName, "fileID", uploadedFile.Id)

	// Return the public URL or shared link
	if fileData.WebContentLink == "" {
		errMsg := "uploaded file missing WebContentLink"
		logger.Error(errMsg, "fileID", uploadedFile.Id)
		return models.Image{}, errors.New(errMsg)
	}

	logger.Info("Storing File information in Local database")

	filedata, err := service.CreateGoogleImage(uploadedFile.Name, uploadedFile.Id, strings.Replace(fileData.WebContentLink, "&export=download", "", 1), uploadedFile.MimeType)

	if err != nil {
		logger.Error("Failed to upload image MetaData in Database", "error", err.Error())
		return models.Image{}, fmt.Errorf("unable to upload image metadata to Database: %v", err)
	}

	logger.Info("Image uploaded and accessible via WebViewLink", "URL", uploadedFile.WebViewLink)
	return filedata, nil
}

// isDuplicateImage checks if a file with the same hash already exists in the folder.
func isDuplicateImage(client *GoogleDriveClient, folderID string, fileHash string) (bool, error) {
	logger.Info("Checking for duplicate files in folder", "folderID", folderID, "fileHash", fileHash)
	query := fmt.Sprintf("'%s' in parents and properties has { key='fileHash' and value='%s' }", folderID, fileHash)
	files, err := client.SearchFiles(query)
	if err != nil {
		logger.Error("Error querying Google Drive for duplicates", "query", query, "error", err.Error())
		return false, err
	}
	logger.Info("Duplicate check completed", "fileCount", len(files))
	return len(files) > 0, nil
}

// SearchFiles searches for files on Google Drive matching the given query.
func (client *GoogleDriveClient) SearchFiles(query string) ([]*drive.File, error) {
	logger.Info("Searching files on Google Drive", "query", query)
	files := []*drive.File{}
	ctx := context.TODO()
	err := client.service.Files.List().
		Q(query).
		Fields("files(id, name, properties)").
		Pages(ctx, func(page *drive.FileList) error {
			files = append(files, page.Files...)
			return nil
		})
	if err != nil {
		logger.Error("Error fetching file list from Google Drive", "query", query, "error", err.Error())
		return nil, err
	}
	logger.Info("Files fetched successfully", "fileCount", len(files))
	return files, nil
}

// isValidImageExtension validates the file extension to allow only specific image types.
func isValidImageExtension(fileName string) bool {
	logger.Info("Validating file extension", "fileName", fileName)
	allowedExtensions := []string{".jpg", ".jpeg", ".png", ".webp"}
	lowerFileName := strings.ToLower(fileName)
	for _, ext := range allowedExtensions {
		if strings.HasSuffix(lowerFileName, ext) {
			logger.Info("File extension is valid", "extension", ext)
			return true
		}
	}
	logger.Warn("File extension is invalid", "fileName", fileName)
	return false
}

// convertToWebP converts an image file to WebP format.
// convertToWebP converts any valid image file to WebP format.
func convertToWebP(file multipart.File, fileName string) (*bytes.Buffer, error) {
	logger.Info("Decoding image for WebP conversion")

	// Validate file extension
	if !isValidImageExtension(fileName) {
		return nil, fmt.Errorf("invalid file extension for WebP conversion")
	}

	// Decode the image based on its format
	var img image.Image
	var err error

	// Read the file extension to decide the decoder
	switch {
	case strings.HasSuffix(fileName, ".jpg") || strings.HasSuffix(fileName, ".jpeg"):
		img, err = jpeg.Decode(file)
	case strings.HasSuffix(fileName, ".png"):
		img, err = png.Decode(file)
	case strings.HasSuffix(fileName, ".webp"):
		img, err = webp.Decode(file)
	case strings.HasSuffix(fileName, ".svg"):
		// SVG decoding needs additional processing (e.g., using a package like `github.com/ajstarks/svgo`)
		// For now, we skip SVG handling in this example.
		return nil, fmt.Errorf("SVG format is not supported for WebP conversion")
	default:
		return nil, fmt.Errorf("unsupported image format")
	}

	if err != nil {
		logger.Error("Failed to decode image", "error", err.Error())
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	logger.Info("Image decoded successfully", "fileName", fileName)

	// Convert to WebP format
	webpBuffer := new(bytes.Buffer)
	err = webp.Encode(webpBuffer, img, &webp.Options{Quality: 80})
	if err != nil {
		logger.Error("Failed to encode image to WebP", "error", err.Error())
		return nil, fmt.Errorf("failed to encode image to WebP: %v", err)
	}

	return webpBuffer, nil
}

// hashBuffer computes the MD5 hash of a buffer for duplicate detection.
func hashBuffer(buffer *bytes.Buffer) string {
	logger.Info("Computing hash for buffer")
	hash := md5.Sum(buffer.Bytes())
	return hex.EncodeToString(hash[:])
}
