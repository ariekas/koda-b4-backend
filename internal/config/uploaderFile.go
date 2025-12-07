package config

import (
	"context"
	"fmt"
	"mime/multipart"

	"github.com/cloudinary/cloudinary-go/api/uploader"
)

func UploaderFile(file *multipart.FileHeader, folder string, fileName string) (string, error) {
	cld, err := ConfigCloudinary()
	if err != nil {
		return "", fmt.Errorf("cloudinary connection failed: %w", err)
	}

	src, err := file.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open file: %w", err)
	}
	defer src.Close()

	params := uploader.UploadParams{
		PublicID: fileName, 
		Folder:   folder,  
	}

	result, err := cld.Upload.Upload(context.Background(), src, params)
	if err != nil {
		return "", fmt.Errorf("cloudinary upload failed: %w", err)
	}

	return result.SecureURL, nil
}
