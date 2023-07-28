package usecase

import (
	"fmt"
	"mime/multipart"

	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
	"be-b-impact.com/csr/repository"
)

type ContentUseCase interface {
	BaseUseCase[model.Content]
	BaseUseCasePaging[model.Content]
	SaveContent(payload *model.Content, tags []string, file multipart.File) error
	UpdateContent(payload *model.Content, tags []string, file multipart.File) error
}

type contentUseCase struct {
	repo           repository.ContentRepository
	tagsContentUC  TagsContentUseCase
	imageUC        ImageUseCase
	notificationUC NotificationUseCase
}

// SaveData implements ContentUseCase.
func (*contentUseCase) SaveData(payload *model.Content) error {
	panic("unimplemented")
}

func (co *contentUseCase) DeleteData(id string) error {
	content, err := co.FindById(id)
	if err != nil {
		return fmt.Errorf("content with ID %s not found", id)
	}

	co.notificationUC.DeleteNotif("content", content.ID)

	return co.repo.Delete(content.ID)
}

func (co *contentUseCase) FindAll() ([]model.Content, error) {
	return co.repo.List()
}

func (co *contentUseCase) FindById(id string) (*model.Content, error) {
	content, err := co.repo.Get(id)
	if err != nil {
		return nil, fmt.Errorf("content with ID %s not found", id)
	}
	return content, nil
}

func (co *contentUseCase) SaveContent(payload *model.Content, tags []string, file multipart.File) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }

	if payload.Status == "1" && (payload.Title == "" || payload.Body == "" || payload.Author == "" || payload.Excerpt == "" || payload.CategoryID == nil) {
		return fmt.Errorf("form is not completed")
	}

	tx := co.repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := co.repo.Save(payload); err != nil {
		tx.Rollback()
		return err
	}

	if len(tags) > 0 {
		for _, v := range tags {
			tcPayload := model.TagsContent{
				TagID:     v,
				ContentID: payload.ID,
			}
			if err := co.tagsContentUC.SaveTagsContent(&tcPayload, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if file != nil {
		imageURL, err := co.imageUC.FirebaseUpload(file)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Create the image payload
		imagePayload := model.Image{
			ImageURL:  imageURL,
			ContentID: payload.ID,
		}
		if err := co.imageUC.SaveImage(&imagePayload, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if payload.Status == "1" {
		notifPayload := model.Notification{
			Body:   fmt.Sprintf("New content: <b>%s</b>", payload.Title),
			UserID: "0",
			Type:   "content",
			TypeID: payload.ID,
		}
		err := co.notificationUC.SaveNotifDetail(&notifPayload, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (co *contentUseCase) UpdateData(payload *model.Content) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }
	// cek jika data sudah ada -> count > 0
	if payload.ID != "" {
		_, err := co.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("content with ID %s not found", payload.ID)
		}
	}

	return co.repo.Update(payload)
}

func (co *contentUseCase) UpdateContent(payload *model.Content, tags []string, file multipart.File) error {
	// err := payload.Vaildate()
	// if err != nil {
	// 	return err
	// }

	if payload.ID != "" {
		_, err := co.FindById(payload.ID)
		if err != nil {
			return fmt.Errorf("content with ID %s not found", payload.ID)
		}
	}

	tx := co.repo.BeginTransaction()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := co.repo.Update(payload); err != nil {
		tx.Rollback()
		return err
	}

	for _, tag := range payload.TagsContent {
		if !contains(tags, tag.ID) {
			if err := co.tagsContentUC.DeleteDataTrx(tag.ID, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	for _, tagId := range tags {
		found := false
		for _, tag := range payload.TagsContent {
			if tag.ID == tagId {
				found = true
				break
			}
		}
		if !found {
			tcPayload := model.TagsContent{
				TagID:     tagId,
				ContentID: payload.ID,
			}
			if err := co.tagsContentUC.SaveTagsContent(&tcPayload, tx); err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	if file != nil {
		if len(payload.Image) > 0 {
			for _, v := range payload.Image {
				if err := co.imageUC.DeleteDataTrx(v.ID, tx); err != nil {
					tx.Rollback()
					return err
				}
			}
		}

		imageURL, err := co.imageUC.FirebaseUpload(file)
		if err != nil {
			tx.Rollback()
			return err
		}

		// Create the image payload
		imagePayload := model.Image{
			ImageURL:  imageURL,
			ContentID: payload.ID,
		}
		if err := co.imageUC.SaveImage(&imagePayload, tx); err != nil {
			tx.Rollback()
			return err
		}
	}

	if payload.Status == "1" && (payload.Title == "" || payload.Body == "" || payload.Author == "" || payload.Excerpt == "" || payload.CategoryID == nil) {
		return fmt.Errorf("form is not completed")
	}

	if payload.Status == "1" {
		notifPayload := model.Notification{
			Body:   fmt.Sprintf("New content: <b>%s</b>", payload.Title),
			UserID: "0",
			Type:   "content",
			TypeID: payload.ID,
		}
		err := co.notificationUC.SaveNotifDetail(&notifPayload, tx)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func contains(slice []string, element string) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

func (co *contentUseCase) SearchBy(by map[string]interface{}) ([]model.Content, error) {
	contents, err := co.repo.Search(by)
	if err != nil {
		return nil, fmt.Errorf("request invalid")
	}
	return contents, nil
}

func (co *contentUseCase) Pagination(requestQueryParams dto.RequestQueryParams) ([]model.Content, dto.Paging, error) {
	if !requestQueryParams.QueryParams.IsSortValid() {
		return nil, dto.Paging{}, fmt.Errorf("invalid sort by: %s", requestQueryParams.QueryParams.Sort)
	}
	return co.repo.Paging(requestQueryParams)
}

func NewContentUseCase(repo repository.ContentRepository, tagsContentUC TagsContentUseCase, imageUC ImageUseCase, notificationUC NotificationUseCase) ContentUseCase {
	return &contentUseCase{
		repo:           repo,
		tagsContentUC:  tagsContentUC,
		imageUC:        imageUC,
		notificationUC: notificationUC,
	}
}
