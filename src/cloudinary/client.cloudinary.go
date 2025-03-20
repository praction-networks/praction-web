package cloudinary

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/admin"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/corona10/goimagehash"
	"github.com/praction-networks/quantum-ISP365/webapp/src/config"
	"github.com/praction-networks/quantum-ISP365/webapp/src/logger"
	"github.com/praction-networks/quantum-ISP365/webapp/src/service"
)

type ReturnImageMeta struct {
	ImageURL  string
	ImageUUID string
}

type CloudinaryClient struct {
	client *cloudinary.Cloudinary
}

func NewCloudinaryCLient() (*CloudinaryClient, error) {
	cfg, err := config.CLOUDINARYEnvGet()

	if err != nil {
		logger.Error("Failed to fetch Cloudinary ENV config", "error", err.Error())
		return nil, fmt.Errorf("failed to initialize Cloudinary client: %w", err)
	}

	cloudinaryCloud := cfg.CloudName
	cloudinaryApiKey := cfg.ApiKey
	cloudinaryApiSecret := cfg.ApiSecret

	cld, err := cloudinary.NewFromParams(cloudinaryCloud, cloudinaryApiKey, cloudinaryApiSecret)

	if err != nil {
		logger.Error("Failed to create Cloudinary client", "error", err.Error())
		return nil, fmt.Errorf("unable to create Cloudinary client: %v", err)
	}
	logger.Info("Cloudinary client successfully created")
	return &CloudinaryClient{client: cld}, nil
}

// ConvertToWebP converts a non-WebP image to WebP format.
func (client *CloudinaryClient) ConvertToWebP(file io.Reader, fileName string) (io.Reader, error) {
	var img image.Image
	var err error
	fileExt := strings.ToLower(filepath.Ext(fileName))

	switch fileExt {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	case ".webp":
		img, err = webp.Decode(file)
	default:
		logger.Error("Unsupported image format for conversion", "fileExtension", fileExt)
		return nil, fmt.Errorf("unsupported image format for conversion to WebP, We only support jpg, jpeg, png, webp")
	}

	if err != nil {
		logger.Error("Failed to decode image", "error", err.Error())
		return nil, fmt.Errorf("failed to decode image: %v", err)
	}

	// Convert to WebP
	webpBuffer := new(bytes.Buffer)
	err = webp.Encode(webpBuffer, img, &webp.Options{Quality: 80})
	if err != nil {
		logger.Error("Failed to encode image to WebP", "error", err.Error())
		return nil, fmt.Errorf("failed to encode image to WebP: %v", err)
	}

	logger.Info("Image successfully converted to WebP")
	return webpBuffer, nil
}

// CalculatePerceptualHash calculates the perceptual hash of an image.
func (client *CloudinaryClient) CalculatePerceptualHash(file io.Reader) (string, error) {
	// Create a buffered reader for the file
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		logger.Error("Failed to read image file into buffer", "error", err.Error())
		return "", fmt.Errorf("failed to read image file into buffer: %v", err)
	}

	// Reset the file pointer for further reads
	file = bytes.NewReader(buf.Bytes())

	// Decode the image for perceptual hash
	img, _, err := image.Decode(file)
	if err != nil {
		logger.Error("Failed to decode image for perceptual hash", "error", err.Error())
		return "", fmt.Errorf("failed to decode image: %v", err)
	}

	// Compute the perceptual hash
	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		logger.Error("Failed to compute perceptual hash", "error", err.Error())
		return "", fmt.Errorf("failed to compute perceptual hash: %v", err)
	}

	// Return the perceptual hash as a string
	hashString := fmt.Sprintf("%x", hash.GetHash())
	logger.Info("Perceptual hash successfully calculated", "hash", hashString)
	return hashString, nil
}

// CheckIfImageExists checks if an image already exists in Cloudinary by its PublicID.
func (client *CloudinaryClient) CheckIfImageExists(publicID string) (bool, error) {
	ctx := context.Background()

	// Call the Admin API to fetch asset details by PublicID
	resp, err := client.client.Admin.Asset(ctx, admin.AssetParams{PublicID: "blog/" + publicID})
	if err != nil {
		// Handle error appropriately (log or return custom error)
		logger.Error("Error fetching asset details from Cloudinary", "error", err.Error())
		return false, fmt.Errorf("error fetching asset details from Cloudinary: %v", err)
	}

	// If the response is not nil and contains valid information, the image exists
	if resp != nil && resp.PublicID != "" {
		// Image found, log the secure URL
		logger.Info("Image already exists", "publicID", resp.PublicID, "secureURL", resp.SecureURL)
		return true, nil
	}

	// Image not found
	logger.Warn("Image not found in Cloudinary", "publicID", publicID)
	return false, nil
}

