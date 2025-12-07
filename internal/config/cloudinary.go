package config

import (
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/joho/godotenv"
)

func ConfigCloudinary() (*cloudinary.Cloudinary, error) {
	godotenv.Load()

	cldName := os.Getenv("CLOUDINARY_NAME")
	cldKey := os.Getenv("CLOUDINARY_API_KEY")
	cldSecret := os.Getenv("CLOUDINARY_API_SECRET")

	return cloudinary.NewFromParams(cldName, cldKey, cldSecret)
}
