package response

import (
	"be-b-impact.com/csr/model"
	"be-b-impact.com/csr/model/dto"
)

func MapContentToResponse(content *model.Content) dto.ContentDTO {
	res := dto.ContentDTO{
		ID:       content.ID,
		Title:    content.Title,
		Body:     content.Body,
		Author:   content.Author,
		Excerpt:  content.Excerpt,
		Status:   content.Status,
		Category: content.Category.Name,
		CategoryDetail: dto.CategoryDTO{
			ID:   content.Category.ID,
			Name: content.Category.Name,
		},
		CreatedAt: content.CreatedAt,
	}
	for _, v := range content.TagsContent {
		res.Tags = append(res.Tags, dto.TagDTO{
			ID:   v.Tag.ID,
			Name: v.Tag.Name,
		})
	}
	for _, v := range content.Image {
		res.ImageURLs = append(res.ImageURLs, dto.ImageDTO{
			ID:       v.ID,
			ImageURL: v.ImageURL,
		})
	}
	return res
}
