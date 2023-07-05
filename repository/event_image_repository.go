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

type EventImageRepository interface {
	BaseRepository[model.EventImage]
	BaseRepositoryCount[model.EventImage]
	BaseRepositoryPaging[model.EventImage]
	FirebaseSave(payload multipart.File) (string, error)
	SaveTrx(payload *model.EventImage, tx *gorm.DB) error
	DeleteTrx(id string, tx *gorm.DB) error
}
type eventImageRepository struct {
	db *gorm.DB
	fb *firebase.App
}

func (ei *eventImageRepository) Delete(id string) error {
	return ei.db.Delete(&model.EventImage{}, "id=?", id).Error
}

func (ei *eventImageRepository) DeleteTrx(id string, tx *gorm.DB) error {
	return tx.Delete(&model.EventImage{}, "id=?", id).Error
}

func (ei *eventImageRepository) Get(id string) (*model.EventImage, error) {
	var eventImage model.EventImage
	result := ei.db.First(&eventImage, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &eventImage, nil
}

func (ei *eventImageRepository) List() ([]model.EventImage, error) {
	var eventImage []model.EventImage
	result := ei.db.Find(&eventImage).Error
	if result != nil {
		return nil, result
	}
	return eventImage, nil
}

func (ei *eventImageRepository) FirebaseSave(payload multipart.File) (string, error) {
	ctx := context.Background()

	// Create a storage reference for the file
	storageClient, err := ei.fb.Storage(ctx)
	if err != nil {
		return "", err
	}

	bucket, err := storageClient.Bucket(os.Getenv("BUCKET_NAME"))
	if err != nil {
		return "", err
	}

	// Generate a unique filename for the file
	filename := generateUniqueEventImagename()

	// Specify the path where the file will be stored in the bucket
	filePath := "eventImages/" + filename

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

func generateUniqueEventImagename() string {
	// Implement your own logic to generate a unique filename
	// You can use a timestamp, random string, or any other method
	// For example:
	return "eventImage_" + time.Now().Format("20060102150405") + ".jpg"
}

func (ei *eventImageRepository) Save(payload *model.EventImage) error {
	return ei.db.Save(payload).Error
}

func (r *eventImageRepository) SaveTrx(payload *model.EventImage, tx *gorm.DB) error {
	// If the provided transaction is not nil, use it for saving the EventImage
	if tx != nil {
		return tx.Create(payload).Error
	}

	// Otherwise, use the default DB connection for saving the EventImage
	return r.db.Create(payload).Error
}

func (ei *eventImageRepository) Update(payload *model.EventImage) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.DeletedAt.Time.String() != "" {
		updateFields["deleted_at"] = payload.DeletedAt
	}

	return ei.db.Model(&model.EventImage{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (ei *eventImageRepository) Search(by map[string]interface{}) ([]model.EventImage, error) {
	var eventImage []model.EventImage
	query := ei.db
	for key, value := range by {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	result := query.Find(&eventImage).Error
	if result != nil {
		return nil, result
	}
	return eventImage, nil
}

func (ei *eventImageRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = ei.db.Model(&model.EventImage{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = ei.db.Model(&model.EventImage{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (ei *eventImageRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.EventImage, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var eventImage []model.EventImage
	err := ei.db.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&eventImage).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var totalRows int64
	err = ei.db.Model(model.EventImage{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return eventImage, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewEventImageRepository(db *gorm.DB, fb *firebase.App) EventImageRepository {
	return &eventImageRepository{
		db: db,
		fb: fb,
	}
}
