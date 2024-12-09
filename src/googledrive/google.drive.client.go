package googleDrive

import (
	"context"
	"fmt"

	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
)

type GoogleDriveClient struct {
	service *drive.Service
}

// NewGoogleDriveClient initializes a Google Drive client using a service account.
func NewGoogleDriveClient(ctx context.Context, credentialsFile string) (*GoogleDriveClient, error) {
	client, err := drive.NewService(ctx, option.WithCredentialsFile(credentialsFile))
	if err != nil {
		return nil, fmt.Errorf("unable to create Drive client with file '%s': %v", credentialsFile, err)
	}
	return &GoogleDriveClient{service: client}, nil
}

func (client *GoogleDriveClient) GetService() *drive.Service {
	return client.service
}

// GetOrCreateFolder checks if a folder exists on Google Drive. If not, it creates it.
func (client *GoogleDriveClient) GetOrCreateFolder(folderName string) (string, error) {
	query := fmt.Sprintf("mimeType = 'application/vnd.google-apps.folder' and name = '%s'", folderName)
	pageToken := ""
	for {
		files, err := client.service.Files.List().Q(query).PageToken(pageToken).Do()
		if err != nil {
			return "", fmt.Errorf("unable to search for folder '%s': %v", folderName, err)
		}

		if len(files.Files) > 0 {
			return files.Files[0].Id, nil // Folder already exists
		}

		if files.NextPageToken == "" {
			break
		}
		pageToken = files.NextPageToken
	}

	// Create a new folder if it doesn't exist
	folder := &drive.File{
		Name:     folderName,
		MimeType: "application/vnd.google-apps.folder",
	}
	createdFolder, err := client.service.Files.Create(folder).Do()
	if err != nil {
		return "", fmt.Errorf("unable to create folder '%s': %v", folderName, err)
	}

	return createdFolder.Id, nil
}
