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

type FileRepository interface {
	BaseRepository[model.File]
	BaseRepositoryCount[model.File]
	BaseRepositoryPaging[model.File]
	FirebaseSave(payload multipart.File) (string, error)
	SaveTrx(payload *model.File, tx *gorm.DB) error
	DeleteTrx(id string, tx *gorm.DB) error
}
type fileRepository struct {
	db *gorm.DB
	fb *firebase.App
}

func (fi *fileRepository) Delete(id string) error {
	return fi.db.Delete(&model.File{}, "id=?", id).Error
}

func (fi *fileRepository) DeleteTrx(id string, tx *gorm.DB) error {
	return tx.Delete(&model.File{}, "id=?", id).Error
}

func (fi *fileRepository) Get(id string) (*model.File, error) {
	var file model.File
	result := fi.db.First(&file, "id=?", id).Error
	if result != nil {
		return nil, result
	}
	return &file, nil
}

func (fi *fileRepository) List() ([]model.File, error) {
	var file []model.File
	result := fi.db.Find(&file).Error
	if result != nil {
		return nil, result
	}
	return file, nil
}

func (fi *fileRepository) FirebaseSave(payload multipart.File) (string, error) {
	ctx := context.Background()

	// Create a storage reference for the file
	storageClient, err := fi.fb.Storage(ctx)
	if err != nil {
		return "", err
	}

	bucket, err := storageClient.Bucket(os.Getenv("BUCKET_NAME"))
	if err != nil {
		return "", err
	}

	// Generate a unique filename for the file
	filename := generateUniqueFilename()

	// Specify the path where the file will be stored in the bucket
	filePath := "files/" + filename

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

func generateUniqueFilename() string {
	// Implement your own logic to generate a unique filename
	// You can use a timestamp, random string, or any other method
	// For example:
	return "file_" + time.Now().Format("20060102150405") + ".pdf"
}

func (fi *fileRepository) Save(payload *model.File) error {
	return fi.db.Save(payload).Error
}

func (fi *fileRepository) SaveTrx(payload *model.File, tx *gorm.DB) error {
	return tx.Create(payload).Error
}

func (fi *fileRepository) Update(payload *model.File) error {
	updateFields := make(map[string]interface{})

	// Add fields to be updated based on the payload
	if payload.DeletedAt.Time.String() != "" {
		updateFields["deleted_at"] = payload.DeletedAt
	}

	return fi.db.Model(&model.File{}).Where("id = ?", payload.ID).Updates(updateFields).Error
}

func (fi *fileRepository) Search(by map[string]interface{}) ([]model.File, error) {
	var file []model.File
	query := fi.db
	for key, value := range by {
		// Perform case-insensitive search using ilike
		query = query.Where(fmt.Sprintf("%s ilike ?", key), fmt.Sprintf("%%%v%%", value))
	}
	result := query.Find(&file).Error
	if result != nil {
		return nil, result
	}
	return file, nil
}

func (fi *fileRepository) CountData(fieldname string, id string) error {
	var count int64
	var result *gorm.DB

	if id != "" {
		result = fi.db.Model(&model.File{}).Where("name ilike ? AND id <> ?", fieldname, id).Count(&count)
	} else {
		result = fi.db.Model(&model.File{}).Where("name ilike ?", fieldname).Count(&count)
	}
	if result.Error != nil {
		return result.Error
	}
	if count > 0 {
		return fmt.Errorf("name %s already exist", fieldname)
	}
	return nil
}

func (fi *fileRepository) Paging(requestQueryParam dto.RequestQueryParams) ([]model.File, dto.Paging, error) {
	paginationQuery, orderQuery := pagingValidate(requestQueryParam)

	var file []model.File
	err := fi.db.Order(orderQuery).Limit(paginationQuery.Take).Offset(paginationQuery.Skip).Find(&file).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	var totalRows int64
	err = fi.db.Model(model.File{}).Count(&totalRows).Error
	if err != nil {
		return nil, dto.Paging{}, err
	}
	return file, common.Paginate(paginationQuery.Page, paginationQuery.Take, int(totalRows)), nil
}

func NewFileRepository(db *gorm.DB, fb *firebase.App) FileRepository {
	return &fileRepository{
		db: db,
		fb: fb,
	}
}
