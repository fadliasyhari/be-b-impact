package dto

import "time"

type ContentDTO struct {
	ID        string     `json:"id"`
	Title     string     `json:"title"`
	Body      string     `json:"body"`
	Author    string     `json:"author"`
	Excerpt   string     `json:"excerpt"`
	Status    string     `json:"status"`
	Category  string     `json:"category"`
	ImageURLs []ImageDTO `json:"image_urls"`
	Tags      []TagDTO   `json:"tags"`
	CreatedAt time.Time  `json:"created_at"`
}

type ImageDTO struct {
	ImageURL string `json:"image_url"`
}

type CategoryDTO struct {
	Name string `json:"name"`
}

type TagDTO struct {
	Name string `json:"name"`
}
