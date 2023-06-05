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
	"be-b-impact.com/csr/model/dto"
	firebase "firebase.google.com/go/v4"

	"be-b-impact.com/csr/utils/common"
	"gorm.io/gorm"
)

type ImageRepository interface {
	BaseRepository[model.Image]
	BaseRepositoryCount[model.Image]
	BaseRepositoryPaging[model.Image]
	FirebaseSave(payload multipart.File) (string, error)
}
type imageRepository struct {
	db *gorm.DB
	fb *firebase.App
}

func (im *imageRepository) Delete(id string) error {
	return im.db.Delete(&model.Image{}, "id=?", id).Error
}

func (im *imageRepository) Get(id string) (*model.Image, error) {
	var image model.Image
	result := im.db.First(&image, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &image, nil
}

func (im *imageRepository) List() ([]model.Image, error) {
	var image []model.Image
	result := im.db.Find(&image).Error
	if result != nil {
		return nil, result
	}
	return image, nil
}

func (im *imageRepository) FirebaseSave(payload multipart.File) (string, error) {
	ctx := context.Background()

	// Create a storage reference for the file
	storageClient, err := im.fb.Storage(ctx)
	if err != nil {
		return "", err
	}

	bucket, err := storageClient.Bucket(os.Getenv("BUCKET_NAME"))
	if err != nil {
		return "", err
	}

	// Generate a unique filename for the file
	filename := generateUniqueImagename()

	// Specify the path where the file will be stored in the bucket
	filePath := "images/" + filename

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

func generateUniqueImagename() string {
	// Implement your own logic to generate a unique filename
	// You can use a timestamp, random string, or any other method
	// For example:
	return "image_" + time.Now().Format("20060102150405") + ".jpg"
}

func (im *imageRepository) Save(payload *model.Image) error {
	return im.db.Save(payload).Error
}

func (im *imageRepository) Update(payload *model.Image) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.DeletedAt.Time.String() != "" {
		updateFields["deleted_at"] = payload.DeletedAt
	}

	return im.db.Model(&model.Image{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (im *imageRepository) Search(by map[string]interface{}) ([]model.Image, error) {
	var image []model.Image
	query := im.db
	for key, value := range by {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	result := query.Find(&image).Error
	if result != nil {
		return nil, result
	}
	return image, nil
}

func (im *imageRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = im.db.Model(&model.Image{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = im.db.Model(&model.Image{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (im *imageRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.Image, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var image []model.Image
	err := im.db.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&image).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var totalRows int64
	err = im.db.Model(model.Image{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return image, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewImageRepository(db *gorm.DB, fb *firebase.App) ImageRepository {
	return &imageRepository{
		db: db,
		fb: fb,
	}
}
