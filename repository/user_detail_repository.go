package repository

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/url"
	"os"
	"time"

	"be-b-impact.com/csr/model"
	firebase "firebase.google.com/go/v4"

	"gorm.io/gorm"
)

type UserDetailRepository interface {
	Get(id string) (*model.UserDetail, error)
	FirebaseSave(payload multipart.File) (string, error)
	SaveTrx(payload *model.UserDetail, tx *gorm.DB) error
	DeleteTrx(id string, tx *gorm.DB) error
}

type userDetailRepository struct {
	db *gorm.DB
	fb *firebase.App
}

func (ud *userDetailRepository) DeleteTrx(id string, tx *gorm.DB) error {
	return tx.Delete(&model.UserDetail{}, "id=?", id).Error
}

func (ud *userDetailRepository) Get(id string) (*model.UserDetail, error) {
	var userDetail model.UserDetail
	result := ud.db.First(&userDetail, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &userDetail, nil
}

func (ud *userDetailRepository) FirebaseSave(payload multipart.File) (string, error) {
	ctx := context.Background()

	// Create a storage reference for the file
	storageClient, err := ud.fb.Storage(ctx)
	if err != nil {
		return "", err
	}

	bucket, err := storageClient.Bucket(os.Getenv("BUCKET_NAME"))
	if err != nil {
		return "", err
	}

	// Generate a unique filename for the file
	filename := generateUniqueUserImagename()

	// Specify the path where the file will be stored in the bucket
	filePath := "userDetails/" + filename

	// Create a storage object reference
	obj := bucket.Object(filePath)

	// Upload the file to Firebase Storage
	wc := obj.NewWriter(ctx)
	if _, err := io.Copy(wc, payload); err != nil {
		wc.Close()
		return "", err
	}
	if err := wc.Close(); err != nil {
		return "", err
	}

	// Get the download URL for the uploaded file
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		return "", err
	}

	firebaseUrl := fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o/%s?alt=media&", attrs.Bucket, url.QueryEscape(attrs.Name))

	return firebaseUrl, nil
}

func generateUniqueUserImagename() string {
	// Implement your own logic to generate a unique filename
	// You can use a timestamp, random string, or any other method
	// For example:
	return "userImage_" + time.Now().Format("20060102150405") + ".jpg"
}

func (r *userDetailRepository) SaveTrx(payload *model.UserDetail, tx *gorm.DB) error {
	// If the provided transaction is not nil, use it for saving the UserDetail
	if tx != nil {
		return tx.Create(payload).Error
	}

	// Otherwise, use the default DB connection for saving the UserDetail
	return r.db.Create(payload).Error
}

func NewUserDetailRepository(db *gorm.DB, fb *firebase.App) UserDetailRepository {
	return &userDetailRepository{
		db: db,
		fb: fb,
	}
}