// UploadImage uploads an image to Cloudinary and returns the URL.
func (client *CloudinaryClient) UploadImage(file io.Reader, fileName string, folder string, name string, tag string) (ReturnImageMeta, error) {
	// Step 1: Ensure file is valid and not empty
	if file == nil {
		logger.Error("File is nil", "fileName", fileName)
		return ReturnImageMeta{}, fmt.Errorf("file is nil: %v", fileName)
	}

	// Step 2: Buffer the file once to ensure multiple reads
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(file)
	if err != nil {
		logger.Error("Error reading file into buffer", "fileName", fileName, "error", err.Error())
		return ReturnImageMeta{}, fmt.Errorf("unable to read file into buffer: %v", fileName)
	}

	// Create a separate buffered file for perceptual hash and conversion
	originalFileForHash := bytes.NewReader(buf.Bytes())
	originalFileForUpload := bytes.NewReader(buf.Bytes())

	// Step 3: Calculate perceptual hash
	imageHash, err := client.CalculatePerceptualHash(originalFileForHash)
	if err != nil {
		logger.Error("Error while hashing the file", "error", err.Error())
		return ReturnImageMeta{}, fmt.Errorf("error checking for duplicate image: %v", err)
	}

	// Check if the image already exists based on its perceptual hash
	exists, err := client.CheckIfImageExists(imageHash)
	if err != nil {
		logger.Error("Error while checking if image exists", "error", err.Error())
		return ReturnImageMeta{}, fmt.Errorf("error checking for duplicate image: %v", err)
	}

	if exists {
		logger.Info("Duplicate image detected, upload aborted", "fileName", fileName)
		return ReturnImageMeta{}, fmt.Errorf("duplicate image detected, upload aborted")
	}

	// Log the perceptual hash for debugging
	logger.Info("Generated perceptual hash", "imageHash", imageHash)

	// Step 4: Convert to WebP format
	webPFile, err := client.ConvertToWebP(originalFileForUpload, fileName)
	if err != nil {
		logger.Error("Error converting image to WebP", "error", err.Error())
		return ReturnImageMeta{}, fmt.Errorf("error converting image to WebP: %v", err)
	}

	uploadParams := uploader.UploadParams{
		Folder:         folder,
		Transformation: "f_auto,q_auto", // Apply auto-format and auto-quality transformations
		PublicID:       imageHash,       // Use image hash as public ID
		PublicIDPrefix: "blog",          // Optional: prefix for organization
	}

	// Step 5: Upload the image to Cloudinary
	ctx := context.Background()
	resp, err := client.client.Upload.Upload(ctx, webPFile, uploadParams)
	if err != nil {
		logger.Error("Unable to upload image to Cloudinary", "fileName", fileName, "error", err.Error())
		return ReturnImageMeta{}, fmt.Errorf("unable to upload image '%s' to Cloudinary: %v", fileName, err)
	}

	// Check for errors in the Cloudinary response
	if resp.Error.Message != "" {
		logger.Error("Error in Cloudinary response", "errorMessage", resp.Error.Message)
		return ReturnImageMeta{}, fmt.Errorf("cloudinary error: %s", resp.Error.Message)
	}

	// Log the successful upload
	logger.Info("Image uploaded successfully", "imageURL", resp.SecureURL)

	// Step 6: Save the image metadata (including perceptual hash) to the database
	ImageMeta, err := service.SaveblogImage(resp.DisplayName, resp.PublicID, resp.SecureURL, name, tag)
	if err != nil {
		logger.Error("Failed to upload image metadata to the database", "error", err.Error())
		return ReturnImageMeta{}, fmt.Errorf("unable to upload image metadata to Database: %v", err)
	}

	returnImageMeta := ReturnImageMeta{
		ImageURL:  resp.SecureURL,
		ImageUUID: ImageMeta.UUID,
	}

	return returnImageMeta, nil
}

func (c *CloudinaryClient) DeleteImage(ctx context.Context, publicID string) error {
	// Perform deletion on Cloudinary
	_, err := c.client.Upload.Destroy(ctx, uploader.DestroyParams{
		PublicID: publicID,
	})
	if err != nil {
		return fmt.Errorf("failed to delete image from Cloudinary: %w", err)
	}
	return nil
}
