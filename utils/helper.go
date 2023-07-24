package utils

import (
	"errors"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"

	"be-b-impact.com/csr/delivery/api"
	"be-b-impact.com/csr/utils/authenticator"
	"github.com/gin-gonic/gin"
)

func AccessInsideToken(b api.BaseApi, c *gin.Context) authenticator.AccessDetail {
	user, exists := c.Get("user")
	if !exists {
		b.NewFailedResponse(c, http.StatusUnauthorized, "token invalid")
		return authenticator.AccessDetail{}
	}
	userTyped, ok := user.(authenticator.AccessDetail)
	if !ok {
		b.NewFailedResponse(c, http.StatusUnauthorized, "user invalid")
		return authenticator.AccessDetail{}
	}
	return userTyped
}

const maxImageSize = 2 << 20 // 2MB in bytes
const maxFileSize = 5 << 20  // 5MB in bytes

func ValidateImage(image multipart.File, header *multipart.FileHeader) error {
	// Check the image size
	if header.Size > maxImageSize {
		return errors.New("image size exceeds the maximum allowed size (2MB)")
	}

	// Check the image format
	ext := strings.ToLower(filepath.Ext(header.Filename))
	validFormats := []string{".png", ".jpg", ".jpeg"}

	validFormat := false
	for _, format := range validFormats {
		if ext == format {
			validFormat = true
			break
		}
	}

	if !validFormat {
		return errors.New("invalid image format. Supported formats: png, jpg, jpeg")
	}

	return nil
}

func ValidateFile(file multipart.File, header *multipart.FileHeader) error {
	// Check the file size
	if header.Size > maxFileSize {
		return errors.New("file size exceeds the maximum allowed size (5MB)")
	}

	// Check the file format
	ext := strings.ToLower(filepath.Ext(header.Filename))
	validFormats := []string{".png", ".jpg", ".jpeg", ".pdf"}

	validFormat := false
	for _, format := range validFormats {
		if ext == format {
			validFormat = true
			break
		}
	}

	if !validFormat {
		return errors.New("invalid file format. Supported formats: pdf, png, jpg, jpeg")
	}

	return nil
}
